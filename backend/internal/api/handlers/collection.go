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

func GetCollections(c *gin.Context) {
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

	collection := db.GetCollection("collections")
	cursor, err := collection.Find(context.Background(), bson.M{"workspace_id": workspaceID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch collections"})
		return
	}
	defer cursor.Close(context.Background())

	var collections []models.Collection
	if err = cursor.All(context.Background(), &collections); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode collections"})
		return
	}

	c.JSON(http.StatusOK, collections)
}

func CreateCollection(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	workspaceID, err := primitive.ObjectIDFromHex(c.Param("workspaceId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var req models.CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	collection := models.Collection{
		Name:        req.Name,
		Description: req.Description,
		WorkspaceID: workspaceID,
		CreatedBy:   userID,
		Requests:    []primitive.ObjectID{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	collectionDB := db.GetCollection("collections")
	result, err := collectionDB.InsertOne(context.Background(), collection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create collection"})
		return
	}

	collection.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, collection)
}

func GetCollection(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	collectionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	collectionDB := db.GetCollection("collections")
	var collection models.Collection
	err = collectionDB.FindOne(context.Background(), bson.M{"_id": collectionID}).Decode(&collection)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	// Check if user has access to workspace
	workspaceCollection := db.GetCollection("workspaces")
	var workspace models.Workspace
	err = workspaceCollection.FindOne(context.Background(), bson.M{
		"_id": collection.WorkspaceID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": userID},
		},
	}).Decode(&workspace)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, collection)
}

func UpdateCollection(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	collectionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	var req models.CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collectionDB := db.GetCollection("collections")
	
	// Check if collection exists and user has access
	var collection models.Collection
	err = collectionDB.FindOne(context.Background(), bson.M{"_id": collectionID}).Decode(&collection)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	// Check workspace access
	workspaceCollection := db.GetCollection("workspaces")
	var workspace models.Workspace
	err = workspaceCollection.FindOne(context.Background(), bson.M{
		"_id": collection.WorkspaceID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": userID},
		},
	}).Decode(&workspace)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Access denied"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"name":        req.Name,
			"description": req.Description,
			"updated_at":  time.Now(),
		},
	}

	_, err = collectionDB.UpdateOne(context.Background(), bson.M{"_id": collectionID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update collection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collection updated successfully"})
}

func DeleteCollection(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	collectionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	collectionDB := db.GetCollection("collections")
	
	// Check if collection exists and user has access
	var collection models.Collection
	err = collectionDB.FindOne(context.Background(), bson.M{"_id": collectionID}).Decode(&collection)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	// Check workspace access
	workspaceCollection := db.GetCollection("workspaces")
	var workspace models.Workspace
	err = workspaceCollection.FindOne(context.Background(), bson.M{
		"_id": collection.WorkspaceID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": userID},
		},
	}).Decode(&workspace)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Access denied"})
		return
	}

	// Delete all requests in this collection first
	requestDB := db.GetCollection("requests")
	_, err = requestDB.DeleteMany(context.Background(), bson.M{"collection_id": collectionID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete collection requests"})
		return
	}

	// Delete the collection
	_, err = collectionDB.DeleteOne(context.Background(), bson.M{"_id": collectionID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete collection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collection deleted successfully"})
}