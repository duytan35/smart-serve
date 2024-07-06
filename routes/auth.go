package routes

import (
	"smart-serve/controllers"
	"smart-serve/middlewares"

	"github.com/gin-gonic/gin"
)

func addAuthRoutes(r *gin.RouterGroup) {
	group := r.Group("auth")

	group.POST("sign-in", controllers.SignIn)
	group.GET("me", middlewares.AuthGuard(), controllers.GetMe)
}
