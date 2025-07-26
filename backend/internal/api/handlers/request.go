package handlers

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"volt-backend/internal/db"
	"volt-backend/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetRequests(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	collectionID, err := primitive.ObjectIDFromHex(c.Param("collectionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	// Check if user has access to collection
	if !hasCollectionAccess(userID, collectionID, c) {
		return
	}

	collection := db.GetCollection("requests")
	cursor, err := collection.Find(context.Background(), bson.M{"collection_id": collectionID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch requests"})
		return
	}
	defer cursor.Close(context.Background())

	var requests []models.Request
	if err = cursor.All(context.Background(), &requests); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode requests"})
		return
	}

	c.JSON(http.StatusOK, requests)
}

func CreateRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	collectionID, err := primitive.ObjectIDFromHex(c.Param("collectionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	var req models.CreateRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user has access to collection
	if !hasCollectionAccess(userID, collectionID, c) {
		return
	}

	request := models.Request{
		Name:         req.Name,
		Method:       req.Method,
		URL:          req.URL,
		Headers:      req.Headers,
		Body:         req.Body,
		Auth:         req.Auth,
		Tests:        req.Tests,
		CollectionID: collectionID,
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	if request.Tests == nil {
		request.Tests = []string{}
	}

	requestDB := db.GetCollection("requests")
	result, err := requestDB.InsertOne(context.Background(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	request.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, request)
}

func GetRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	requestID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	requestDB := db.GetCollection("requests")
	var request models.Request
	err = requestDB.FindOne(context.Background(), bson.M{"_id": requestID}).Decode(&request)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	// Check if user has access to collection
	if !hasCollectionAccess(userID, request.CollectionID, c) {
		return
	}

	c.JSON(http.StatusOK, request)
}

func UpdateRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	requestID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var req models.CreateRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestDB := db.GetCollection("requests")
	
	// Check if request exists and user has access
	var request models.Request
	err = requestDB.FindOne(context.Background(), bson.M{"_id": requestID}).Decode(&request)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	// Check collection access
	if !hasCollectionAccess(userID, request.CollectionID, c) {
		return
	}

	update := bson.M{
		"$set": bson.M{
			"name":       req.Name,
			"method":     req.Method,
			"url":        req.URL,
			"headers":    req.Headers,
			"body":       req.Body,
			"auth":       req.Auth,
			"tests":      req.Tests,
			"updated_at": time.Now(),
		},
	}

	_, err = requestDB.UpdateOne(context.Background(), bson.M{"_id": requestID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request updated successfully"})
}

func DeleteRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	requestID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	requestDB := db.GetCollection("requests")
	
	// Check if request exists and user has access
	var request models.Request
	err = requestDB.FindOne(context.Background(), bson.M{"_id": requestID}).Decode(&request)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	// Check collection access
	if !hasCollectionAccess(userID, request.CollectionID, c) {
		return
	}

	// Delete the request
	_, err = requestDB.DeleteOne(context.Background(), bson.M{"_id": requestID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request deleted successfully"})
}

func ExecuteRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	requestID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	requestDB := db.GetCollection("requests")
	var request models.Request
	err = requestDB.FindOne(context.Background(), bson.M{"_id": requestID}).Decode(&request)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	// Check collection access
	if !hasCollectionAccess(userID, request.CollectionID, c) {
		return
	}

	// Execute the HTTP request
	startTime := time.Now()
	
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	var body io.Reader
	if request.Body.Content != "" {
		body = strings.NewReader(request.Body.Content)
	}

	httpReq, err := http.NewRequest(string(request.Method), request.URL, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request URL or method"})
		return
	}

	// Add headers
	for key, value := range request.Headers {
		httpReq.Header.Set(key, value)
	}

	// Add authentication
	switch request.Auth.Type {
	case models.Bearer:
		if token, ok := request.Auth.Credentials["token"]; ok {
			httpReq.Header.Set("Authorization", "Bearer "+token)
		}
	case models.Basic:
		if username, ok := request.Auth.Credentials["username"]; ok {
			if password, ok := request.Auth.Credentials["password"]; ok {
				httpReq.SetBasicAuth(username, password)
			}
		}
	case models.APIKey:
		if key, ok := request.Auth.Credentials["key"]; ok {
			if header, ok := request.Auth.Credentials["header"]; ok {
				httpReq.Header.Set(header, key)
			}
		}
	}

	// Execute request
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to execute request: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	executionTime := time.Since(startTime).Milliseconds()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	// Convert headers to map
	responseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	response := models.ExecuteRequestResponse{
		Status:  resp.StatusCode,
		Headers: responseHeaders,
		Body:    string(responseBody),
		Time:    executionTime,
	}

	// Store execution history
	history := map[string]interface{}{
		"request_id":    requestID,
		"user_id":       userID,
		"request_data":  request,
		"response_data": response,
		"executed_at":   time.Now(),
	}

	historyDB := db.GetCollection("execution_history")
	historyDB.InsertOne(context.Background(), history)

	c.JSON(http.StatusOK, response)
}

// Helper function to check if user has access to collection
func hasCollectionAccess(userID, collectionID primitive.ObjectID, c *gin.Context) bool {
	collectionDB := db.GetCollection("collections")
	var collection models.Collection
	err := collectionDB.FindOne(context.Background(), bson.M{"_id": collectionID}).Decode(&collection)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return false
	}

	// Check workspace access
	workspaceDB := db.GetCollection("workspaces")
	var workspace models.Workspace
	err = workspaceDB.FindOne(context.Background(), bson.M{
		"_id": collection.WorkspaceID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": userID},
		},
	}).Decode(&workspace)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Access denied"})
		return false
	}

	return true
}