package controllers

import (
	"auth-api-go/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Structs
type userRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register POST /register
func Register(c *gin.Context) {
	var newUser userRequest

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	userEntry, err := services.CreateUser(newUser.Username, newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	token, err := services.CreateToken(userEntry.Username)
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

	user, err := services.GetUserByUsername(userReq.Username)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not found!"})
		return
	}

	isMatch := services.CheckPasswordHash(userReq.Password, user.Hash)
	if isMatch {
		token, err := services.CreateToken(userReq.Username)
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
	tokenHeader := c.GetHeader("x-auth-token")

	isValid, err := services.VerifyToken(tokenHeader, jwtKey)
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

// DeleteUser DELETE /
func DeleteUser(c *gin.Context) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	tokenHeader := c.GetHeader("x-auth-token")

	token, err := services.ParseToken(tokenHeader, jwtKey)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
		return
	}
	var username = token.Claims.(jwt.MapClaims)["username"]

	// Delete active sessions, if any
	_, err = services.DeleteSessionInRedis(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	err = services.DeleteUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Deleted user": username})
}

// DeleteUserSession DeleteUser DELETE /session
func DeleteUserSession(c *gin.Context) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	tokenHeader := c.GetHeader("x-auth-token")

	token, err := services.ParseToken(tokenHeader, jwtKey)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
		return
	}
	var username = token.Claims.(jwt.MapClaims)["username"]

	_, err = services.DeleteSessionInRedis(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Deleted session for user": username})
}
