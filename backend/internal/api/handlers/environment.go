package handlers

import (
	"context"
	"net/http"
	"time"

	"volt-backend/internal/db"
	"volt-backend/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetEnvironments gets all environments for a workspace
func GetEnvironments(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	workspaceID, err := primitive.ObjectIDFromHex(c.Param("workspaceId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	// Check if user has access to workspace
	workspaceCollection := db.GetCollection("workspaces")
	var workspace models.Workspace
	err = workspaceCollection.FindOne(context.Background(), bson.M{
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

	collection := db.GetCollection("environments")
	cursor, err := collection.Find(context.Background(), bson.M{
		"workspace_id": workspaceID.Hex(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch environments"})
		return
	}
	defer cursor.Close(context.Background())

	var environments []models.Environment
	if err = cursor.All(context.Background(), &environments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode environments"})
		return
	}

	if environments == nil {
		environments = []models.Environment{}
	}

	c.JSON(http.StatusOK, environments)
}

// CreateEnvironment creates a new environment
func CreateEnvironment(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	workspaceID, err := primitive.ObjectIDFromHex(c.Param("workspaceId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	// Check if user has access to workspace
	workspaceCollection := db.GetCollection("workspaces")
	var workspace models.Workspace
	err = workspaceCollection.FindOne(context.Background(), bson.M{
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

	var req models.CreateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Initialize variables slice if nil
	if req.Variables == nil {
		req.Variables = []models.Variable{}
	}

	environment := models.Environment{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		WorkspaceID: workspaceID.Hex(),
		Variables:   req.Variables,
		CreatedBy:   userID.Hex(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	collection := db.GetCollection("environments")
	_, err = collection.InsertOne(context.Background(), environment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create environment"})
		return
	}

	c.JSON(http.StatusCreated, environment)
}

// GetEnvironment gets a specific environment
func GetEnvironment(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	environmentID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(environmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	collection := db.GetCollection("environments")
	var environment models.Environment
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&environment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch environment"})
		return
	}

	// Check if user has access to workspace
	workspaceID, err := primitive.ObjectIDFromHex(environment.WorkspaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID in environment"})
		return
	}

	workspaceCollection := db.GetCollection("workspaces")
	var workspace models.Workspace
	err = workspaceCollection.FindOne(context.Background(), bson.M{
		"_id": workspaceID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": userID},
		},
	}).Decode(&workspace)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to environment"})
		return
	}

	c.JSON(http.StatusOK, environment)
}

// UpdateEnvironment updates an environment
func UpdateEnvironment(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	environmentID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(environmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	var req models.UpdateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := db.GetCollection("environments")

	// First, get the environment to check workspace access
	var environment models.Environment
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&environment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch environment"})
		return
	}

	// Check if user has access to workspace
	workspaceID, err := primitive.ObjectIDFromHex(environment.WorkspaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID in environment"})
		return
	}

	workspaceCollection := db.GetCollection("workspaces")
	var workspace models.Workspace
	err = workspaceCollection.FindOne(context.Background(), bson.M{
		"_id": workspaceID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": userID},
		},
	}).Decode(&workspace)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to environment"})
		return
	}

	// Initialize variables slice if nil
	if req.Variables == nil {
		req.Variables = []models.Variable{}
	}

	// Update the environment
	update := bson.M{
		"$set": bson.M{
			"name":       req.Name,
			"variables":  req.Variables,
			"updated_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update environment"})
		return
	}

	// Return updated environment
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&environment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated environment"})
		return
	}

	c.JSON(http.StatusOK, environment)
}

// DeleteEnvironment deletes an environment
func DeleteEnvironment(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	environmentID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(environmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	collection := db.GetCollection("environments")

	// First, get the environment to check workspace access
	var environment models.Environment
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&environment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch environment"})
		return
	}

	// Check if user has access to workspace
	workspaceID, err := primitive.ObjectIDFromHex(environment.WorkspaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID in environment"})
		return
	}

	workspaceCollection := db.GetCollection("workspaces")
	var workspace models.Workspace
	err = workspaceCollection.FindOne(context.Background(), bson.M{
		"_id": workspaceID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": userID},
		},
	}).Decode(&workspace)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to environment"})
		return
	}

	// Delete the environment
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete environment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Environment deleted successfully"})
}