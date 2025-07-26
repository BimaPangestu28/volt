package main

import (
	"log"
	"os"

	"volt-backend/internal/api"
	"volt-backend/internal/db"
	"volt-backend/internal/services/redis"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to database
	if err := db.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Disconnect()

	// Connect to Redis (optional - for caching and sessions)
	if err := redis.Connect(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Println("⚠️ Continuing without Redis - some features may be limited")
	} else {
		defer redis.Disconnect()
	}

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Setup and start API server
	router := api.SetupRouter()
	log.Printf("🚀 Volt API server starting on port %s", port)
	
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}