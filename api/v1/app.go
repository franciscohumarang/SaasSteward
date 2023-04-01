package api

import (
	"saasmanagement/api/v1/routes"
	"saasmanagement/config"

	"github.com/gin-gonic/gin"
)

func Start() {
	// Initialize the application configuration
	config.Init()

	// Set the Gin router
	router := gin.Default()

	// Set up the API routes
	apiV1Routes := router.Group("/api/v1")
	{
		routes.SetupUsersRoutes(apiV1Routes)
	}

	// Start the server
	router.Run()
}
