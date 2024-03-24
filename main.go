package main

import (
	"auth-api-go/controllers"
	"auth-api-go/models"
	"auth-api-go/redis"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	redis.ConnectRedis()
	models.ConnectDatabase()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	// cors, allow all and new header
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("x-auth-token")
	config.AddAllowHeaders("X-API-Token")
	router.Use(cors.New(config))

	router.GET("/", controllers.Index)
	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.Register)
	router.GET("/verify", controllers.Verify)
	router.DELETE("/", controllers.DeleteUser)
	router.DELETE("/session", controllers.DeleteUserSession)

	router.GET("/roles", controllers.GetRoles)
	router.GET("/roles/:role", controllers.DoesUserHaveRole)
	router.POST("/roles", controllers.AddRole)

	appRoutes := router.Group("/app")
	{
		appRoutes.GET("/verify", controllers.AppVerify)
		appRoutes.DELETE("/user/:username", controllers.AppDeleteUser)
	}

	// By default, it serves on :8080 unless a
	// PORT environment variable was defined.
	err := router.Run()
	if err != nil {
		fmt.Println("Error starting Server")
		return
	}
}
