package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WebhookEndpoint represents a webhook endpoint for testing
type WebhookEndpoint struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	WorkspaceID primitive.ObjectID `bson:"workspace_id" json:"workspace_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	URL         string             `bson:"url" json:"url"`
	Token       string             `bson:"token" json:"token"`
	Status      string             `bson:"status" json:"status"` // active, inactive
	CreatedBy   primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	ExpiresAt   *time.Time         `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	
	// Configuration
	Config WebhookConfig `bson:"config" json:"config"`
	
	// Statistics
	Stats WebhookStats `bson:"stats" json:"stats"`
}

// WebhookConfig holds configuration for the webhook endpoint
type WebhookConfig struct {
	Methods          []string          `bson:"methods" json:"methods"`                     // GET, POST, PUT, etc.
	ResponseStatus   int               `bson:"response_status" json:"response_status"`     // Default response status
	ResponseHeaders  map[string]string `bson:"response_headers" json:"response_headers"`   // Custom response headers
	ResponseBody     string            `bson:"response_body" json:"response_body"`         // Custom response body
	ResponseDelay    int               `bson:"response_delay" json:"response_delay"`       // Delay in milliseconds
	ContentType      string            `bson:"content_type" json:"content_type"`           // Response content type
	CORSEnabled      bool              `bson:"cors_enabled" json:"cors_enabled"`           // Enable CORS
	LogRequests      bool              `bson:"log_requests" json:"log_requests"`           // Log incoming requests
	MaxRequests      int               `bson:"max_requests" json:"max_requests"`           // Rate limiting
	Authentication   WebhookAuth       `bson:"authentication" json:"authentication"`       // Auth requirements
}

// WebhookAuth represents authentication configuration
type WebhookAuth struct {
	Enabled   bool              `bson:"enabled" json:"enabled"`
	Type      string            `bson:"type" json:"type"`           // none, basic, bearer, signature
	Username  string            `bson:"username,omitempty" json:"username,omitempty"`
	Password  string            `bson:"password,omitempty" json:"password,omitempty"`
	Token     string            `bson:"token,omitempty" json:"token,omitempty"`
	Secret    string            `bson:"secret,omitempty" json:"secret,omitempty"`    // For signature validation
	Headers   map[string]string `bson:"headers,omitempty" json:"headers,omitempty"` // Required headers
}

// WebhookStats holds statistics for the webhook endpoint
type WebhookStats struct {
	TotalRequests    int64     `bson:"total_requests" json:"total_requests"`
	SuccessRequests  int64     `bson:"success_requests" json:"success_requests"`
	FailedRequests   int64     `bson:"failed_requests" json:"failed_requests"`
	LastRequestAt    *time.Time `bson:"last_request_at,omitempty" json:"last_request_at,omitempty"`
	AverageResponse  int64     `bson:"average_response" json:"average_response"` // in milliseconds
	RequestsToday    int64     `bson:"requests_today" json:"requests_today"`
	RequestsThisWeek int64     `bson:"requests_this_week" json:"requests_this_week"`
}

// WebhookRequest represents an incoming webhook request
type WebhookRequest struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	WebhookID  primitive.ObjectID `bson:"webhook_id" json:"webhook_id"`
	Method     string             `bson:"method" json:"method"`
	Path       string             `bson:"path" json:"path"`
	Query      map[string]string  `bson:"query" json:"query"`
	Headers    map[string]string  `bson:"headers" json:"headers"`
	Body       string             `bson:"body" json:"body"`
	UserAgent  string             `bson:"user_agent" json:"user_agent"`
	IP         string             `bson:"ip" json:"ip"`
	Size       int64              `bson:"size" json:"size"`
	
	// Response data
	ResponseStatus int               `bson:"response_status" json:"response_status"`
	ResponseTime   int64             `bson:"response_time" json:"response_time"` // milliseconds
	ResponseSize   int64             `bson:"response_size" json:"response_size"`
	
	// Metadata
	Timestamp  time.Time `bson:"timestamp" json:"timestamp"`
	Processed  bool      `bson:"processed" json:"processed"`
	Error      string    `bson:"error,omitempty" json:"error,omitempty"`
}

// CreateWebhookRequest represents the request to create a webhook
type CreateWebhookRequest struct {
	Name         string        `json:"name" validate:"required,min=1,max=100"`
	Description  string        `json:"description,omitempty" validate:"max=500"`
	ExpiresInDays int          `json:"expires_in_days,omitempty" validate:"min=1,max=365"`
	Config       WebhookConfig `json:"config"`
}

// UpdateWebhookRequest represents the request to update a webhook
type UpdateWebhookRequest struct {
	Name        string        `json:"name,omitempty" validate:"min=1,max=100"`
	Description string        `json:"description,omitempty" validate:"max=500"`
	Status      string        `json:"status,omitempty" validate:"oneof=active inactive"`
	Config      WebhookConfig `json:"config,omitempty"`
}

// WebhookResponse represents the response when creating/getting a webhook
type WebhookResponse struct {
	*WebhookEndpoint
	URL string `json:"url"` // Full webhook URL
}

// WebhookTestRequest represents a request to test a webhook
type WebhookTestRequest struct {
	Method  string            `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
	Query   map[string]string `json:"query,omitempty"`
}

// WebhookTestResponse represents the response from testing a webhook
type WebhookTestResponse struct {
	Success      bool              `json:"success"`
	Status       int               `json:"status"`
	Headers      map[string]string `json:"headers"`
	Body         string            `json:"body"`
	ResponseTime int64             `json:"response_time"` // milliseconds
	Error        string            `json:"error,omitempty"`
	Timestamp    time.Time         `json:"timestamp"`
}

// WebhookListResponse represents a paginated list of webhooks
type WebhookListResponse struct {
	Webhooks []WebhookResponse `json:"webhooks"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	HasMore  bool              `json:"has_more"`
}

// WebhookRequestListResponse represents a paginated list of webhook requests
type WebhookRequestListResponse struct {
	Requests []WebhookRequest `json:"requests"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
	HasMore  bool             `json:"has_more"`
}