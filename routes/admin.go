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
		group.GET("", controllers.GetRestaurants)
		group.GET("/:id", controllers.GetRestaurant)
		group.DELETE("/:id", controllers.DeleteRestaurant)
	}
}
