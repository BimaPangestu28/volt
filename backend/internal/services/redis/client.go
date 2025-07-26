package redis

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

func Connect() error {
	redisURI := os.Getenv("REDIS_URI")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	
	if redisURI == "" {
		redisURI = "redis://localhost:6379"
	}

	// Parse Redis URI to get host and port
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		log.Printf("❌ Failed to parse Redis URI: %v", err)
		return err
	}

	// Set password if provided
	if redisPassword != "" {
		opt.Password = redisPassword
	}

	// Configure additional options for external Redis
	opt.PoolSize = 10
	opt.MinIdleConns = 5
	opt.DialTimeout = 10 * time.Second
	opt.ReadTimeout = 30 * time.Second
	opt.WriteTimeout = 30 * time.Second
	opt.PoolTimeout = 30 * time.Second

	Client = redis.NewClient(opt)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < 3; i++ {
		_, err = Client.Ping(ctx).Result()
		if err != nil {
			log.Printf("⚠️ Redis ping attempt %d failed: %v", i+1, err)
			if i == 2 {
				return err
			}
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	log.Println("🔗 Connected to external Redis successfully")
	return nil
}

func Disconnect() {
	if Client != nil {
		if err := Client.Close(); err != nil {
			log.Printf("Error disconnecting from Redis: %v", err)
		} else {
			log.Println("📴 Disconnected from Redis")
		}
	}
}

func GetClient() *redis.Client {
	return Client
}