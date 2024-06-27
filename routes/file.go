package routes

import (
	"smart-serve/controllers"
	"smart-serve/middlewares"

	"github.com/gin-gonic/gin"
)

func addFileRoutes(r *gin.RouterGroup) {
	group := r.Group("files")

	group.GET("/:id", controllers.GetFile)

	group.Use(middlewares.JWTAuth())
	{
		group.POST("", controllers.UploadFile)
		group.PUT("/:id", controllers.UpdateFile)
		group.DELETE("/:id", controllers.DeleteFile)
	}
}
