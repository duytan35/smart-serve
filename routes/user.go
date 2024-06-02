package routes

import (
	"smart-serve/controllers"

	"github.com/gin-gonic/gin"
)

func addUserRoutes(r *gin.RouterGroup) {
	group := r.Group("users")

	group.POST("", controllers.CreateUser)
	group.GET("", controllers.GetUsers)
	group.GET("/:id", controllers.GetUser)
	group.PATCH("/:id", controllers.UpdateUser)
	group.DELETE("/:id", controllers.DeleteUser)
}
