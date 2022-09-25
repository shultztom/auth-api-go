package controllers

import (
	"auth-api-go/models"
	"auth-api-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
)

// Uses some global vars from users.go

// Structs
type roleRequest struct {
	Role string `json:"role"`
}

// GetRoles GET /roles
func GetRoles(c *gin.Context) {
	token := utils.ParseToken(c, jwtKey)
	var username = token.Claims.(jwt.MapClaims)["username"]

	var roles []models.Roles
	result := models.DB.Find(&roles, "username = ?", username)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Roles for User not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Roles": roles})
}

// AddRole POST /roles
func AddRole(c *gin.Context) {
	// Get user from token
	token := utils.ParseToken(c, jwtKey)
	var username = token.Claims.(jwt.MapClaims)["username"]

	var newRole roleRequest
	if err := c.BindJSON(&newRole); err != nil {
		return
	}

	// Look to see if user already has role
	var roles []models.Roles
	result := models.DB.Find(&roles, "username = ?", username)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Roles for User not found!"})
		return
	}
	var hasRoleAlready = false
	for i := 0; i < len(roles); i++ {
		if roles[i].Role == newRole.Role {
			hasRoleAlready = true
			break
		}
	}
	if hasRoleAlready {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already has role!"})
		return
	}

	// Add Role
	roleEntry := &models.Roles{
		Username: username.(string),
		Role:     newRole.Role,
	}

	_, err := models.DB.Create(roleEntry).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Added Role": newRole.Role})
}
