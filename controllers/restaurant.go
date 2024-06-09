package controllers

import (
	"smart-serve/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

// @Tags Restaurants
// @Accept  json
// @Produce  json
// @Param data body models.CreateRestaurantInput true "Restaurant Data"
// @Success 201 {object} Response{data=models.Restaurant}
// @Router /restaurants [post]
func CreateRestaurant(c *gin.Context) {
	var createRestaurant models.CreateRestaurantInput

	if err := c.ShouldBindJSON(&createRestaurant); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	restaurant := models.Restaurant{
		Name:     createRestaurant.Name,
		Email:    createRestaurant.Email,
		Password: createRestaurant.Password,
		Phone:    createRestaurant.Phone,
		Address:  createRestaurant.Address,
	}

	res, err := models.CreateRestaurant(restaurant)

	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Restaurant created successfully",
		Data:    res,
	})
}

// @Tags Restaurants
// @Accept  json
// @Produce  json
// @Param restaurant body models.UpdateRestaurantInput true "Restaurant Data"
// @Success 200 {object} Response{data=models.Restaurant}
// @Router /restaurants [patch]
// @Security BearerAuth
func UpdateRestaurant(c *gin.Context) {
	id := c.GetString("id")
	var restaurant models.UpdateRestaurantInput

	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	updatedRestaurant, err := models.UpdateRestaurant(id, restaurant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Restaurant updated successfully",
		Data:    updatedRestaurant,
	})
}
