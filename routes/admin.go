package routes

import (
	"smart-serve/controllers"
	"smart-serve/middlewares"

	"github.com/gin-gonic/gin"
)

func addAdminRoutes(r *gin.RouterGroup) {
	group := r.Group("admin")

	group.Use(middlewares.JWTAuth())
	{
		group.GET("restaurants", controllers.GetRestaurants)
		group.GET("restaurants/:id", controllers.GetRestaurant)
		group.DELETE("restaurants/:id", controllers.DeleteRestaurant)
	}
}
