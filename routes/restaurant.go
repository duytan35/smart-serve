package routes

import (
	"smart-serve/controllers"
	"smart-serve/middlewares"

	"github.com/gin-gonic/gin"
)

func addRestaurantRoutes(r *gin.RouterGroup) {
	group := r.Group("restaurants")

	group.POST("", controllers.CreateRestaurant)

	group.Use(middlewares.JWTAuth())
	{

		// group.GET("/:id", controllers.GetRestaurant)
		group.PATCH("/", controllers.UpdateRestaurant)
		// group.DELETE("/:id", controllers.DeleteRestaurant)
	}
}
