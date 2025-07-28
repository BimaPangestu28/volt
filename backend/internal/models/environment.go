package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Variable represents an environment variable
type Variable struct {
	Key          string `json:"key" bson:"key"`
	Value        string `json:"value" bson:"value"`
	InitialValue string `json:"initial_value" bson:"initial_value"`
	Type         string `json:"type" bson:"type"` // default, secret
	Enabled      bool   `json:"enabled" bson:"enabled"`
}

// Environment represents an environment configuration
type Environment struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	WorkspaceID string             `json:"workspace_id" bson:"workspace_id"`
	Variables   []Variable         `json:"variables" bson:"variables"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// CreateEnvironmentRequest represents the request to create an environment
type CreateEnvironmentRequest struct {
	Name      string     `json:"name" binding:"required"`
	Variables []Variable `json:"variables"`
}

// UpdateEnvironmentRequest represents the request to update an environment
type UpdateEnvironmentRequest struct {
	Name      string     `json:"name" binding:"required"`
	Variables []Variable `json:"variables"`
}