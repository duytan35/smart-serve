package middlewares

import (
	"net/http"
	"os"
	"strings"

	"smart-serve/constants"
	"smart-serve/controllers"
	"smart-serve/utils"

	"github.com/gin-gonic/gin"
)

func AuthGuard(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(roles) == 0 {
			roles = append(roles, constants.Restaurant)
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			if utils.Contains(roles, constants.Client) {
				c.Set("role", constants.Client)
				c.Next()
				return
			}

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

		var role string
		if claims.Email == os.Getenv("ADMIN_MAIL") {
			role = constants.Admin
		} else {
			role = constants.Restaurant
		}

		if !utils.Contains(roles, role) {
			c.JSON(http.StatusUnauthorized, controllers.Response{
				Success: false,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Set("restaurantId", claims.RestaurantID)
		c.Set("role", role)
		c.Next()
	}
}
