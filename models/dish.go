package models

import (
	"github.com/google/uuid"
)

type CreateDishInput struct {
	DishGroupID string  `json:"dishGroupId" binding:"required" example:"1"`
	Name        string  `json:"name" binding:"required" example:"Phở"`
	Description string  `json:"description" example:"Phở bò Việt Nam"`
	Price       float64 `json:"price" binding:"required" example:"50000"`
	Status      uint    `json:"status" example:"1"`

	ImageIds []uuid.UUID `json:"imageIds" binding:"required,minLen=1" example:"7c5a4b8f-fcf6-48d3-b21a-d81ebdfdf6f1,1566e532-72d4-49d4-8fca-c9142816006a"`
}

type UpdateDishInput struct {
	Name        string  `json:"name" binding:"required" example:"Phở"`
	Description string  `json:"description" example:"Phở bò Việt Nam"`
	Price       float64 `json:"price" binding:"required" example:"50000"`
	Status      uint    `json:"status" example:"1"`

	ImageIds []uuid.UUID `json:"imageIds" binding:"required,minLen=1" example:"7c5a4b8f-fcf6-48d3-b21a-d81ebdfdf6f1,1566e532-72d4-49d4-8fca-c9142816006a"`
}

func GetDishes(dishGroupId string) []Dish {
	var dishes []Dish
	DB.Preload("Images").Where("dish_group_id = ?", dishGroupId).Find(&dishes)

	return dishes
}

func CreateDish(dish Dish, images []uuid.UUID) (Dish, error) {
	tx := DB.Begin()

	if tx.Error != nil {
		return Dish{}, tx.Error
	}

	if err := tx.Create(&dish).Error; err != nil {
		tx.Rollback()
		return Dish{}, err
	}

	for _, imageId := range images {
		dishImage := DishImage{
			DishID: dish.ID,
			FileID: imageId,
		}

		if err := tx.Create(&dishImage).Error; err != nil {
			tx.Rollback()
			return Dish{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return Dish{}, err
	}

	return dish, nil
}

func UpdateDish(id string, restaurantId uuid.UUID, updateDishInput UpdateDishInput) (Dish, error) {
	var updatedDish Dish
	if err := DB.Joins("DishGroup").
		Where("dishes.id = ? AND DishGroup.restaurant_id = ?", id, restaurantId).
		First(&updatedDish).Error; err != nil {
		return Dish{}, err
	}

	dish := Dish{
		Name:        updateDishInput.Name,
		Description: updateDishInput.Description,
		Price:       updateDishInput.Price,
		Status:      updateDishInput.Status,
	}

	tx := DB.Begin()

	if tx.Error != nil {
		return Dish{}, tx.Error
	}

	if err := tx.Model(&updatedDish).Updates(dish).Error; err != nil {
		tx.Rollback()
		return Dish{}, err
	}

	if err := tx.Where("dish_id = ?", updatedDish.ID).Delete(&DishImage{}).Error; err != nil {
		tx.Rollback()
		return Dish{}, err
	}

	for _, imageId := range updateDishInput.ImageIds {
		dishImage := DishImage{
			DishID: updatedDish.ID,
			FileID: imageId,
		}

		if err := tx.Create(&dishImage).Error; err != nil {
			tx.Rollback()
			return Dish{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return Dish{}, err
	}

	return updatedDish, nil
}

func GetDish(id string) (Dish, error) {
	var dish Dish
	if err := DB.Preload("Images").Where("id = ?", id).First(&dish).Error; err != nil {
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
