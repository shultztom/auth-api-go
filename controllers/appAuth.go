package controllers

import (
	"auth-api-go/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AppClaims struct {
	AppName string `json:"appName"`
	jwt.StandardClaims
}

// AppVerify GET /app/verify
func AppVerify(c *gin.Context) {
	appJwtKey := []byte(os.Getenv("JWT_APP_SECRET"))
	tokenHeader := c.GetHeader("X-API-Token")

	isValid, err := services.VerifyToken(tokenHeader, appJwtKey)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
		return
	}

	if isValid {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
	}
}

// AppDeleteUser Delete /app/user
func AppDeleteUser(c *gin.Context) {
	appJwtKey := []byte(os.Getenv("JWT_APP_SECRET"))
	tokenHeader := c.GetHeader("X-API-Token")

	isValid, err := services.VerifyToken(tokenHeader, appJwtKey)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
		return
	}

	if !isValid {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
		return
	}

	// Get user from request
	username := c.Param("username")

	// Delete active sessions, if any
	_, err = services.DeleteSessionInRedis(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	err = services.DeleteUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Deleted user": username})
}
