package routes

import (
	"saasmanagement/api/v1/controllers"
	"saasmanagement/api/v1/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAccountRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/accounts")
	{
		// Add the authentication middleware to all routes except login

		//userRoutes.Use(middlewares.AuthMiddleware()).GET("/", controllers.GetAllUsers)
		userRoutes.Use(middlewares.AuthMiddleware()).GET("/:user_id", controllers.GetAllAccountsByUserID)
		userRoutes.Use(middlewares.AuthMiddleware()).GET("/account", controllers.CreateAccount)
		userRoutes.Use(middlewares.AuthMiddleware()).GET("/:id", controllers.GetAccountByID)
		userRoutes.Use(middlewares.AuthMiddleware()).PUT("/:id", controllers.UpdateAccount)
		userRoutes.Use(middlewares.AuthMiddleware()).DELETE("/:id", controllers.DeleteAccount)

	}
}
