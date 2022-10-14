package controllers

import (
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
	token := utils.ParseToken(c, appJwtKey, "X-API-Token")

	isValid := token.Valid

	if isValid {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
	}
}
