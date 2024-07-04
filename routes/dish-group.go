package routes

import (
	"smart-serve/controllers"
	"smart-serve/middlewares"

	"github.com/gin-gonic/gin"
)

func addDishGroupRoutes(r *gin.RouterGroup) {
	group := r.Group("dish-groups")

	group.Use(middlewares.JWTAuth())
	{
		group.POST("", controllers.CreateDishGroup)
		group.GET("", controllers.GetDishGroups)
		group.GET("/:id", controllers.GetDishGroup)
		group.PUT("/:id", controllers.UpdateDishGroup)
		group.DELETE("/:id", controllers.DeleteDishGroup)
	}
}
