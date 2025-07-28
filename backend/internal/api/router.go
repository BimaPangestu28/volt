package api

import (
	"os"
	"volt-backend/internal/api/handlers"
	"volt-backend/internal/auth"
	"volt-backend/internal/db"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// Set Gin mode
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}

	router := gin.Default()

	// Initialize handlers with database connection
	webhookHandler := handlers.NewWebhookHandler(db.GetDB())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		origin := os.Getenv("CORS_ORIGIN")
		if origin == "" {
			origin = "http://localhost:5173" // Default frontend port
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
			
			// Workspace invitations
			protected.POST("/workspaces/:id/invite-link", handlers.GenerateInvitationLink)
			protected.POST("/invitations/:token/accept", handlers.AcceptInvitation)

			// Collections - workspace specific
			protected.GET("/workspace-collections/:workspaceId", handlers.GetCollections)
			protected.POST("/workspace-collections/:workspaceId", handlers.CreateCollection)

			// Collections - individual operations
			collectionRoutes := protected.Group("/collections")
			{
				collectionRoutes.GET("/:id", handlers.GetCollection)
				collectionRoutes.PUT("/:id", handlers.UpdateCollection)
				collectionRoutes.DELETE("/:id", handlers.DeleteCollection)
				collectionRoutes.POST("/:id/favorite", handlers.ToggleCollectionFavorite)
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

			// Environments - workspace specific
			protected.GET("/workspace-environments/:workspaceId", handlers.GetEnvironments)
			protected.POST("/workspace-environments/:workspaceId", handlers.CreateEnvironment)

			// Environments - individual operations
			environmentRoutes := protected.Group("/environments")
			{
				environmentRoutes.GET("/:id", handlers.GetEnvironment)
				environmentRoutes.PUT("/:id", handlers.UpdateEnvironment)
				environmentRoutes.DELETE("/:id", handlers.DeleteEnvironment)
			}

			// Webhooks - workspace specific
			protected.GET("/workspace-webhooks/:workspaceId", webhookHandler.GetWebhooks)
			protected.POST("/workspace-webhooks/:workspaceId", webhookHandler.CreateWebhook)

			// Webhooks - individual operations
			webhookRoutes := protected.Group("/webhooks")
			{
				webhookRoutes.GET("/:id", webhookHandler.GetWebhook)
				webhookRoutes.PUT("/:id", webhookHandler.UpdateWebhook)
				webhookRoutes.DELETE("/:id", webhookHandler.DeleteWebhook)
				webhookRoutes.GET("/:id/requests", webhookHandler.GetWebhookRequests)
			}

			// Dashboard routes
			dashboardRoutes := protected.Group("/dashboard")
			{
				dashboardRoutes.GET("/stats", handlers.GetDashboardStats)
				dashboardRoutes.GET("/api-health", handlers.GetAPIHealth)
				dashboardRoutes.GET("/activity", handlers.GetDashboardActivity)
				dashboardRoutes.GET("/team", handlers.GetDashboardTeam)
			}

			// Notifications routes
			notificationRoutes := protected.Group("/notifications")
			{
				notificationRoutes.GET("", handlers.GetNotifications)
				notificationRoutes.PUT("/:id/read", handlers.MarkNotificationAsRead)
				notificationRoutes.PUT("/read-all", handlers.MarkAllNotificationsAsRead)
			}

			// Workspace search - TODO: implement SearchWorkspace handler
			// protected.GET("/workspaces/:id/search", handlers.SearchWorkspace)
		}
	}

	// Public webhook endpoint (no authentication required)
	router.Any("/webhook/:token", webhookHandler.HandleWebhookRequest)

	return router
}