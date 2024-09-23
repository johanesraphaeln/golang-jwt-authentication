package main

import (
	"golang-jwt-auth/handlers"
	"golang-jwt-auth/middleware"
	"golang-jwt-auth/models"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	models.ConnectDatabase()

	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	// Public routes (do not require authentication)
	publicRoutes := router.Group("/api/public")
	{
		publicRoutes.POST("/login", handlers.Login)
		publicRoutes.POST("/register", handlers.Register)
	}

	// Protected routes (require authentication)
	protectedRoutes := router.Group("/api/protected")
	protectedRoutes.Use(middleware.AuthenticationMiddleware())
	{
		protectedRoutes.GET("/users", handlers.GetUsers)
	}

	router.Run(":3000")
}
