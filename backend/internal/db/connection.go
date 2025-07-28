package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func Connect() error {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017/volt"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Configure client options for external MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	
	// Set additional options for external connection
	clientOptions.SetMaxPoolSize(10)
	clientOptions.SetMinPoolSize(5)
	clientOptions.SetMaxConnIdleTime(30 * time.Second)
	clientOptions.SetServerSelectionTimeout(10 * time.Second)
	clientOptions.SetSocketTimeout(30 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("❌ Failed to connect to MongoDB: %v", err)
		return err
	}

	// Test the connection with retry
	for i := 0; i < 3; i++ {
		if err = client.Ping(ctx, nil); err != nil {
			log.Printf("⚠️ MongoDB ping attempt %d failed: %v", i+1, err)
			if i == 2 {
				return err
			}
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	Client = client
	Database = client.Database("volt")
	
	log.Println("🔗 Connected to external MongoDB successfully")
	return nil
}

func Disconnect() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		if err := Client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		} else {
			log.Println("📴 Disconnected from MongoDB")
		}
	}
}

func GetDB() *mongo.Database {
	return Database
}

func GetCollection(name string) *mongo.Collection {
	return Database.Collection(name)
}