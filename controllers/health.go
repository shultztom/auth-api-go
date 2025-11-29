package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index GET /
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": ""})
}
