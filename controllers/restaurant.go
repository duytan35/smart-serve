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
// @Success 201 {object} models.Restaurant
// @Failure 400 {object} models.ErrorResponse
// @Router /restaurants [post]
func CreateRestaurant(c *gin.Context) {
	var restaurant models.Restaurant

	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurant, err := models.CreateRestaurant(restaurant)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, restaurant)
}

// @Tags Restaurants
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Restaurant
// @Router /restaurants [get]
func GetRestaurants(c *gin.Context) {
	restaurants := models.GetRestaurants()
	c.JSON(http.StatusOK, restaurants)
}

// @Tags Restaurants
// @Accept  json
// @Produce  json
// @Param id path string true "Restaurant ID"
// @Success 200 {object} models.Restaurant
// @Failure 404 {object} models.ErrorResponse
// @Router /restaurants/{id} [get]
func GetRestaurant(c *gin.Context) {
	id := c.Param("id")

	restaurant, err := models.GetRestaurant(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, restaurant)
}

// @Tags Restaurants
// @Accept  json
// @Produce  json
// @Param id path string true "Restaurant ID"
// @Param restaurant body models.UpdateRestaurantInput true "Restaurant Data"
// @Success 200 {object} models.Restaurant
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /restaurants/{id} [patch]
func UpdateRestaurant(c *gin.Context) {
	id := c.Param("id")
	var restaurant models.UpdateRestaurantInput

	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updatedRestaurant, err := models.UpdateRestaurant(id, restaurant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedRestaurant)
}

// @Tags Restaurants
// @Accept  json
// @Produce  json
// @Param id path string true "Restaurant ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Router /restaurants/{id} [delete]
func DeleteRestaurant(c *gin.Context) {
	id := c.Param("id")
	if err := models.DeleteRestaurant(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "restaurant deleted successfully"})
}
