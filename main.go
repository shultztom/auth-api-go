package main

import (
	"auth-api-go/controllers"
	"auth-api-go/models"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDatabase()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.Register)
	router.GET("/verify", controllers.Verify)

	// By default, it serves on :8080 unless a
	// PORT environment variable was defined.
	err := router.Run()
	if err != nil {
		fmt.Println("Error starting Server")
		return
	}
}
