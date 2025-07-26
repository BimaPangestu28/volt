package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HTTPMethod string

const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
	PATCH  HTTPMethod = "PATCH"
)

type AuthType string

const (
	Bearer AuthType = "bearer"
	Basic  AuthType = "basic"
	APIKey AuthType = "apikey"
	None   AuthType = "none"
)

type BodyType string

const (
	JSON   BodyType = "json"
	Form   BodyType = "form"
	Raw    BodyType = "raw"
	Binary BodyType = "binary"
)

type Request struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `bson:"name" json:"name" binding:"required,min=1,max=200"`
	Method       HTTPMethod         `bson:"method" json:"method" binding:"required"`
	URL          string             `bson:"url" json:"url" binding:"required"`
	Headers      map[string]string  `bson:"headers" json:"headers"`
	Body         RequestBody        `bson:"body" json:"body"`
	Auth         Auth               `bson:"auth" json:"auth"`
	Tests        []string           `bson:"tests" json:"tests"`
	CollectionID primitive.ObjectID `bson:"collection_id" json:"collection_id"`
	CreatedBy    primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type RequestBody struct {
	Type    BodyType `bson:"type" json:"type"`
	Content string   `bson:"content" json:"content"`
}

type Auth struct {
	Type        AuthType          `bson:"type" json:"type"`
	Credentials map[string]string `bson:"credentials" json:"credentials"`
}

type CreateRequestRequest struct {
	Name    string            `json:"name" binding:"required,min=1,max=200"`
	Method  HTTPMethod        `json:"method" binding:"required"`
	URL     string            `json:"url" binding:"required"`
	Headers map[string]string `json:"headers"`
	Body    RequestBody       `json:"body"`
	Auth    Auth              `json:"auth"`
	Tests   []string          `json:"tests"`
}

type ExecuteRequestResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Time    int64             `json:"time"` // milliseconds
}