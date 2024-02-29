package main

import (
	"net/http"

	"github.com/AdrianTworek/go-tasks-manager/handlers"
	"github.com/AdrianTworek/go-tasks-manager/initializers"
	"github.com/AdrianTworek/go-tasks-manager/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.InitDb()
	initializers.SyncDb()
}

func main() {
	r := gin.Default()

	api := r.Group("/api")

	// Healthcheck
	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Tasks Manager API!",
		})
	})

	// Users
	api.POST("/users", handlers.HandleRegister)

	// Auth
	api.POST("/auth/login", handlers.HandleLogin)

	// Protected routes
	authenticated := r.Group("/api")
	authenticated.Use(middleware.Authenticate)

	r.Run()
}
