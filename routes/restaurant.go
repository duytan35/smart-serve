package routes

import (
	"smart-serve/controllers"
	"smart-serve/middlewares"

	"github.com/gin-gonic/gin"
)

func addRestaurantRoutes(r *gin.RouterGroup) {
	group := r.Group("restaurants")

	group.POST("", controllers.CreateRestaurant)

	group.Use(middlewares.AuthGuard())
	{
		group.PATCH("/", controllers.UpdateRestaurant)
		group.PATCH("/steps", controllers.UpdateSteps)
	}
}
