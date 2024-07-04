package controllers

import (
	"smart-serve/models"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Tags DishGroups
// @Accept  json
// @Produce  json
// @Param data body models.DishGroupInput true "DishGroup Data"
// @Success 201 {object} Response{data=models.DishGroup}
// @Router /dish-groups [post]
// @Security BearerAuth
func CreateDishGroup(c *gin.Context) {
	var createDishGroup models.DishGroupInput

	if err := c.ShouldBindJSON(&createDishGroup); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	restaurantId, _ := uuid.Parse(c.GetString("restaurantId"))
	dishGroup := models.DishGroup{
		Name:         createDishGroup.Name,
		RestaurantID: restaurantId,
	}

	res, err := models.CreateDishGroup(dishGroup)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "DishGroup created successfully",
		Data:    res,
	})
}

// @Tags DishGroups
// @Accept  json
// @Produce  json
// @Param id path string true "DishGroup ID"
// @Success 200 {object} Response{data=models.DishGroup}
// @Router /dish-groups/{id} [get]
// @Security BearerAuth
func GetDishGroup(c *gin.Context) {
	id := c.Param("id")
	dishGroup, err := models.GetDishGroup(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    dishGroup,
	})
}

// @Tags DishGroups
// @Accept  json
// @Produce  json
// @Success 200 {object} Response{data=[]models.DishGroup}
// @Router /dish-groups [get]
// @Security BearerAuth
func GetDishGroups(c *gin.Context) {
	restaurantId := c.GetString("restaurantId")

	dishGroups := models.GetDishGroups(restaurantId)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    dishGroups,
	})
}

// @Tags DishGroups
// @Accept  json
// @Produce  json
// @Param id path string true "DishGroup ID"
// @Param dishGroup body models.DishGroupInput true "DishGroup Data"
// @Success 200 {object} Response{data=models.DishGroup}
// @Router /dish-groups/{id} [put]
// @Security BearerAuth
func UpdateDishGroup(c *gin.Context) {
	restaurantId, _ := uuid.Parse(c.GetString("restaurantId"))

	id := c.Param("id")

	var dishGroupInput models.DishGroupInput

	if err := c.ShouldBindJSON(&dishGroupInput); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	updatedDishGroup, err := models.UpdateDishGroup(id, restaurantId, dishGroupInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "DishGroup updated successfully",
		Data:    updatedDishGroup,
	})
}

// @Tags DishGroups
// @Accept  json
// @Produce  json
// @Param id path string true "DishGroup ID"
// @Success 200 {object} Response{data=nil}
// @Router /dish-groups/{id} [delete]
// @Security BearerAuth
func DeleteDishGroup(c *gin.Context) {
	id := c.Param("id")
	restaurantId, _ := uuid.Parse(c.GetString("restaurantId"))

	if err := models.DeleteDishGroup(id, restaurantId); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Dish group deleted successfully",
	})
}
