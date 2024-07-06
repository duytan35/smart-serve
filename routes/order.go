package routes

import (
	"smart-serve/controllers"
	"smart-serve/middlewares"

	"github.com/gin-gonic/gin"
)

func addOrderRoutes(r *gin.RouterGroup) {
	group := r.Group("orders")

	group.POST("", controllers.CreateOrder)
	group.GET("/:id", controllers.GetOrder)

	group.Use(middlewares.AuthGuard())
	{
		group.GET("", controllers.GetOrders)
		group.PATCH("/:id", controllers.UpdateOrder)
		group.DELETE("/:id", controllers.DeleteOrder)
	}
}
