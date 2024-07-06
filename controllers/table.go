package controllers

import (
	"smart-serve/models"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Tags Tables
// @Accept  json
// @Produce  json
// @Param data body models.CreateTableInput true "Table Data"
// @Success 201 {object} Response{data=models.Table}
// @Router /tables [post]
// @Security BearerAuth
func CreateTable(c *gin.Context) {
	var createTable models.CreateTableInput

	if err := c.ShouldBindJSON(&createTable); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	restaurantId, _ := uuid.Parse(c.GetString("restaurantId"))

	table := models.Table{
		RestaurantID: restaurantId,
		Name:         createTable.Name,
		Seats:        createTable.Seats,
	}

	res, err := models.CreateTable(table)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Table created successfully",
		Data:    res,
	})
}

// @Tags Tables
// @Accept  json
// @Produce  json
// @Param id path string true "Table ID"
// @Success 200 {object} Response{data=models.Table}
// @Router /tables/{id} [get]
// @Security BearerAuth
func GetTable(c *gin.Context) {
	restaurantId := c.GetString("restaurantId")
	id := c.Param("id")

	dish, err := models.GetTable(id, restaurantId)
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

// @Tags Tables
// @Accept  json
// @Produce  json
// @Success 200 {object} Response{data=[]models.Table}
// @Router /tables [get]
// @Security BearerAuth
func GetTables(c *gin.Context) {
	restaurantId := c.GetString("restaurantId")

	tables := models.GetTables(restaurantId)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    tables,
	})
}

// @Tags Tables
// @Accept  json
// @Produce  json
// @Param id path string true "Table ID"
// @Param TableInput body models.UpdateTableInput true "Table Data"
// @Success 200 {object} Response{data=models.Table}
// @Router /tables/{id} [put]
// @Security BearerAuth
func UpdateTable(c *gin.Context) {
	restaurantId, _ := uuid.Parse(c.GetString("restaurantId"))

	id := c.Param("id")

	var tableInput models.UpdateTableInput

	if err := c.ShouldBindJSON(&tableInput); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	updatedTable, err := models.UpdateTable(id, restaurantId, tableInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Table updated successfully",
		Data:    updatedTable,
	})
}

// @Tags Tables
// @Accept  json
// @Produce  json
// @Param id path string true "Table ID"
// @Success 200 {object} Response{data=nil}
// @Router /tables/{id} [delete]
// @Security BearerAuth
func DeleteTable(c *gin.Context) {
	restaurantId, _ := uuid.Parse(c.GetString("restaurantId"))

	id := c.Param("id")

	if err := models.DeleteTable(id, restaurantId); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Table deleted successfully",
	})
}
