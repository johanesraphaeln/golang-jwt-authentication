package main

import (
	"golang-jwt-auth/handlers"
	"golang-jwt-auth/middleware"
	"golang-jwt-auth/models"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	models.ConnectRelationalDatabase()
	models.ConnectRealtimeDatabase()

	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World!"})
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
		protectedRoutes.POST("/projects", func(c *gin.Context) {
			handlers.SaveProject(c, models.FirebaseDB)
		})
		protectedRoutes.GET("/projects", func(c *gin.Context) {
			handlers.GetProject(c, models.FirebaseDB)
		})
		protectedRoutes.DELETE("/projects", func(c *gin.Context) {
			handlers.DeleteProject(c, models.FirebaseDB)
		})
	}

	router.Run(":3000")
}
