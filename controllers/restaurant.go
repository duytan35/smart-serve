package controllers

import (
	"smart-serve/models"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Tags Restaurants
// @Accept  json
// @Produce  json
// @Param data body models.CreateRestaurantInput true "Restaurant Data"
// @Success 201 {object} Response{data=models.RestaurantResponse}
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
// @Success 200 {object} Response{data=models.RestaurantResponse}
// @Router /restaurants [patch]
// @Security BearerAuth
func UpdateRestaurant(c *gin.Context) {
	restaurantId := c.GetString("restaurantId")
	var restaurant models.UpdateRestaurantInput

	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	updatedRestaurant, err := models.UpdateRestaurant(restaurantId, restaurant)
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

// @Tags Restaurants
// @Accept  json
// @Produce  json
// @Param restaurant body models.UpdateStepsInput true "Steps Data"
// @Success 200 {object} Response{data=nil}
// @Router /restaurants/steps [patch]
// @Security BearerAuth
func UpdateSteps(c *gin.Context) {
	restaurantId := c.GetString("restaurantId")
	restaurantUUID, _ := uuid.Parse(restaurantId)
	var input models.UpdateStepsInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := models.RemoveAllOrderSteps(restaurantId); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	var newSteps []models.OrderStep
	for i, stepName := range input.Steps {
		newSteps = append(newSteps, models.OrderStep{
			RestaurantID: restaurantUUID,
			Step:         uint(i),
			Name:         stepName,
		})
	}

	if err := models.CreateOrderSteps(newSteps, nil); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})

		c.JSON(http.StatusOK, Response{
			Success: true,
			Message: "Steps updated successfully",
		})
	}
}
