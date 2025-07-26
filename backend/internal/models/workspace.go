package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workspace struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string               `bson:"name" json:"name" binding:"required,min=1,max=100"`
	Description string               `bson:"description" json:"description"`
	OwnerID     primitive.ObjectID   `bson:"owner_id" json:"owner_id"`
	Members     []primitive.ObjectID `bson:"members" json:"members"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updated_at"`
}

type CreateWorkspaceRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
}

type Collection struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name" binding:"required,min=1,max=100"`
	Description string             `bson:"description" json:"description"`
	WorkspaceID primitive.ObjectID `bson:"workspace_id" json:"workspace_id"`
	CreatedBy   primitive.ObjectID `bson:"created_by" json:"created_by"`
	Requests    []primitive.ObjectID `bson:"requests" json:"requests"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateCollectionRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
}