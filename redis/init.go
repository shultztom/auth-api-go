package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var REDIS *redis.Client

func ConnectRedis() {
	ctx := context.Background()

	// Load env vars
	err := godotenv.Load()

	isCloud := os.Getenv("IS_CLOUD")

	if err != nil {
		// Ignore error if running in k8s
		if isCloud != "true" {
			log.Fatal("Error loading .env file")
			return
		}
	}

	redisUrl := os.Getenv("REDIS_URL")

	client := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	maxRetries := 3
	retryDelay := 15 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		pong, err := client.Ping(ctx).Result()
		if err != nil {
			fmt.Printf("Attempt %d/%d: Could not connect to Redis: %v\n", attempt, maxRetries, err)
			if attempt < maxRetries {
				fmt.Printf("Retrying in %v...\n", retryDelay)
				time.Sleep(retryDelay)
			} else {
				log.Fatal("Failed to connect to Redis after all retries. Exiting application.")
			}
		} else {
			fmt.Println("Connected to Redis:", pong)
			REDIS = client
			return
		}
	}
}
