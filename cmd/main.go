package main

import (
	"net/http"

	"github.com/AdrianTworek/go-tasks-manager/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
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

	r.Run()
}
