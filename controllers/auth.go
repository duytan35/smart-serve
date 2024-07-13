package controllers

import (
	"smart-serve/models"
	"smart-serve/utils"

	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignInData struct {
	Email    string `json:"email" binding:"required,email" example:"example@gmail.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

type SignInResponse struct {
	models.RestaurantResponse
	AccessToken string `json:"accessToken"`
}

// @Tags Auth
// @Accept json
// @Produce json
// @Param data body SignInData true "Sign in data"
// @Success 200 {object} Response{data=SignInResponse}
// @Router /auth/sign-in [post]
func SignIn(c *gin.Context) {
	var signInData SignInData

	if err := c.ShouldBindJSON(&signInData); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	restaurant, err := models.GetRestaurantByEmail(signInData.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Email not found",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(restaurant.Password), []byte(signInData.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid password",
		})
		return
	}

	accessToken, _ := utils.GenerateJWT(restaurant.ID.String(), restaurant.Email)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: SignInResponse{
			AccessToken: accessToken,
			RestaurantResponse: models.RestaurantResponse{
				ID:      restaurant.ID.String(),
				Name:    restaurant.Name,
				Phone:   restaurant.Phone,
				Email:   restaurant.Email,
				Address: restaurant.Address,
				Avatar:  restaurant.Avatar,
			},
		},
	})
}

// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=models.RestaurantResponse}
// @Router /auth/me [get]
// @Security BearerAuth
func GetMe(c *gin.Context) {
	restaurantId := c.GetString("restaurantId")

	restaurant, err := models.GetRestaurant(restaurantId)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Restaurant not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    restaurant,
	})
}
