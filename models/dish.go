package models

import (
	"smart-serve/utils"

	"github.com/google/uuid"
)

type CreateDishInput struct {
	DishGroupID string  `json:"dishGroupId" binding:"required" example:"1"`
	Name        string  `json:"name" binding:"required" example:"Phở"`
	Description string  `json:"description" example:"Phở bò Việt Nam"`
	Price       float64 `json:"price" binding:"required" example:"50000"`
	Status      uint    `json:"status" example:"1"`
}

type UpdateDishInput struct {
	Name        string  `json:"name" binding:"required" example:"Phở"`
	Description string  `json:"description" example:"Phở bò Việt Nam"`
	Price       float64 `json:"price" binding:"required" example:"50000"`
	Status      uint    `json:"status" example:"1"`
}

func GetDishes(dishGroupId string) []Dish {
	var dishes []Dish
	DB.Model(&Dish{}).Where("dish_group_id = ?", dishGroupId).Find(&dishes)

	return dishes
}

func CreateDish(dish Dish) (Dish, error) {
	if err := DB.Create(&dish).Error; err != nil {
		return Dish{}, err
	}

	return dish, nil
}

func UpdateDish(id string, restaurantId uuid.UUID, dish UpdateDishInput) (Dish, error) {
	var updatedDish Dish
	if err := DB.Joins("DishGroup").
		Where("dishes.id = ? AND DishGroup.restaurant_id = ?", id, restaurantId).
		First(&updatedDish).Error; err != nil {
		return Dish{}, err
	}

	if err := DB.Model(&updatedDish).Updates(utils.ToMap(dish)).Error; err != nil {
		return Dish{}, err
	}

	return updatedDish, nil
}

func GetDish(id string) (Dish, error) {
	var dish Dish
	if err := DB.Where("id = ?", id).First(&dish).Error; err != nil {
		return Dish{}, err
	}

	return dish, nil
}

func DeleteDish(id string, restaurantId uuid.UUID) error {
	var dish Dish

	if err := DB.Joins("DishGroup").
		Where("dishes.id = ? AND DishGroup.restaurant_id = ?", id, restaurantId).
		First(&dish).Error; err != nil {
		return err
	}

	if err := DB.Delete(&dish, id).Error; err != nil {
		return err
	}

	return nil
}
