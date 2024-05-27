package routes

import (
	"smart-serve/controllers"

	"github.com/gin-gonic/gin"
)

func Config(r *gin.Engine) {
	r.GET("/users", controllers.GetUsers)
	r.POST("/users", controllers.CreateUser)
	r.GET("/users/:id", controllers.GetUser)
	r.PUT("/users/:id", controllers.UpdateUser)
	r.DELETE("/users/:id", controllers.DeleteUser)
}
