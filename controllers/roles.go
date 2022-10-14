package controllers

import (
	"auth-api-go/models"
	"auth-api-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
)

// Uses some global vars from users.go

// Structs
type roleRequest struct {
	Role string `json:"role"`
}

// GetRoles GET /roles
func GetRoles(c *gin.Context) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	token := utils.ParseToken(c, jwtKey, "x-auth-token")
	var username = token.Claims.(jwt.MapClaims)["username"]

	var roles []models.Roles
	result := models.DB.Find(&roles, "username = ?", username)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Roles for User not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Roles": roles})
}

// DoesUserHaveRole GET /roles/<role>
func DoesUserHaveRole(c *gin.Context) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	// Get role from url
	role := c.Param("role")

	// Get user from token
	token := utils.ParseToken(c, jwtKey, "x-auth-token")
	var username = token.Claims.(jwt.MapClaims)["username"]

	// Look to see if user already has role
	var hasRoleAlready = utils.RoleCheck(c, role, username.(string))

	c.JSON(http.StatusOK, gin.H{"hasRoleAlready": hasRoleAlready})
}

// AddRole POST /roles
func AddRole(c *gin.Context) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	
	// Get user from token
	token := utils.ParseToken(c, jwtKey, "x-auth-token")
	var username = token.Claims.(jwt.MapClaims)["username"]

	var newRole roleRequest
	if err := c.BindJSON(&newRole); err != nil {
		return
	}

	// Look to see if user already has role
	var hasRoleAlready = utils.RoleCheck(c, newRole.Role, username.(string))

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
