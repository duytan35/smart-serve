package routes

import (
	"smart-serve/controllers"

	"github.com/gin-gonic/gin"
)

func addRestaurantRoutes(r *gin.RouterGroup) {
	group := r.Group("restaurants")

	group.POST("", controllers.CreateRestaurant)
	group.GET("", controllers.GetRestaurants)
	group.GET("/:id", controllers.GetRestaurant)
	group.PATCH("/:id", controllers.UpdateRestaurant)
	group.DELETE("/:id", controllers.DeleteRestaurant)
}
