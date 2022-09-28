package utils

import (
	"auth-api-go/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RoleCheck(c *gin.Context, roleToCheck string, username string) bool {
	// Set up array
	var roles []models.Roles

	// Query DB and add to roles var
	result := models.DB.Find(&roles, "username = ?", username)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting roles for User!"})
		// TODO look into if this ever gets called
		return true
	}

	// Check if role already exists
	var hasRoleAlready = false
	for i := 0; i < len(roles); i++ {
		if roles[i].Role == roleToCheck {
			hasRoleAlready = true
			break
		}
	}

	return hasRoleAlready
}
