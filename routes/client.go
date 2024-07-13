package routes

import (
	"smart-serve/controllers"

	"github.com/gin-gonic/gin"
)

func addClientRoutes(r *gin.RouterGroup) {
	group := r.Group("client")

	group.GET("/menu", controllers.GetMenu)
	group.GET("/order", controllers.GetOrderByClient)
}
