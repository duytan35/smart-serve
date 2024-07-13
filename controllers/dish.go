package controllers

import (
	"fmt"
	"smart-serve/models"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Tags Dishes
// @Accept  json
// @Produce  json
// @Param data body models.CreateDishInput true "Dish Data"
// @Success 201 {object} Response{data=models.Dish}
// @Router /dishes [post]
// @Security BearerAuth
func CreateDish(c *gin.Context) {
	var createDish models.CreateDishInput

	if err := c.ShouldBindJSON(&createDish); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	restaurantId, _ := uuid.Parse(c.GetString("restaurantId"))
	dishGroup, err := models.GetDishGroup(createDish.DishGroupID)
	if err != nil || dishGroup.RestaurantID != restaurantId {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Dish group not found",
		})
		return
	}

	var dishGroupID uint
	fmt.Sscanf(createDish.DishGroupID, "%d", &dishGroupID)

	dish := models.Dish{
		DishGroupID: dishGroupID,
		Name:        createDish.Name,
		Description: createDish.Description,
		Price:       createDish.Price,
		Status:      createDish.Status,
	}

	newDish, err := models.CreateDish(dish, createDish.ImageIds)
	newDish.ImageIds = createDish.ImageIds

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Dish created successfully",
		Data:    newDish,
	})
}

// @Tags Dishes
// @Accept  json
// @Produce  json
// @Param id path string true "Dish ID"
// @Success 200 {object} Response{data=models.Dish}
// @Router /dishes/{id} [get]
// @Security BearerAuth
func GetDish(c *gin.Context) {
	id := c.Param("id")
	dish, err := models.GetDish(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    dish,
	})
}

// @Tags Dishes
// @Accept  json
// @Produce  json
// @Success 200 {object} Response{data=[]models.Dish}
// @Param dishGroupId query string true "Dish Group ID"
// @Router /dishes [get]
// @Security BearerAuth
func GetDishes(c *gin.Context) {
	restaurantId, _ := uuid.Parse(c.GetString("restaurantId"))

	dishGroupId := c.Query("dishGroupId")
	if dishGroupId == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "dishGroupId is required",
		})
		return
	}

	dishGroup, err := models.GetDishGroup(dishGroupId)
	if err != nil || dishGroup.RestaurantID != restaurantId {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Dish group not found",
		})
		return
	}

	dishes := models.GetDishes(dishGroupId)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    dishes,
	})
}

// @Tags Dishes
// @Accept  json
// @Produce  json
// @Param id path string true "Dish ID"
// @Param DishInput body models.UpdateDishInput true "Dish Data"
// @Success 200 {object} Response{data=models.Dish}
// @Router /dishes/{id} [put]
// @Security BearerAuth
func UpdateDish(c *gin.Context) {
	restaurantId, _ := uuid.Parse(c.GetString("restaurantId"))

	id := c.Param("id")

	var updateDishInput models.UpdateDishInput

	if err := c.ShouldBindJSON(&updateDishInput); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	updatedDish, err := models.UpdateDish(id, restaurantId, updateDishInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	updatedDish.ImageIds = updateDishInput.ImageIds

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Dish updated successfully",
		Data:    updatedDish,
	})
}

// @Tags Dishes
// @Accept  json
// @Produce  json
// @Param id path string true "Dish ID"
// @Success 200 {object} Response{data=nil}
// @Router /dishes/{id} [delete]
// @Security BearerAuth
func DeleteDish(c *gin.Context) {
	restaurantId, _ := uuid.Parse(c.GetString("restaurantId"))

	id := c.Param("id")

	if err := models.DeleteDish(id, restaurantId); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Dish deleted successfully",
	})
}
