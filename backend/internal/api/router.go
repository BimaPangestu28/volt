package api

import (
	"os"
	"volt-backend/internal/api/handlers"
	"volt-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// Set Gin mode
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}

	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		origin := os.Getenv("CORS_ORIGIN")
		if origin == "" {
			origin = "http://localhost:5176" // Update to match current frontend port
		}
		
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "volt-api"})
	})

	// API routes
	api := router.Group("/api")
	{
		// Authentication routes
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", handlers.Register)
			authRoutes.POST("/login", handlers.Login)
			authRoutes.POST("/refresh", handlers.RefreshToken)
			authRoutes.GET("/me", auth.AuthMiddleware(), handlers.GetMe)
			authRoutes.DELETE("/logout", auth.AuthMiddleware(), handlers.Logout)
		}

		// Protected routes
		protected := api.Group("/", auth.AuthMiddleware())
		{
			// Workspaces
			workspaceRoutes := protected.Group("/workspaces")
			{
				workspaceRoutes.GET("", handlers.GetWorkspaces)
				workspaceRoutes.POST("", handlers.CreateWorkspace)
			}
			
			// Individual workspace routes (separate to avoid conflicts)
			protected.GET("/workspace-detail/:id", handlers.GetWorkspace)
			protected.PUT("/workspace-detail/:id", handlers.UpdateWorkspace)
			protected.DELETE("/workspace-detail/:id", handlers.DeleteWorkspace)

			// Collections - workspace specific
			protected.GET("/workspace-collections/:workspaceId", handlers.GetCollections)
			protected.POST("/workspace-collections/:workspaceId", handlers.CreateCollection)

			// Collections - individual operations
			collectionRoutes := protected.Group("/collections")
			{
				collectionRoutes.GET("/:id", handlers.GetCollection)
				collectionRoutes.PUT("/:id", handlers.UpdateCollection)
				collectionRoutes.DELETE("/:id", handlers.DeleteCollection)
			}

			// Requests - collection specific
			protected.GET("/collection-requests/:collectionId", handlers.GetRequests)
			protected.POST("/collection-requests/:collectionId", handlers.CreateRequest)

			// Requests - individual operations
			requestRoutes := protected.Group("/requests")
			{
				requestRoutes.GET("/:id", handlers.GetRequest)
				requestRoutes.PUT("/:id", handlers.UpdateRequest)
				requestRoutes.DELETE("/:id", handlers.DeleteRequest)
				requestRoutes.POST("/:id/execute", handlers.ExecuteRequest)
			}

			// Execute raw request without saving
			protected.POST("/execute-request", handlers.ExecuteRawRequest)
		}
	}

	return router
}