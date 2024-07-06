package routes

import (
	"smart-serve/controllers"
	"smart-serve/middlewares"

	"github.com/gin-gonic/gin"
)

func addDishRoutes(r *gin.RouterGroup) {
	group := r.Group("dishes")

	group.Use(middlewares.AuthGuard())
	{
		group.POST("", controllers.CreateDish)
		group.GET("", controllers.GetDishes)
		group.GET("/:id", controllers.GetDish)
		group.PUT("/:id", controllers.UpdateDish)
		group.DELETE("/:id", controllers.DeleteDish)
	}
}
