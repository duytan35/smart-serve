package models

import (
	"strconv"

	"github.com/google/uuid"
)

type MenuResponse struct {
	RestaurantID      uuid.UUID `json:"restaurantId"`
	RestaurantName    string    `json:"restaurantName"`
	RestaurantAddress string    `json:"restaurantAddress"`
	RestaurantAvatar  string    `json:"restaurantAvatar"`

	Menu []MenuDishGroup `json:"menu"`
}

type MenuDishGroup struct {
	GroupID   uint   `json:"groupId"`
	GroupName string `json:"groupName"`

	Dishes []MenuDish `json:"dishes"`
}

type MenuDish struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name" gorm:"not null" binding:"required"`
	Description string  `json:"description"` // optional
	Price       float64 `json:"price" gorm:"not null" binding:"required"`
}

func GetMenu(restaurantId string) MenuResponse {
	var restaurant Restaurant
	DB.
		Preload("DishGroup").
		Preload("DishGroup.Dishes").
		Where("id = ?", restaurantId).First(&restaurant)

	var menu []MenuDishGroup

	for _, group := range restaurant.DishGroup {
		dishes := []MenuDish{}

		for _, dish := range group.Dishes {
			dishes = append(dishes, MenuDish{
				ID:          dish.ID,
				Name:        dish.Name,
				Description: dish.Description,
				Price:       dish.Price,
			})
		}

		menu = append(menu, MenuDishGroup{
			GroupID:   group.ID,
			GroupName: group.Name,
			Dishes:    dishes,
		})
	}

	menuResponse := MenuResponse{
		RestaurantID:      restaurant.ID,
		RestaurantName:    restaurant.Name,
		RestaurantAddress: restaurant.Address,
		RestaurantAvatar:  restaurant.Avatar,
		Menu:              menu,
	}

	return menuResponse
}

func GetOrderAtTable(tableId string) *OrderResponse {
	var order Order

	if err := DB.Where("table_id = ? AND status = ?", tableId, "InProgress").
		Order("created_at DESC").
		First(&order).Error; err != nil {
		return nil
	}

	orderResponse, err := GetOrder(strconv.FormatUint(uint64(order.ID), 10))
	if err != nil {
		return nil
	}

	return &orderResponse
}

func GetOrderIdAtTable(tableId string) *uint {
	var order Order

	if err := DB.Where("table_id = ? AND status = ?", tableId, "InProgress").
		Order("created_at DESC").
		First(&order).Error; err != nil {
		return nil
	}

	return &order.ID
}
