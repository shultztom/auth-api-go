package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var REDIS *redis.Client

func ConnectRedis() {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     "10.0.0.98:6379",
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
