package controllers

import (
	"auth-api-go/models"
	"auth-api-go/redis"
	"auth-api-go/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Structs
type userRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Utils
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
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
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
	now := time.Now().Add(-5 * time.Minute) // Shorter than expiration time to account for latency
	duration := expirationTime.Sub(now)
	err = redis.REDIS.Set(ctx, username+"-token", tokenString, duration).Err()
	if err != nil {
		fmt.Println("error with redis set", err.Error())
		return "", fmt.Errorf("error with redis set: %v", err)
	}

	return tokenString, nil
}

func DeleteSessionInRedis(username string) (bool, error) {
	// Delete sessions in redis
	ctx := context.Background()
	err := redis.REDIS.Del(ctx, username+"-token").Err()
	if err != nil {
		fmt.Println("error with redis del", err.Error())
		return false, fmt.Errorf("error with redis del: %v", err)
	}
	return true, nil
}

// Register POST /register
func Register(c *gin.Context) {
	var newUser userRequest

	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	hash, _ := HashPassword(newUser.Password)

	userEntry := &models.User{
		Username: newUser.Username,
		Hash:     hash,
	}

	err := models.DB.Create(userEntry).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	token, err := CreateToken(userEntry.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})

}

// Login POST /login
func Login(c *gin.Context) {
	var userReq userRequest

	if err := c.BindJSON(&userReq); err != nil {
		return
	}

	var user models.User

	if err := models.DB.Where("username = ?", userReq.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not found!"})
		return
	}

	isMatch := CheckPasswordHash(userReq.Password, user.Hash)
	if isMatch {
		token, err := CreateToken(userReq.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Password is incorrect!"})
	}
}

// Verify GET /verify
func Verify(c *gin.Context) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := utils.ParseToken(c, jwtKey, "x-auth-token")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
		return
	}

	isValid := token.Valid

	if isValid {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
	}
}

// DeleteUser DELETE /
func DeleteUser(c *gin.Context) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	// Get user from token
	token, err := utils.ParseToken(c, jwtKey, "x-auth-token")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
		return
	}
	var username = token.Claims.(jwt.MapClaims)["username"]

	// Delete active sessions, if any
	_, err = DeleteSessionInRedis(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var user models.User

	result := models.DB.Where("username = ?", username).Delete(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Deleted user": username})
}

// DeleteUser DELETE /session
func DeleteUserSession(c *gin.Context) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	// Get user from token
	token, err := utils.ParseToken(c, jwtKey, "x-auth-token")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
		return
	}
	var username = token.Claims.(jwt.MapClaims)["username"]

	_, err = DeleteSessionInRedis(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Deleted session for user": username})
}
