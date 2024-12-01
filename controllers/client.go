package controllers

import (
	"smart-serve/models"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

// @Tags Client
// @Accept  json
// @Produce  json
// @Param restaurantId query string true "Restaurant ID"
// @Param tableId query string true "table ID"
// @Success 200 {object} Response{data=models.MenuResponse}
// @Router /client/menu [get]
func GetMenu(c *gin.Context) {
	restaurantId := c.Query("restaurantId")
	tableId := c.Query("tableId")

	checkTable := models.CheckTableExist(tableId, restaurantId)

	if !checkTable {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Restaurant or table not found",
		})
		return
	}

	menuResponse := models.GetMenu(restaurantId)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    menuResponse,
	})
}

// @Tags Client
// @Accept  json
// @Produce  json
// @Param tableId query string true "Table ID"
// @Success 200 {object} Response{data=models.OrderResponse}
// @Router /client/order [get]
// @Security BearerAuth
func GetOrderByClient(c *gin.Context) {
	tableId := c.Query("tableId")

	orderId := models.GetOrderIdAtTable(tableId)

	if orderId == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "No Order found",
		})
		return
	}

	order, _ := models.GetOrder(strconv.FormatUint(uint64(*orderId), 10))

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    order,
	})
}
