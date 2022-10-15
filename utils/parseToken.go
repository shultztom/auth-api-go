package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
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

	return token, nil
}
