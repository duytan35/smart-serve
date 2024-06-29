package middlewares

import (
	"net/http"
	"strings"

	"smart-serve/controllers"
	"smart-serve/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, controllers.Response{
				Success: false,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		if len(strings.Split(authHeader, "Bearer ")) != 2 {
			c.JSON(http.StatusUnauthorized, controllers.Response{
				Success: false,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, controllers.Response{
				Success: false,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Set("restaurantId", claims.RestaurantID)
		c.Next()
	}
}
