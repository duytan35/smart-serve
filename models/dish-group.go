package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type DishGroupInput struct {
	Name string `json:"name" binding:"required" example:"Noodles"`
}

type DishGroupResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	RestaurantID string    `json:"restaurantId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func GetDishGroups(restaurantId string) []DishGroup {
	var dishGroups []DishGroup
	DB.
		Preload("Dishes").
		Where("restaurant_id = ?", restaurantId).Find(&dishGroups)

	return dishGroups
}

func CreateDishGroup(dishGroup DishGroup) (DishGroup, error) {
	if err := DB.Create(&dishGroup).Error; err != nil {
		return DishGroup{}, err
	}

	return dishGroup, nil
}

func UpdateDishGroup(id string, restaurantId uuid.UUID, dishGroup DishGroupInput) (DishGroup, error) {
	var updatedDishGroup DishGroup
	if err := DB.Where("id = ?", id).First(&updatedDishGroup).Error; err != nil {
		return DishGroup{}, err
	}

	if updatedDishGroup.RestaurantID != restaurantId {
		return DishGroup{}, errors.New("record not found")
	}

	if err := DB.Model(&updatedDishGroup).Updates(dishGroup).Error; err != nil {
		return DishGroup{}, err
	}

	return updatedDishGroup, nil
}

func GetDishGroup(id string) (DishGroup, error) {
	var dishGroup DishGroup
	if err := DB.Preload("Dishes").Where("id = ?", id).First(&dishGroup).Error; err != nil {
		return DishGroup{}, err
	}

	return dishGroup, nil
}

func DeleteDishGroup(id string, restaurantId uuid.UUID) error {
	var dishGroup DishGroup

	if err := DB.First(&dishGroup, id).Error; err != nil {
		return err
	}

	if dishGroup.RestaurantID != restaurantId {
		return errors.New("record not found")
	}

	if err := DB.Delete(&dishGroup, id).Error; err != nil {
		return err
	}

	return nil
}
