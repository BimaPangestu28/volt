package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"volt-backend/internal/auth"
	"volt-backend/internal/db"
)

func GetDashboardStats(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Count workspaces
	workspaceCount, _ := db.GetCollection("workspaces").CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": bson.M{"$elemMatch": bson.M{"user_id": userID}}},
		},
	})

	// Count collections
	collectionCount, _ := db.GetCollection("collections").CountDocuments(ctx, bson.M{
		"owner_id": userID,
	})

	// Count requests
	requestCount, _ := db.GetCollection("requests").CountDocuments(ctx, bson.M{
		"owner_id": userID,
	})

	// Count environments
	environmentCount, _ := db.GetCollection("environments").CountDocuments(ctx, bson.M{
		"owner_id": userID,
	})

	c.JSON(http.StatusOK, gin.H{
		"workspaces":   workspaceCount,
		"collections":  collectionCount,
		"requests":     requestCount,
		"environments": environmentCount,
		"activity_today": 0, // TODO: implement activity tracking
	})
}

func GetAPIHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check MongoDB connection
	mongoHealthy := true
	if err := db.Client.Ping(ctx, nil); err != nil {
		mongoHealthy = false
	}

	// Overall health status
	healthy := mongoHealthy

	status := "healthy"
	if !healthy {
		status = "unhealthy"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"services": gin.H{
			"mongodb": map[string]interface{}{
				"status": func() string {
					if mongoHealthy {
						return "healthy"
					}
					return "unhealthy"
				}(),
				"response_time": "< 5ms", // Mock response time
			},
			"api": gin.H{
				"status":        "healthy",
				"response_time": "< 1ms",
			},
		},
		"timestamp": time.Now().Unix(),
	})
}

func GetDashboardActivity(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Mock activity data for now
	activities := []gin.H{
		{
			"id":          "1",
			"type":        "request_created",
			"description": "Created new API request 'Get Users'",
			"timestamp":   time.Now().Add(-2 * time.Hour).Unix(),
			"user":        "You",
		},
		{
			"id":          "2",
			"type":        "collection_updated",
			"description": "Updated collection 'User Management'",
			"timestamp":   time.Now().Add(-4 * time.Hour).Unix(),
			"user":        "You",
		},
		{
			"id":          "3",
			"type":        "workspace_created",
			"description": "Created workspace 'E-commerce API'",
			"timestamp":   time.Now().Add(-1 * 24 * time.Hour).Unix(),
			"user":        "You",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"activities": activities,
		"total":      len(activities),
	})
}

func GetDashboardTeam(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get workspaces where user is owner or member
	cursor, err := db.GetCollection("workspaces").Find(ctx, bson.M{
		"$or": []bson.M{
			{"owner_id": userID},
			{"members": bson.M{"$elemMatch": bson.M{"user_id": userID}}},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team data"})
		return
	}
	defer cursor.Close(ctx)

	var workspaces []bson.M
	if err = cursor.All(ctx, &workspaces); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse team data"})
		return
	}

	// Extract unique team members
	memberMap := make(map[string]gin.H)
	totalWorkspaces := len(workspaces)

	for _, workspace := range workspaces {
		// Add owner
		if ownerID, ok := workspace["owner_id"].(string); ok && ownerID != "" {
			memberMap[ownerID] = gin.H{
				"id":     ownerID,
				"name":   "Workspace Owner", // TODO: get actual user name
				"email":  "owner@example.com", // TODO: get actual email
				"role":   "owner",
				"avatar": "",
				"status": "online",
			}
		}

		// Add members
		if members, ok := workspace["members"].([]interface{}); ok {
			for _, member := range members {
				if memberDoc, ok := member.(bson.M); ok {
					if memberID, ok := memberDoc["user_id"].(string); ok {
						memberMap[memberID] = gin.H{
							"id":     memberID,
							"name":   "Team Member", // TODO: get actual user name
							"email":  "member@example.com", // TODO: get actual email
							"role":   "member",
							"avatar": "",
							"status": "offline",
						}
					}
				}
			}
		}
	}

	// Convert map to slice
	var teamMembers []gin.H
	for _, member := range memberMap {
		teamMembers = append(teamMembers, member)
	}

	c.JSON(http.StatusOK, gin.H{
		"team_members":    teamMembers,
		"total_members":   len(teamMembers),
		"total_workspaces": totalWorkspaces,
		"active_members":  1, // Mock data
	})
}