package controllers

import (
	"fmt"
	"smart-serve/models"
	"smart-serve/utils"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Tags Orders
// @Summary Use for both restaurant and client
// @Accept  json
// @Produce  json
// @Param data body models.CreateOrderInput true "Order Data"
// @Success 201 {object} Response{data=models.OrderResponse}
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	var createOrder models.CreateOrderInput

	if err := c.ShouldBindJSON(&createOrder); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	table, err := models.GetTableById(fmt.Sprintf("%d", createOrder.TableID))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Table not found",
		})
		return
	}

	order, err := models.CreateOrder(createOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	newOrder, _ := models.GetOrder(strconv.FormatUint(uint64(order.ID), 10))

	utils.SendMessageToRoom(table.RestaurantID.String(), utils.SocketMessage{
		Event: "PLACE_ORDER",
		Data:  newOrder,
	})

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Order created successfully",
		Data:    newOrder,
	})
}

// @Tags Orders
// @Accept  json
// @Produce  json
// @Param tableId query string false "Table ID"
// @Param status query string false "Status"
// @Success 200 {object} Response{data=[]models.OrderResponse}
// @Router /orders [get]
// @Security BearerAuth
func GetOrders(c *gin.Context) {
	restaurantId := c.GetString("restaurantId")
	tableId := c.Query("tableId")
	status := c.Query("status")

	orders := models.GetOrders(restaurantId, tableId, status)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    orders,
	})
}

// @Tags Orders
// @Summary Use for both restaurant and client
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} Response{data=models.OrderResponse}
// @Router /orders/{id} [get]
func GetOrder(c *gin.Context) {
	id := c.Param("id")

	order, err := models.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    order,
	})
}

// @Tags Orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Param OrderInput body models.UpdateOrderInput true "Order Data"
// @Success 200 {object} Response{data=models.OrderResponse}
// @Router /orders/{id} [patch]
// @Security BearerAuth
func UpdateOrder(c *gin.Context) {
	restaurantId, _ := uuid.Parse(c.GetString("restaurantId"))

	id := c.Param("id")

	var orderInput models.UpdateOrderInput

	if err := c.ShouldBindJSON(&orderInput); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	updatedOrder, err := models.UpdateOrder(id, restaurantId, orderInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	newOrder, _ := models.GetOrder(strconv.FormatUint(uint64(updatedOrder.ID), 10))

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Order updated successfully",
		Data:    newOrder,
	})
}

// @Tags Orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} Response{data=nil}
// @Router /orders/{id} [delete]
// @Security BearerAuth
func DeleteOrder(c *gin.Context) {
	restaurantId, _ := uuid.Parse(c.GetString("restaurantId"))

	id := c.Param("id")

	if err := models.DeleteOrder(id, restaurantId); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Order deleted successfully",
	})
}

// @Tags Orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order Detail ID"
// @Param UpdateOrderStepInput body models.UpdateOrderStepInput true "Order Step Data"
// @Success 200 {object} Response{data=nil}
// @Router /orders/order-details/{id} [patch]
// @Security BearerAuth
func UpdateOrderDetailStep(c *gin.Context) {
	restaurantId, _ := uuid.Parse(c.GetString("restaurantId"))
	orderDetailId := c.Param("id")

	var updateOrderStepInput models.UpdateOrderStepInput

	if err := c.ShouldBindJSON(&updateOrderStepInput); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	orderDetail, err := models.UpdateOrderDetailStep(restaurantId, orderDetailId, *updateOrderStepInput.Step)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	utils.SendMessageToRoom(fmt.Sprintf("%s_%s", restaurantId.String(), strconv.Itoa(int(orderDetail.Order.TableID))), utils.SocketMessage{
		Event: "ORDER_UPDATED",
		Data:  orderDetail,
	})

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Order step updated successfully",
	})
}
