package routes

import (
	"saasmanagement/api/v1/controllers"
	"saasmanagement/api/v1/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupBillingRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/billings")
	{

		userRoutes.Use(middlewares.AuthMiddleware()).GET("/:user_id", controllers.GetAllBillingsByUserID)
		userRoutes.Use(middlewares.AuthMiddleware()).GET("/billing", controllers.CreateBilling)
		userRoutes.Use(middlewares.AuthMiddleware()).GET("/:id", controllers.GetBillingByID)
		userRoutes.Use(middlewares.AuthMiddleware()).PUT("/:id", controllers.UpdateBilling)
		userRoutes.Use(middlewares.AuthMiddleware()).DELETE("/:id", controllers.DeleteBilling)

	}
}
