package services

import (
	"auth-api-go/redis"
	"context"
	"fmt"
)

func DeleteSessionInRedis(username string) (bool, error) {
	ctx := context.Background()
	err := redis.REDIS.Del(ctx, username+"-token").Err()
	if err != nil {
		fmt.Println("error with redis del", err.Error())
		return false, fmt.Errorf("error with redis del: %v", err)
	}
	return true, nil
}
