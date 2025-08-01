package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"volt-backend/internal/auth"
	"volt-backend/internal/db"
	"volt-backend/internal/models"
)

// GetProjects retrieves all projects for the authenticated user
func GetProjects(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find projects where user is owner or member
	filter := bson.M{
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": bson.M{"$elemMatch": bson.M{"user_id": userID}}},
		},
	}

	cursor, err := db.GetCollection("projects").Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}
	defer cursor.Close(ctx)

	var projects []models.Project
	if err = cursor.All(ctx, &projects); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse projects"})
		return
	}

	if projects == nil {
		projects = []models.Project{}
	}

	c.JSON(http.StatusOK, projects)
}

// GetProject retrieves a specific project by ID
func GetProject(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user has access to this project
	filter := bson.M{
		"_id": objID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": bson.M{"$elemMatch": bson.M{"user_id": userID}}},
		},
	}

	var project models.Project
	err = db.GetCollection("projects").FindOne(ctx, filter).Decode(&project)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// CreateProject creates a new project
func CreateProject(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Name         string                   `json:"name" binding:"required"`
		Description  string                   `json:"description"`
		Type         string                   `json:"type" binding:"required"`
		Priority     string                   `json:"priority"`
		StartDate    time.Time                `json:"start_date"`
		DueDate      *time.Time               `json:"due_date"`
		Tags         []string                 `json:"tags"`
		Settings     models.ProjectSettings   `json:"settings"`
		TeamMembers  []string                 `json:"team_members"`
		Milestones   []models.ProjectMilestone `json:"milestones"`
		Tasks        []models.ProjectTask     `json:"tasks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	project := models.Project{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Status:      "planning",
		Priority:    req.Priority,
		Progress:    0,
		StartDate:   req.StartDate,
		DueDate:     req.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
		OwnerID:     userID,
		Tags:        req.Tags,
		Settings:    req.Settings,
		Milestones:  req.Milestones,
		Tasks:       req.Tasks,
		Stats: models.ProjectStats{
			TotalTasks:     int64(len(req.Tasks)),
			CompletedTasks: 0,
			OverdueTasks:   0,
			ActiveMembers:  1,
			TotalAPICalls:  0,
			LastActivity:   now,
		},
	}

	// Add owner as first member
	project.Members = []models.ProjectMember{
		{
			ID:       userID,
			Name:     "Project Owner", // TODO: Get actual user name from user collection
			Email:    "owner@example.com",
			Role:     "owner",
			Status:   "online",
			LastSeen: now,
		},
	}

	// Add team members if provided
	for _, memberEmail := range req.TeamMembers {
		// TODO: Look up user by email and add as member
		project.Members = append(project.Members, models.ProjectMember{
			ID:       memberEmail, // This should be user ID, not email
			Name:     memberEmail,
			Email:    memberEmail,
			Role:     "member",
			Status:   "offline",
			LastSeen: now,
		})
	}

	// Add project ID to milestones and tasks
	for i := range project.Milestones {
		project.Milestones[i].ID = primitive.NewObjectID().Hex()
		project.Milestones[i].ProjectID = project.ID.Hex()
	}

	for i := range project.Tasks {
		project.Tasks[i].ID = primitive.NewObjectID().Hex()
		project.Tasks[i].ProjectID = project.ID.Hex()
		project.Tasks[i].CreatedAt = now
		project.Tasks[i].UpdatedAt = now
	}

	_, err := db.GetCollection("projects").InsertOne(ctx, project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// UpdateProject updates an existing project
func UpdateProject(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req struct {
		Name        string     `json:"name"`
		Description string     `json:"description"`
		Type        string     `json:"type"`
		Priority    string     `json:"priority"`
		Status      string     `json:"status"`
		Progress    int        `json:"progress"`
		DueDate     *time.Time `json:"due_date"`
		Tags        []string   `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user is owner or has admin role
	filter := bson.M{
		"_id": objID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": bson.M{"$elemMatch": bson.M{"user_id": userID, "role": bson.M{"$in": []string{"owner", "admin"}}}}},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	// Only update fields that are provided
	if req.Name != "" {
		update["$set"].(bson.M)["name"] = req.Name
	}
	if req.Description != "" {
		update["$set"].(bson.M)["description"] = req.Description
	}
	if req.Type != "" {
		update["$set"].(bson.M)["type"] = req.Type
	}
	if req.Priority != "" {
		update["$set"].(bson.M)["priority"] = req.Priority
	}
	if req.Status != "" {
		update["$set"].(bson.M)["status"] = req.Status
	}
	if req.Progress >= 0 {
		update["$set"].(bson.M)["progress"] = req.Progress
	}
	if req.DueDate != nil {
		update["$set"].(bson.M)["due_date"] = req.DueDate
	}
	if req.Tags != nil {
		update["$set"].(bson.M)["tags"] = req.Tags
	}

	result, err := db.GetCollection("projects").UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}

// DeleteProject deletes a project
func DeleteProject(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Only project owner can delete
	filter := bson.M{
		"_id":      objID,
		"owner_id": userID,
	}

	result, err := db.GetCollection("projects").DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// AddProjectMember adds a member to a project
func AddProjectMember(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req struct {
		Email string `json:"email" binding:"required"`
		Role  string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user has permission to add members (owner or admin)
	filter := bson.M{
		"_id": objID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": bson.M{"$elemMatch": bson.M{"user_id": userID, "role": bson.M{"$in": []string{"owner", "admin"}}}}},
		},
	}

	// TODO: Look up user by email to get user ID and name
	newMember := models.ProjectMember{
		ID:       req.Email, // This should be actual user ID
		Name:     req.Email, // This should be actual user name
		Email:    req.Email,
		Role:     req.Role,
		Status:   "offline",
		LastSeen: time.Now(),
	}

	update := bson.M{
		"$push": bson.M{
			"members": newMember,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := db.GetCollection("projects").UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}

// GetProjectTasks retrieves all tasks for a project
func GetProjectTasks(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verify user has access to project
	filter := bson.M{
		"_id": objID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": bson.M{"$elemMatch": bson.M{"user_id": userID}}},
		},
	}

	var project models.Project
	err = db.GetCollection("projects").FindOne(ctx, filter).Decode(&project)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project.Tasks)
}

// CreateProjectTask creates a new task in a project
func CreateProjectTask(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req models.ProjectTask
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user has access to project
	filter := bson.M{
		"_id": objID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": bson.M{"$elemMatch": bson.M{"user_id": userID}}},
		},
	}

	now := time.Now()
	req.ID = primitive.NewObjectID().Hex()
	req.ProjectID = projectID
	req.CreatedAt = now
	req.UpdatedAt = now

	update := bson.M{
		"$push": bson.M{
			"tasks": req,
		},
		"$inc": bson.M{
			"stats.total_tasks": 1,
		},
		"$set": bson.M{
			"updated_at":           now,
			"stats.last_activity": now,
		},
	}

	result, err := db.GetCollection("projects").UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or access denied"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// UpdateProjectTask updates a task status and details
func UpdateProjectTask(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectID := c.Param("id")
	taskID := c.Param("taskId")
	
	objID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		Priority    string     `json:"priority"`
		AssignedTo  string     `json:"assigned_to"`
		DueDate     *time.Time `json:"due_date"`
		Tags        []string   `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build update query for the specific task in the array
	update := bson.M{
		"$set": bson.M{
			"updated_at":           time.Now(),
			"stats.last_activity": time.Now(),
		},
	}

	// Update specific task fields
	if req.Title != "" {
		update["$set"].(bson.M)["tasks.$.title"] = req.Title
	}
	if req.Description != "" {
		update["$set"].(bson.M)["tasks.$.description"] = req.Description
	}
	if req.Status != "" {
		update["$set"].(bson.M)["tasks.$.status"] = req.Status
	}
	if req.Priority != "" {
		update["$set"].(bson.M)["tasks.$.priority"] = req.Priority
	}
	if req.AssignedTo != "" {
		update["$set"].(bson.M)["tasks.$.assigned_to"] = req.AssignedTo
	}
	if req.DueDate != nil {
		update["$set"].(bson.M)["tasks.$.due_date"] = req.DueDate
	}
	if req.Tags != nil {
		update["$set"].(bson.M)["tasks.$.tags"] = req.Tags
	}

	update["$set"].(bson.M)["tasks.$.updated_at"] = time.Now()

	filter := bson.M{
		"_id":      objID,
		"tasks.id": taskID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": bson.M{"$elemMatch": bson.M{"user_id": userID}}},
		},
	}

	result, err := db.GetCollection("projects").UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project or task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

// GetProjectMilestones retrieves all milestones for a project
func GetProjectMilestones(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verify user has access to project
	filter := bson.M{
		"_id": objID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": bson.M{"$elemMatch": bson.M{"user_id": userID}}},
		},
	}

	var project models.Project
	err = db.GetCollection("projects").FindOne(ctx, filter).Decode(&project)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project.Milestones)
}

// CreateProjectMilestone creates a new milestone in a project
func CreateProjectMilestone(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req models.ProjectMilestone
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user has access to project
	filter := bson.M{
		"_id": objID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": bson.M{"$elemMatch": bson.M{"user_id": userID}}},
		},
	}

	req.ID = primitive.NewObjectID().Hex()
	req.ProjectID = projectID

	update := bson.M{
		"$push": bson.M{
			"milestones": req,
		},
		"$set": bson.M{
			"updated_at":           time.Now(),
			"stats.last_activity": time.Now(),
		},
	}

	result, err := db.GetCollection("projects").UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create milestone"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or access denied"})
		return
	}

	c.JSON(http.StatusCreated, req)
}