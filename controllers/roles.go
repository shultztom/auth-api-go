package controllers

import (
	"auth-api-go/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Structs
type roleRequest struct {
	Role string `json:"role"`
}

// GetRoles GET /roles
func GetRoles(c *gin.Context) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	tokenHeader := c.GetHeader("x-auth-token")

	token, err := services.ParseToken(tokenHeader, jwtKey)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
		return
	}
	var username = token.Claims.(jwt.MapClaims)["username"]

	roles, err := services.GetRolesByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Roles for User not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Roles": roles})
}

// DoesUserHaveRole GET /roles/<role>
func DoesUserHaveRole(c *gin.Context) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	tokenHeader := c.GetHeader("x-auth-token")

	// Get role from url
	role := c.Param("role")

	// Get user from token
	token, err := services.ParseToken(tokenHeader, jwtKey)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
		return
	}
	var username = token.Claims.(jwt.MapClaims)["username"]

	// Look to see if user already has role
	hasRoleAlready, err := services.RoleCheck(role, username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting roles for User!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hasRoleAlready": hasRoleAlready})
}

// AddRole POST /roles
func AddRole(c *gin.Context) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	tokenHeader := c.GetHeader("x-auth-token")

	// Get user from token
	token, err := services.ParseToken(tokenHeader, jwtKey)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
		return
	}
	var username = token.Claims.(jwt.MapClaims)["username"]

	var newRole roleRequest
	if err := c.BindJSON(&newRole); err != nil {
		return
	}

	// Look to see if user already has role
	hasRoleAlready, err := services.RoleCheck(newRole.Role, username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting roles for User!"})
		return
	}

	if hasRoleAlready {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already has role!"})
		return
	}

	// Add Role
	err = services.AddRole(username.(string), newRole.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Added Role": newRole.Role})
}
