package models

import (
	"golang.org/x/crypto/bcrypt"
)

func CreateFile(file File) (File, error) {
	if err := DB.Create(&file).Error; err != nil {
		return File{}, err
	}

	return file, nil
}

func UpdateFile(id string, restaurant UpdateRestaurantInput) (Restaurant, error) {
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

func GetFile(id string) (Restaurant, error) {
	var restaurant Restaurant
	if err := DB.First(&restaurant, id).Error; err != nil {
		return Restaurant{}, err
	}

	return restaurant, nil
}

func DeleteFile(id string) error {
	var restaurant Restaurant
	if err := DB.Delete(&restaurant, id).Error; err != nil {
		return err
	}

	return nil
}
