package redis

import (
	"context"
	"fmt"
	"log"
	"os"

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

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Could not connect to Redis:", err)
	} else {
		fmt.Println("Connected to Redis:", pong)
		REDIS = client
	}
}
