package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"volt-backend/internal/db"
	"volt-backend/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetWorkspaces(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)

	collection := db.GetCollection("workspaces")
	cursor, err := collection.Find(context.Background(), bson.M{
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": userID},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch workspaces"})
		return
	}
	defer cursor.Close(context.Background())

	var workspaces []models.Workspace
	if err = cursor.All(context.Background(), &workspaces); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode workspaces"})
		return
	}

	c.JSON(http.StatusOK, workspaces)
}

func CreateWorkspace(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)

	var req models.CreateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workspace := models.Workspace{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID,
		Members:     []primitive.ObjectID{userID},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	collection := db.GetCollection("workspaces")
	result, err := collection.InsertOne(context.Background(), workspace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workspace"})
		return
	}

	workspace.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, workspace)
}

func GetWorkspace(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	workspaceID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	collection := db.GetCollection("workspaces")
	var workspace models.Workspace
	err = collection.FindOne(context.Background(), bson.M{
		"_id": workspaceID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": userID},
		},
	}).Decode(&workspace)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}

	c.JSON(http.StatusOK, workspace)
}

func UpdateWorkspace(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	workspaceID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var req models.CreateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := db.GetCollection("workspaces")
	update := bson.M{
		"$set": bson.M{
			"name":        req.Name,
			"description": req.Description,
			"updated_at":  time.Now(),
		},
	}

	result, err := collection.UpdateOne(context.Background(), bson.M{
		"_id":      workspaceID,
		"owner_id": userID,
	}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workspace"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found or not authorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workspace updated successfully"})
}

func DeleteWorkspace(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	workspaceID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	collection := db.GetCollection("workspaces")
	result, err := collection.DeleteOne(context.Background(), bson.M{
		"_id":      workspaceID,
		"owner_id": userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workspace"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found or not authorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workspace deleted successfully"})
}

// generateRandomToken creates a cryptographically secure random token
func generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateInvitationLink(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	workspaceID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var req models.CreateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate role
	validRoles := map[string]bool{
		"member": true,
		"admin":  true,
		"viewer": true,
	}
	if !validRoles[req.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	// Check if user has permission to invite (owner or admin)
	workspaceCollection := db.GetCollection("workspaces")
	var workspace models.Workspace
	err = workspaceCollection.FindOne(context.Background(), bson.M{
		"_id": workspaceID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": userID}, // For now, any member can invite
		},
	}).Decode(&workspace)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found or no permission"})
		return
	}

	// Generate random token
	token, err := generateRandomToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Calculate expiry time
	var expiresAt *time.Time
	if req.ExpiresInHours > 0 {
		expiry := time.Now().Add(time.Duration(req.ExpiresInHours) * time.Hour)
		expiresAt = &expiry
	}

	// Create invitation
	invitation := models.WorkspaceInvitation{
		WorkspaceID: workspaceID,
		Token:       token,
		Role:        req.Role,
		CreatedBy:   userID,
		ExpiresAt:   expiresAt,
		Used:        false,
		CreatedAt:   time.Now(),
	}

	// Save to database
	invitationCollection := db.GetCollection("workspace_invitations")
	result, err := invitationCollection.InsertOne(context.Background(), invitation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invitation"})
		return
	}

	invitation.ID = result.InsertedID.(primitive.ObjectID)

	// Generate full invitation link
	baseURL := c.GetHeader("Origin")
	if baseURL == "" {
		baseURL = "http://localhost:5173" // Default for development
	}
	inviteLink := fmt.Sprintf("%s/invite/%s", baseURL, token)

	response := models.InvitationResponse{
		Token:       token,
		InviteLink:  inviteLink,
		Role:        req.Role,
		ExpiresAt:   expiresAt,
		WorkspaceID: workspaceID.Hex(),
	}

	c.JSON(http.StatusCreated, response)
}

func AcceptInvitation(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	token := c.Param("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invitation token"})
		return
	}

	// Find invitation
	invitationCollection := db.GetCollection("workspace_invitations")
	var invitation models.WorkspaceInvitation
	err := invitationCollection.FindOne(context.Background(), bson.M{
		"token": token,
		"used":  false,
	}).Decode(&invitation)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired invitation"})
		return
	}

	// Check if invitation is expired
	if invitation.ExpiresAt != nil && time.Now().After(*invitation.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invitation has expired"})
		return
	}

	// Check if user is already a member
	workspaceCollection := db.GetCollection("workspaces")
	var workspace models.Workspace
	err = workspaceCollection.FindOne(context.Background(), bson.M{
		"_id":     invitation.WorkspaceID,
		"members": userID,
	}).Decode(&workspace)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You are already a member of this workspace"})
		return
	}

	// Add user to workspace
	now := time.Now()
	_, err = workspaceCollection.UpdateOne(context.Background(), bson.M{
		"_id": invitation.WorkspaceID,
	}, bson.M{
		"$addToSet": bson.M{"members": userID},
		"$set":      bson.M{"updated_at": now},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to workspace"})
		return
	}

	// Mark invitation as used
	_, err = invitationCollection.UpdateOne(context.Background(), bson.M{
		"_id": invitation.ID,
	}, bson.M{
		"$set": bson.M{
			"used":    true,
			"used_by": userID,
			"used_at": now,
		},
	})
	if err != nil {
		// Log error but don't fail the request since user was already added
		fmt.Printf("Failed to mark invitation as used: %v\n", err)
	}

	// Get workspace details to return
	err = workspaceCollection.FindOne(context.Background(), bson.M{
		"_id": invitation.WorkspaceID,
	}).Decode(&workspace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get workspace details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Successfully joined workspace",
		"workspace": workspace,
		"role":      invitation.Role,
	})
}