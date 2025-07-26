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