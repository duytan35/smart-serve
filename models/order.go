package models

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateOrderInput struct {
	TableID      uint               `json:"tableId" binding:"required" example:"1"`
	OrderDetails []OrderDetailInput `json:"orderDetails" binding:"required,dive"`
}

type UpdateOrderInput struct {
	Status       OrderStatus        `json:"status" example:"InProgress"  binding:"orderStatus"`
	OrderDetails []OrderDetailInput `json:"orderDetails" binding:"dive"`
}

type UpdateOrderStepInput struct {
	Step *uint `json:"step" example:"1" binding:"required"`
}

type OrderDetailInput struct {
	DishID      uint   `json:"dishId" binding:"required" example:"1"`
	Quantity    uint   `json:"quantity" binding:"required" example:"2"`
	DiscountIDs []uint `json:"discountIds" example:"1,2"`
	Note        string `json:"note" example:"Note"`
}

type OrderDetailResponse struct {
	ID               uint        `json:"id"`
	Quantity         uint        `json:"quantity"`
	Step             uint        `json:"step"`
	DiscountPercent  float64     `json:"discountPercent"`
	DishID           uint        `json:"dishId"`
	DishName         string      `json:"dishName"`
	DishPrice        float64     `json:"dishPrice"`
	DishDescription  string      `json:"dishDescription"`
	Note             string      `json:"note"`
	GroupOrderNumber uint        `json:"groupOrderNumber"`
	ImageIds         []uuid.UUID `json:"imageIds"`
	CreatedAt        time.Time   `json:"createdAt"`
	UpdatedAt        time.Time   `json:"updatedAt"`
}

type OrderStatus string

const (
	StatusInProgress OrderStatus = "InProgress"
	StatusComplete   OrderStatus = "Complete"
	StatusCancel     OrderStatus = "Cancel"
)

type OrderResponse struct {
	ID        uint        `json:"id"`
	TableId   uint        `json:"tableId"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`

	OrderDetails []OrderDetailResponse `json:"orderDetails"`
}

func GetOrders(restaurantId, tableId string, status string) []OrderResponse {
	var orders []Order
	var orderResponses = []OrderResponse{}

	query := DB.
		Preload("OrderDetails", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, order_id, quantity, step, dish_id, discount_percent, note, group_order_number, created_at, updated_at")
		}).
		Preload("OrderDetails.Dish", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, price, description")
		}).
		Joins("JOIN tables ON tables.id = orders.table_id").
		Where("tables.restaurant_id = ?", restaurantId)
	if tableId != "" {
		query = query.Where("table_id = ?", tableId)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Find(&orders)

	for _, order := range orders {
		orderResponse := OrderResponse{
			ID:        order.ID,
			TableId:   order.TableID,
			Status:    order.Status,
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
		}

		for _, detail := range order.OrderDetails {
			detailResponse := OrderDetailResponse{
				ID:               detail.ID,
				Quantity:         detail.Quantity,
				Step:             detail.Step,
				DiscountPercent:  detail.DiscountPercent,
				DishID:           detail.DishID,
				DishName:         detail.Dish.Name,
				DishPrice:        detail.Dish.Price,
				DishDescription:  detail.Dish.Description,
				Note:             detail.Note,
				GroupOrderNumber: detail.GroupOrderNumber,
				CreatedAt:        detail.CreatedAt,
				UpdatedAt:        detail.UpdatedAt,
			}

			orderResponse.OrderDetails = append(orderResponse.OrderDetails, detailResponse)
		}

		orderResponses = append(orderResponses, orderResponse)
	}

	return orderResponses
}

func CreateOrder(createOrder CreateOrderInput) (OrderResponse, error) {
	existOrderIdInProgress := GetOrderIdAtTable(strconv.FormatUint(uint64(createOrder.TableID), 10))

	tx := DB.Begin()

	if tx.Error != nil {
		return OrderResponse{}, tx.Error
	}

	var orderId uint
	maxGroupOrderNumber := uint(0)

	if existOrderIdInProgress == nil {
		order := Order{
			TableID: createOrder.TableID,
		}

		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			return OrderResponse{}, err
		}

		orderId = order.ID
	} else {
		orderId = *existOrderIdInProgress
		maxGroupOrderNumber = GetMaxGroupOrderNumber(orderId)
	}

	for _, detail := range createOrder.OrderDetails {
		orderDetail := OrderDetail{
			OrderID:          orderId,
			DishID:           detail.DishID,
			Quantity:         detail.Quantity,
			Note:             detail.Note,
			GroupOrderNumber: maxGroupOrderNumber + 1,
		}

		if err := tx.Create(&orderDetail).Error; err != nil {
			tx.Rollback()
			return OrderResponse{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return OrderResponse{}, err
	}

	order, _ := GetOrder(strconv.FormatUint(uint64(orderId), 10))

	return order, nil
}

func GetOrder(id string) (OrderResponse, error) {
	var order Order

	if err := DB.
		Preload("OrderDetails", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, order_id, quantity, step, dish_id, discount_percent, note, group_order_number, created_at, updated_at")
		}).
		Preload("OrderDetails.Dish", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, price, description")
		}).
		Preload("OrderDetails.Dish.Images").
		Where("orders.id = ?", id).
		First(&order).Error; err != nil {
		return OrderResponse{}, err
	}

	orderResponse := OrderResponse{
		ID:        order.ID,
		TableId:   order.TableID,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}

	for _, detail := range order.OrderDetails {
		detailResponse := OrderDetailResponse{
			ID:               detail.ID,
			Quantity:         detail.Quantity,
			Step:             detail.Step,
			DiscountPercent:  detail.DiscountPercent,
			DishID:           detail.DishID,
			DishName:         detail.Dish.Name,
			DishPrice:        detail.Dish.Price,
			DishDescription:  detail.Dish.Description,
			Note:             detail.Note,
			GroupOrderNumber: detail.GroupOrderNumber,
			ImageIds:         detail.Dish.ImageIds,
			CreatedAt:        detail.CreatedAt,
			UpdatedAt:        detail.UpdatedAt,
		}

		orderResponse.OrderDetails = append(orderResponse.OrderDetails, detailResponse)
	}

	return orderResponse, nil
}

func UpdateOrder(id string, restaurantId uuid.UUID, orderInput UpdateOrderInput) (Order, error) {
	var updatedOrder Order

	tx := DB.Begin()

	if err := tx.
		Joins("JOIN tables ON tables.id = orders.table_id").
		Where("orders.id = ? AND tables.restaurant_id = ?", id, restaurantId).
		First(&updatedOrder).Error; err != nil {
		tx.Rollback()
		return Order{}, err
	}

	updatedOrder.Status = orderInput.Status
	if err := tx.Save(&updatedOrder).Error; err != nil {
		tx.Rollback()
		return Order{}, err
	}

	if len(orderInput.OrderDetails) > 0 {
		// Delete existing order details
		if err := tx.Where("order_id = ?", updatedOrder.ID).Delete(&OrderDetail{}).Error; err != nil {
			tx.Rollback()
			return Order{}, err
		}

		// Create new order details
		for _, detail := range orderInput.OrderDetails {
			orderDetail := OrderDetail{
				OrderID:  updatedOrder.ID,
				DishID:   detail.DishID,
				Quantity: detail.Quantity,
			}
			if err := tx.Create(&orderDetail).Error; err != nil {
				tx.Rollback()
				return Order{}, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return Order{}, err
	}

	return updatedOrder, nil
}

func DeleteOrder(id string, restaurantId uuid.UUID) error {
	var order Order

	if err := DB.
		Joins("JOIN tables ON tables.id = orders.table_id").
		Where("orders.id = ? AND tables.restaurant_id = ?", id, restaurantId).
		First(&order).Error; err != nil {
		return err
	}

	if err := DB.Delete(&order, id).Error; err != nil {
		return err
	}

	return nil
}

func UpdateOrderDetailStep(restaurantId uuid.UUID, orderDetailId string, step uint) (OrderDetail, error) {
	var orderDetail OrderDetail

	if err := DB.
		Joins("JOIN orders ON orders.id = order_details.order_id").
		Joins("JOIN tables ON tables.id = orders.table_id").
		Preload("Order").
		Where("tables.restaurant_id = ?", restaurantId).
		First(&orderDetail, orderDetailId).Error; err != nil {
		return OrderDetail{}, err
	}

	orderDetail.Step = step

	if err := DB.Save(&orderDetail).Error; err != nil {
		return OrderDetail{}, err
	}

	return orderDetail, nil
}

func GetMaxGroupOrderNumber(orderId uint) uint {
	var orderDetail OrderDetail
	DB.Where("order_id = ?", orderId).Order("group_order_number desc").First(&orderDetail)

	return orderDetail.GroupOrderNumber
}
