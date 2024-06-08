package routes

import (
	"smart-serve/controllers"
	"smart-serve/middlewares"

	"github.com/gin-gonic/gin"
)

func addUserRoutes(r *gin.RouterGroup) {
	group := r.Group("users")

	group.POST("", controllers.CreateUser)

	group.Use(middlewares.JWTAuth())
	{
		group.GET("", controllers.GetUsers)
		group.GET("/:id", controllers.GetUser)
		group.PATCH("/:id", controllers.UpdateUser)
		group.DELETE("/:id", controllers.DeleteUser)
	}
}
