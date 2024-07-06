package routes

import (
	"smart-serve/controllers"
	"smart-serve/middlewares"

	"github.com/gin-gonic/gin"
)

func addTableRoutes(r *gin.RouterGroup) {
	group := r.Group("tables")

	group.Use(middlewares.AuthGuard())
	{
		group.POST("", controllers.CreateTable)
		group.GET("", controllers.GetTables)
		group.GET("/:id", controllers.GetTable)
		group.PUT("/:id", controllers.UpdateTable)
		group.DELETE("/:id", controllers.DeleteTable)
	}
}
