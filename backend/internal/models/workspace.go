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
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string               `bson:"name" json:"name" binding:"required,min=1,max=100"`
	Description string               `bson:"description" json:"description"`
	WorkspaceID primitive.ObjectID   `bson:"workspace_id" json:"workspace_id"`
	CreatedBy   primitive.ObjectID   `bson:"created_by" json:"created_by"`
	Requests    []primitive.ObjectID `bson:"requests" json:"requests"`
	FavoritedBy []primitive.ObjectID `bson:"favorited_by" json:"favorited_by"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updated_at"`
}

type CreateCollectionRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
}

type WorkspaceInvitation struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	WorkspaceID   primitive.ObjectID `bson:"workspace_id" json:"workspace_id"`
	Token         string             `bson:"token" json:"token"`
	Role          string             `bson:"role" json:"role"`
	CreatedBy     primitive.ObjectID `bson:"created_by" json:"created_by"`
	ExpiresAt     *time.Time         `bson:"expires_at" json:"expires_at"`
	Used          bool               `bson:"used" json:"used"`
	UsedBy        *primitive.ObjectID `bson:"used_by,omitempty" json:"used_by,omitempty"`
	UsedAt        *time.Time         `bson:"used_at,omitempty" json:"used_at,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

type CreateInvitationRequest struct {
	Role           string `json:"role" binding:"required"`
	ExpiresInHours int    `json:"expiresInHours"`
}

type InvitationResponse struct {
	Token       string     `json:"token"`
	InviteLink  string     `json:"inviteLink"`
	Role        string     `json:"role"`
	ExpiresAt   *time.Time `json:"expiresAt"`
	WorkspaceID string     `json:"workspaceId"`
}