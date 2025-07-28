package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"volt-backend/internal/auth"
)

func GetNotifications(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Mock notifications data for now
	notifications := []gin.H{
		{
			"id":          "1",
			"title":       "Welcome to Volt!",
			"message":     "Your account has been successfully created. Start building your API collections.",
			"type":        "welcome",
			"read":        false,
			"timestamp":   time.Now().Add(-1 * time.Hour).Unix(),
			"action_url":  "/dashboard",
		},
		{
			"id":          "2",
			"title":       "New Team Member",
			"message":     "John Doe has joined your workspace 'E-commerce API'",
			"type":        "team",
			"read":        false,
			"timestamp":   time.Now().Add(-2 * time.Hour).Unix(),
			"action_url":  "/workspaces/1",
		},
		{
			"id":          "3",
			"title":       "API Request Failed",
			"message":     "Your scheduled request 'Health Check' failed with status 500",
			"type":        "error",
			"read":        true,
			"timestamp":   time.Now().Add(-4 * time.Hour).Unix(),
			"action_url":  "/requests/123",
		},
	}

	// Count unread notifications
	unreadCount := 0
	for _, notification := range notifications {
		if read, ok := notification["read"].(bool); ok && !read {
			unreadCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"total":         len(notifications),
		"unread_count":  unreadCount,
	})
}

func MarkNotificationAsRead(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	notificationID := c.Param("id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notification ID is required"})
		return
	}

	// TODO: Implement actual database update
	// For now, just return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Notification marked as read",
		"id":      notificationID,
	})
}

func MarkAllNotificationsAsRead(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// TODO: Implement actual database update
	// For now, just return success
	c.JSON(http.StatusOK, gin.H{
		"message": "All notifications marked as read",
	})
}