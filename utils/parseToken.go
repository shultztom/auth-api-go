package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
)

func ParseToken(c *gin.Context, jwtKey []byte) *jwt.Token {
	tokenHeader := c.GetHeader("x-auth-token")

	if tokenHeader == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Missing Token!"})
		return nil
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unable to Parse Token!"})
		return nil
	}

	return token
}
