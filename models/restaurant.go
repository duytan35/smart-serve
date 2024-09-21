package models

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateRestaurantInput struct {
	Name     string `json:"name" binding:"required" example:"Example Restaurant"`
	Phone    string `json:"phone" binding:"required" example:"1234567890"`
	Email    string `json:"email" binding:"required,email" example:"example@gmail.com"`
	Address  string `json:"address" binding:"required" example:"36 Pasteur, Ben Nghe, Quan 1, Ho Chi Minh City"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

type UpdateRestaurantInput struct {
	Name     string `json:"name" binding:"omitempty"`
	Phone    string `json:"phone" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email"`
	Address  string `json:"address" binding:"omitempty"`
	Password string `json:"password" binding:"omitempty,min=8"`
	Avatar   string `json:"avatar" binding:"omitempty"`
}

type RestaurantResponse struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Phone   string      `json:"phone"`
	Email   string      `json:"email"`
	Address string      `json:"address"`
	Avatar  string      `json:"avatar"`
	Steps   []OrderStep `json:"steps"`
}

type UpdateStepsInput struct {
	Steps []string `json:"steps" binding:"required"`
}

func GetRestaurants() []Restaurant {
	var restaurants []Restaurant
	DB.Model(&Restaurant{}).Find(&restaurants)

	return restaurants
}

func CreateRestaurant(restaurant Restaurant) (Restaurant, error) {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&restaurant).Error; err != nil {
			return err
		}

		if restaurant.ID == uuid.Nil {
			return fmt.Errorf("restaurant ID not generated")
		}

		defaultSteps := []OrderStep{
			{Step: 0, Name: "New", RestaurantID: restaurant.ID},
			{Step: 1, Name: "Confirmed", RestaurantID: restaurant.ID},
			{Step: 2, Name: "InProgress", RestaurantID: restaurant.ID},
			{Step: 3, Name: "Done", RestaurantID: restaurant.ID},
		}

		if err := CreateOrderSteps(defaultSteps); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return Restaurant{}, err
	}

	return restaurant, nil
}

func UpdateRestaurant(id string, restaurant UpdateRestaurantInput) (Restaurant, error) {
	var updatedRestaurant Restaurant
	if err := DB.First(&updatedRestaurant, id).Error; err != nil {
		return Restaurant{}, err
	}

	if updatedRestaurant.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedRestaurant.Password), bcrypt.DefaultCost)
		if err != nil {
			return Restaurant{}, err
		}
		updatedRestaurant.Password = string(hashedPassword)
	}

	if err := DB.Model(&updatedRestaurant).Updates(restaurant).Error; err != nil {
		return Restaurant{}, err
	}

	return updatedRestaurant, nil
}

func GetRestaurant(id string) (Restaurant, error) {
	var restaurant Restaurant

	if err := DB.Where("id = ?", id).
		Preload("Steps").
		First(&restaurant).Error; err != nil {
		return Restaurant{}, err
	}

	return restaurant, nil
}

func DeleteRestaurant(id string) error {
	var restaurant Restaurant
	if err := DB.Delete(&restaurant, id).Error; err != nil {
		return err
	}

	return nil
}

func GetRestaurantByEmail(email string) (Restaurant, error) {
	var restaurant Restaurant

	err := DB.Where(&Restaurant{Email: email}).First(&restaurant).Error

	if err != nil {
		return restaurant, err
	}

	return restaurant, nil
}

func RemoveAllOrderSteps(restaurantId string) error {
	return DB.Where("restaurant_id = ?", restaurantId).Delete(&OrderStep{}).Error
}

func CreateOrderSteps(steps []OrderStep) error {
	if err := DB.Create(&steps).Error; err != nil {
		return err
	}

	return nil
}
