package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Try to get token from cookie first (more secure)
		if cookie, err := c.Cookie("volt_token"); err == nil && cookie != "" {
			tokenString = cookie
		} else {
			// Fallback to Authorization header for API clients
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
				c.Abort()
				return
			}

			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
				c.Abort()
				return
			}
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if objectID, ok := userID.(primitive.ObjectID); ok {
			return objectID.Hex()
		}
	}
	return ""
}

func GetUserObjectID(c *gin.Context) primitive.ObjectID {
	if userID, exists := c.Get("user_id"); exists {
		if objectID, ok := userID.(primitive.ObjectID); ok {
			return objectID
		}
	}
	return primitive.NilObjectID
}