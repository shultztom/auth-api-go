package utils

import (
	"auth-api-go/redis"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func ParseToken(c *gin.Context, jwtKey []byte, headerName string) (*jwt.Token, error) {
	tokenHeader := c.GetHeader(headerName)

	if tokenHeader == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Missing Token!"})
		return nil, errors.New("missing token")
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unable to Parse Token!"})
		return nil, errors.New("forbidden")
	}

	// Verify session exists
	var username = token.Claims.(jwt.MapClaims)["username"]

	ctx := context.Background()
	val, err := redis.REDIS.Get(ctx, username.(string)+"-token").Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			// Token not found in redis; throw forbidden error
			fmt.Println("Token not found in redis; returning forbidden")
			return nil, errors.New("forbidden")
		} else {
			fmt.Println("error with redis get", err.Error())
			return nil, fmt.Errorf("error with redis get: %v", err)
		}
	}

	// Ensure token passed is what is in session; throw error if not
	if val != tokenHeader {
		return nil, errors.New("forbidden")
	}

	return token, nil
}
