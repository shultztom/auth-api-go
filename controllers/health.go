package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Index GET /
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": ""})
}
