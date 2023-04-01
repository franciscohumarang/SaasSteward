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

		userRoutes.Use(middlewares.AuthMiddleware()).GET("/", controllers.GetAllUsers)
		userRoutes.Use(middlewares.AuthMiddleware()).GET("/:id", controllers.GetUserById)
		userRoutes.Use(middlewares.AuthMiddleware()).PUT("/:id", controllers.UpdateUser)
		userRoutes.Use(middlewares.AuthMiddleware()).DELETE("/:id", controllers.DeleteUser)
		userRoutes.POST("/login", controllers.Login)
		userRoutes.POST("/", controllers.CreateUser)
	}
}
