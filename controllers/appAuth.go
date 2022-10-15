package controllers

import (
	"auth-api-go/models"
	"auth-api-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
)

type AppClaims struct {
	AppName string `json:"appName"`
	jwt.StandardClaims
}

// AppVerify GET /app/verify
func AppVerify(c *gin.Context) {
	appJwtKey := []byte(os.Getenv("JWT_APP_SECRET"))
	token, err := utils.ParseToken(c, appJwtKey, "X-API-Token")

	isValid := token.Valid

	if isValid && err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
	}
}

// AppDeleteUser Delete /app/user
func AppDeleteUser(c *gin.Context) {
	appJwtKey := []byte(os.Getenv("JWT_APP_SECRET"))
	token, err := utils.ParseToken(c, appJwtKey, "X-API-Token")
	if err != nil {
		return
	}

	if !token.Valid {
		return
	}

	// Get user from request
	username := c.Param("username")

	var user models.User

	result := models.DB.Where("username = ?", username).Delete(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Deleted user": username})
}
