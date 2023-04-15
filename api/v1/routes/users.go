package routes

import (
	"saasmanagement/api/v1/controllers"
	"saasmanagement/api/v1/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUsersRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/users")
	{
		// Add the authentication middleware to all routes except login

		//userRoutes.Use(middlewares.AuthMiddleware()).GET("/", controllers.GetAllUsers)
		userRoutes.GET("/", controllers.GetAllUsers)
		userRoutes.POST("/login", controllers.Login)
		userRoutes.POST("/user", controllers.CreateUser)
		userRoutes.GET("/refresh", controllers.RefreshToken)

		userRoutes.Use(middlewares.AuthMiddleware()).GET("/:id", controllers.GetUserById)
		userRoutes.Use(middlewares.AuthMiddleware()).PUT("/:id", controllers.UpdateUser)
		userRoutes.Use(middlewares.AuthMiddleware()).DELETE("/:id", controllers.DeleteUser)

	}
}

func SetupCustomEndpointRoutes(router *gin.RouterGroup) {
	customRoutes := router.Group("/services")
	{

		customRoutes.Use(middlewares.AuthMiddleware()).GET("/external-endpoint", func(c *gin.Context) {
			url := "http://localhost:8082/api/v1/users"
			controllers.ExternalEndpoint(c, url)
		})

	}
}
