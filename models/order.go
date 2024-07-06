package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateOrderInput struct {
	TableID      uint               `json:"tableId" binding:"required" example:"1"`
	OrderDetails []OrderDetailInput `json:"orderDetails" binding:"required,dive"`
}

type UpdateOrderInput struct {
	Status       uint               `json:"status" example:"1"`
	OrderDetails []OrderDetailInput `json:"orderDetails" binding:"dive"`
}

type OrderDetailInput struct {
	DishID      uint   `json:"dishId" binding:"required" example:"1"`
	Quantity    uint   `json:"quantity" binding:"required" example:"2"`
	DiscountIDs []uint `json:"discountIds" binding:"required" example:"1,2"`
}

type OrderDetailResponse struct {
	ID              uint      `json:"id"`
	Quantity        uint      `json:"quantity"`
	DiscountPercent float64   `json:"discountPercent"`
	DishID          uint      `json:"dishId"`
	DishName        string    `json:"dishName"`
	DishPrice       float64   `json:"dishPrice"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type OrderResponse struct {
	ID        uint      `json:"id"`
	TableId   uint      `json:"tableId"`
	Status    uint      `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	OrderDetails []OrderDetailResponse `json:"orderDetails"`
}

func GetOrders(restaurantId, tableId string) []OrderResponse {
	var orders []Order
	var orderResponses = []OrderResponse{}

	query := DB.
		Preload("OrderDetails", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, order_id, quantity, dish_id, discount_percent, created_at, updated_at")
		}).
		Preload("OrderDetails.Dish", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, price")
		}).
		Joins("JOIN tables ON tables.id = orders.table_id").
		Where("tables.restaurant_id = ?", restaurantId)
	if tableId != "" {
		query = query.Where("table_id = ?", tableId)
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
				ID:              detail.ID,
				Quantity:        detail.Quantity,
				DiscountPercent: detail.DiscountPercent,
				DishID:          detail.DishID,
				DishName:        detail.Dish.Name,
				DishPrice:       detail.Dish.Price,
				CreatedAt:       detail.CreatedAt,
				UpdatedAt:       detail.UpdatedAt,
			}

			orderResponse.OrderDetails = append(orderResponse.OrderDetails, detailResponse)
		}

		orderResponses = append(orderResponses, orderResponse)
	}

	return orderResponses
}

func CreateOrder(createOrder CreateOrderInput) (Order, error) {
	tx := DB.Begin()

	if tx.Error != nil {
		return Order{}, tx.Error
	}

	order := Order{
		TableID: createOrder.TableID,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return Order{}, err
	}

	for _, detail := range createOrder.OrderDetails {
		orderDetail := OrderDetail{
			OrderID:  order.ID,
			DishID:   detail.DishID,
			Quantity: detail.Quantity,
		}

		if err := tx.Create(&orderDetail).Error; err != nil {
			tx.Rollback()
			return Order{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return Order{}, err
	}

	return order, nil
}

func GetOrder(id string) (OrderResponse, error) {
	var order Order

	if err := DB.
		Preload("OrderDetails", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, order_id, quantity, dish_id, discount_percent, created_at, updated_at")
		}).
		Preload("OrderDetails.Dish", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, price")
		}).
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
			ID:              detail.ID,
			Quantity:        detail.Quantity,
			DiscountPercent: detail.DiscountPercent,
			DishID:          detail.DishID,
			DishName:        detail.Dish.Name,
			DishPrice:       detail.Dish.Price,
			CreatedAt:       detail.CreatedAt,
			UpdatedAt:       detail.UpdatedAt,
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
