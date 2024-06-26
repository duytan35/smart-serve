package routes

import (
	"smart-serve/controllers"
	"smart-serve/middlewares"

	"github.com/gin-gonic/gin"
)

func addFileRoutes(r *gin.RouterGroup) {
	group := r.Group("files")

	group.Use(middlewares.JWTAuth())
	{
		group.POST("", controllers.Upload)
		group.GET("/:id", controllers.Upload)
		group.PUT("/:id", controllers.Upload)
		group.DELETE("/:id", controllers.Upload)
	}
}
