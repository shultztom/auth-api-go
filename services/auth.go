package services

import (
	"auth-api-go/redis"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateToken(username string) (string, error) {
	ctx := context.Background()

	// See if token exists in redis
	val, err := redis.REDIS.Get(ctx, username+"-token").Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			fmt.Println("Token not found in redis; will make new one")
		} else {
			fmt.Println("error with redis get", err.Error())
			return "", fmt.Errorf("error with redis get: %v", err)
		}
	}

	if val != "" {
		fmt.Println("token found in redis; returning it")
		return val, nil
	}

	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	expirationTime := time.Now().Add(8 * time.Hour)

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("error with creating token: %v", err)
	}

	// Save as session in redis
	now := time.Now().Add(-1 * time.Minute) // Shorter than expiration time to account for latency
	duration := expirationTime.Sub(now)
	err = redis.REDIS.Set(ctx, username+"-token", tokenString, duration).Err()
	if err != nil {
		fmt.Println("error with redis set", err.Error())
		return "", fmt.Errorf("error with redis set: %v", err)
	}

	return tokenString, nil
}

func ParseToken(tokenHeader string, jwtKey []byte) (*jwt.Token, error) {
	if tokenHeader == "" {
		return nil, errors.New("missing token")
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired())

	if err != nil {
		return nil, errors.New("unable to parse token")
	}

	// Verify session exists
	username, ok := token.Claims.(jwt.MapClaims)["username"].(string)
	if !ok || username == "" {
		return nil, errors.New("unable to parse token")
	}

	ctx := context.Background()
	val, err := redis.REDIS.Get(ctx, username+"-token").Result()
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

func VerifyToken(tokenHeader string, jwtKey []byte) (bool, error) {
	token, err := ParseToken(tokenHeader, jwtKey)
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}

func GetUsernameFromToken(tokenHeader string, jwtKey []byte) (string, error) {
	token, err := ParseToken(tokenHeader, jwtKey)
	if err != nil {
		return "", err
	}
	username, ok := token.Claims.(jwt.MapClaims)["username"].(string)
	if !ok {
		return "", errors.New("unable to parse token")
	}
	return username, nil
}
