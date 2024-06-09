package controllers

import (
	"net/http"
	"smart-serve/models"

	"github.com/gin-gonic/gin"
)

// @Tags Admin
// @Accept  json
// @Produce  json
// @Success 200 {object} Response{data=[]models.Restaurant}
// @Router /restaurants [get]
// @Security BearerAuth
func GetRestaurants(c *gin.Context) {
	restaurants := models.GetRestaurants()
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    restaurants,
	})
}

// @Tags Admin
// @Accept  json
// @Produce  json
// @Param id path string true "Restaurant ID"
// @Success 200 {object} Response{data=models.Restaurant}
// @Router /restaurants/{id} [get]
// @Security BearerAuth
func GetRestaurant(c *gin.Context) {
	id := c.Param("id")

	restaurant, err := models.GetRestaurant(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    restaurant,
	})
}

// @Tags Admin
// @Accept  json
// @Produce  json
// @Param id path string true "Restaurant ID"
// @Success 200 {object} Response{data=nil}
// @Router /restaurants/{id} [delete]
// @Security BearerAuth
func DeleteRestaurant(c *gin.Context) {
	id := c.Param("id")
	if err := models.DeleteRestaurant(id); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "restaurant deleted successfully",
	})
}
