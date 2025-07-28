package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"volt/internal/models"
)

type WebhookHandler struct {
	db *mongo.Database
}

func NewWebhookHandler(db *mongo.Database) *WebhookHandler {
	return &WebhookHandler{db: db}
}

// generateWebhookToken generates a secure random token for webhook URL
func generateWebhookToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreateWebhook creates a new webhook endpoint
func (h *WebhookHandler) CreateWebhook(c *gin.Context) {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workspace ID is required"})
		return
	}

	workspaceObjID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObjID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.CreateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate webhook token
	token, err := generateWebhookToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate webhook token"})
		return
	}

	// Set default config values
	if req.Config.ResponseStatus == 0 {
		req.Config.ResponseStatus = 200
	}
	if req.Config.ContentType == "" {
		req.Config.ContentType = "application/json"
	}
	if req.Config.Methods == nil || len(req.Config.Methods) == 0 {
		req.Config.Methods = []string{"POST"}
	}
	if req.Config.ResponseHeaders == nil {
		req.Config.ResponseHeaders = make(map[string]string)
	}
	req.Config.LogRequests = true // Default to logging

	// Calculate expiry date
	var expiresAt *time.Time
	if req.ExpiresInDays > 0 {
		expiry := time.Now().AddDate(0, 0, req.ExpiresInDays)
		expiresAt = &expiry
	}

	webhook := models.WebhookEndpoint{
		ID:          primitive.NewObjectID(),
		WorkspaceID: workspaceObjID,
		Name:        req.Name,
		Description: req.Description,
		Token:       token,
		Status:      "active",
		CreatedBy:   userObjID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
		Config:      req.Config,
		Stats: models.WebhookStats{
			TotalRequests:    0,
			SuccessRequests:  0,
			FailedRequests:   0,
			AverageResponse:  0,
			RequestsToday:    0,
			RequestsThisWeek: 0,
		},
	}

	// Insert webhook into database
	collection := h.db.Collection("webhooks")
	_, err = collection.InsertOne(context.Background(), webhook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create webhook"})
		return
	}

	// Generate full webhook URL
	baseURL := c.Request.Header.Get("X-Forwarded-Host")
	if baseURL == "" {
		baseURL = c.Request.Host
	}
	scheme := "https"
	if strings.Contains(baseURL, "localhost") || strings.Contains(baseURL, "127.0.0.1") {
		scheme = "http"
	}

	webhookURL := fmt.Sprintf("%s://%s/webhook/%s", scheme, baseURL, token)

	response := models.WebhookResponse{
		WebhookEndpoint: &webhook,
		URL:             webhookURL,
	}

	c.JSON(http.StatusCreated, response)
}

// GetWebhooks returns all webhooks for a workspace
func (h *WebhookHandler) GetWebhooks(c *gin.Context) {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workspace ID is required"})
		return
	}

	workspaceObjID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	skip := (page - 1) * limit

	// Get total count
	collection := h.db.Collection("webhooks")
	filter := bson.M{"workspace_id": workspaceObjID}
	
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count webhooks"})
		return
	}

	// Get webhooks with pagination
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get webhooks"})
		return
	}
	defer cursor.Close(context.Background())

	var webhooks []models.WebhookEndpoint
	if err := cursor.All(context.Background(), &webhooks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode webhooks"})
		return
	}

	// Generate URLs for each webhook
	baseURL := c.Request.Header.Get("X-Forwarded-Host")
	if baseURL == "" {
		baseURL = c.Request.Host
	}
	scheme := "https"
	if strings.Contains(baseURL, "localhost") || strings.Contains(baseURL, "127.0.0.1") {
		scheme = "http"
	}

	webhookResponses := make([]models.WebhookResponse, len(webhooks))
	for i, webhook := range webhooks {
		webhookURL := fmt.Sprintf("%s://%s/webhook/%s", scheme, baseURL, webhook.Token)
		webhookResponses[i] = models.WebhookResponse{
			WebhookEndpoint: &webhook,
			URL:             webhookURL,
		}
	}

	response := models.WebhookListResponse{
		Webhooks: webhookResponses,
		Total:    total,
		Page:     page,
		Limit:    limit,
		HasMore:  int64(skip+limit) < total,
	}

	c.JSON(http.StatusOK, response)
}

// GetWebhook returns a specific webhook
func (h *WebhookHandler) GetWebhook(c *gin.Context) {
	webhookID := c.Param("id")
	if webhookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook ID is required"})
		return
	}

	webhookObjID, err := primitive.ObjectIDFromHex(webhookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook ID"})
		return
	}

	collection := h.db.Collection("webhooks")
	var webhook models.WebhookEndpoint
	err = collection.FindOne(context.Background(), bson.M{"_id": webhookObjID}).Decode(&webhook)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Webhook not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get webhook"})
		return
	}

	// Generate full webhook URL
	baseURL := c.Request.Header.Get("X-Forwarded-Host")
	if baseURL == "" {
		baseURL = c.Request.Host
	}
	scheme := "https"
	if strings.Contains(baseURL, "localhost") || strings.Contains(baseURL, "127.0.0.1") {
		scheme = "http"
	}

	webhookURL := fmt.Sprintf("%s://%s/webhook/%s", scheme, baseURL, webhook.Token)

	response := models.WebhookResponse{
		WebhookEndpoint: &webhook,
		URL:             webhookURL,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateWebhook updates an existing webhook
func (h *WebhookHandler) UpdateWebhook(c *gin.Context) {
	webhookID := c.Param("id")
	if webhookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook ID is required"})
		return
	}

	webhookObjID, err := primitive.ObjectIDFromHex(webhookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook ID"})
		return
	}

	var req models.UpdateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build update document
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	if req.Name != "" {
		update["$set"].(bson.M)["name"] = req.Name
	}
	if req.Description != "" {
		update["$set"].(bson.M)["description"] = req.Description
	}
	if req.Status != "" {
		update["$set"].(bson.M)["status"] = req.Status
	}
	if req.Config.Methods != nil {
		update["$set"].(bson.M)["config"] = req.Config
	}

	collection := h.db.Collection("webhooks")
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": webhookObjID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update webhook"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Webhook not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook updated successfully"})
}

// DeleteWebhook deletes a webhook
func (h *WebhookHandler) DeleteWebhook(c *gin.Context) {
	webhookID := c.Param("id")
	if webhookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook ID is required"})
		return
	}

	webhookObjID, err := primitive.ObjectIDFromHex(webhookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook ID"})
		return
	}

	collection := h.db.Collection("webhooks")
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": webhookObjID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete webhook"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Webhook not found"})
		return
	}

	// Also delete related webhook requests
	requestsCollection := h.db.Collection("webhook_requests")
	_, err = requestsCollection.DeleteMany(context.Background(), bson.M{"webhook_id": webhookObjID})
	if err != nil {
		// Log error but don't fail the webhook deletion
		fmt.Printf("Warning: Failed to delete webhook requests: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook deleted successfully"})
}

// HandleWebhookRequest handles incoming webhook requests
func (h *WebhookHandler) HandleWebhookRequest(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook token is required"})
		return
	}

	startTime := time.Now()

	// Find webhook by token
	collection := h.db.Collection("webhooks")
	var webhook models.WebhookEndpoint
	err := collection.FindOne(context.Background(), bson.M{"token": token}).Decode(&webhook)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Webhook not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find webhook"})
		return
	}

	// Check if webhook is active
	if webhook.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Webhook is inactive"})
		return
	}

	// Check if webhook has expired
	if webhook.ExpiresAt != nil && time.Now().After(*webhook.ExpiresAt) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Webhook has expired"})
		return
	}

	// Check if method is allowed
	methodAllowed := false
	for _, allowedMethod := range webhook.Config.Methods {
		if strings.ToUpper(c.Request.Method) == strings.ToUpper(allowedMethod) {
			methodAllowed = true
			break
		}
	}
	if !methodAllowed {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	// Read request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	body := string(bodyBytes)

	// Extract headers
	headers := make(map[string]string)
	for name, values := range c.Request.Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}

	// Extract query parameters
	query := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			query[key] = values[0]
		}
	}

	// Authenticate if required
	if webhook.Config.Authentication.Enabled {
		if !h.authenticateRequest(c, webhook.Config.Authentication, headers) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			return
		}
	}

	// Apply response delay if configured
	if webhook.Config.ResponseDelay > 0 {
		time.Sleep(time.Duration(webhook.Config.ResponseDelay) * time.Millisecond)
	}

	// Log request if enabled
	if webhook.Config.LogRequests {
		h.logWebhookRequest(webhook.ID, c.Request, body, headers, query, startTime)
	}

	// Update webhook stats
	h.updateWebhookStats(webhook.ID, true, time.Since(startTime))

	// Set CORS headers if enabled
	if webhook.Config.CORSEnabled {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", strings.Join(webhook.Config.Methods, ", "))
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	}

	// Set custom response headers
	for key, value := range webhook.Config.ResponseHeaders {
		c.Header(key, value)
	}

	// Set content type
	c.Header("Content-Type", webhook.Config.ContentType)

	// Send response
	responseBody := webhook.Config.ResponseBody
	if responseBody == "" {
		responseBody = `{"message": "Webhook received successfully", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`
	}

	c.String(webhook.Config.ResponseStatus, responseBody)
}

// authenticateRequest validates authentication based on webhook config
func (h *WebhookHandler) authenticateRequest(c *gin.Context, auth models.WebhookAuth, headers map[string]string) bool {
	switch auth.Type {
	case "basic":
		username, password, ok := c.Request.BasicAuth()
		return ok && username == auth.Username && password == auth.Password
	
	case "bearer":
		authHeader := headers["Authorization"]
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return false
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		return token == auth.Token
	
	case "signature":
		// Implement signature validation if needed
		return true
	
	default:
		return true
	}
}

// logWebhookRequest logs an incoming webhook request
func (h *WebhookHandler) logWebhookRequest(webhookID primitive.ObjectID, req *http.Request, body string, headers, query map[string]string, startTime time.Time) {
	requestLog := models.WebhookRequest{
		ID:         primitive.NewObjectID(),
		WebhookID:  webhookID,
		Method:     req.Method,
		Path:       req.URL.Path,
		Query:      query,
		Headers:    headers,
		Body:       body,
		UserAgent:  req.UserAgent(),
		IP:         req.RemoteAddr,
		Size:       int64(len(body)),
		ResponseTime: time.Since(startTime).Milliseconds(),
		Timestamp:  time.Now(),
		Processed:  true,
	}

	collection := h.db.Collection("webhook_requests")
	_, err := collection.InsertOne(context.Background(), requestLog)
	if err != nil {
		fmt.Printf("Warning: Failed to log webhook request: %v\n", err)
	}
}

// updateWebhookStats updates webhook statistics
func (h *WebhookHandler) updateWebhookStats(webhookID primitive.ObjectID, success bool, responseTime time.Duration) {
	collection := h.db.Collection("webhooks")
	
	now := time.Now()
	update := bson.M{
		"$inc": bson.M{
			"stats.total_requests": 1,
		},
		"$set": bson.M{
			"stats.last_request_at": &now,
		},
	}

	if success {
		update["$inc"].(bson.M)["stats.success_requests"] = 1
	} else {
		update["$inc"].(bson.M)["stats.failed_requests"] = 1
	}

	// Update average response time (simplified calculation)
	update["$set"].(bson.M)["stats.average_response"] = responseTime.Milliseconds()

	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": webhookID}, update)
	if err != nil {
		fmt.Printf("Warning: Failed to update webhook stats: %v\n", err)
	}
}

// GetWebhookRequests returns webhook request logs
func (h *WebhookHandler) GetWebhookRequests(c *gin.Context) {
	webhookID := c.Param("id")
	if webhookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook ID is required"})
		return
	}

	webhookObjID, err := primitive.ObjectIDFromHex(webhookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook ID"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	skip := (page - 1) * limit

	// Get total count
	collection := h.db.Collection("webhook_requests")
	filter := bson.M{"webhook_id": webhookObjID}
	
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count webhook requests"})
		return
	}

	// Get requests with pagination
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get webhook requests"})
		return
	}
	defer cursor.Close(context.Background())

	var requests []models.WebhookRequest
	if err := cursor.All(context.Background(), &requests); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode webhook requests"})
		return
	}

	response := models.WebhookRequestListResponse{
		Requests: requests,
		Total:    total,
		Page:     page,
		Limit:    limit,
		HasMore:  int64(skip+limit) < total,
	}

	c.JSON(http.StatusOK, response)
}