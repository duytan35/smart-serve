package models

type CreateRestaurantInput struct {
	Name    string `json:"name" binding:"required" example:"Example Restaurant"`
	Phone   string `json:"phone" binding:"required,phone" example:"1234567890"`
	Email   string `json:"email" binding:"required,email" example:"example@gmail.com"`
	Address string `json:"address" binding:"required" example:"36 Pasteur, Ben Nghe, Quan 1, Ho Chi Minh City"`
}

type UpdateRestaurantInput struct {
	Name    string `json:"name" binding:"omitempty"`
	Phone   string `json:"phone" binding:"omitempty"`
	Email   string `json:"email" binding:"omitempty,email"`
	Address string `json:"address" binding:"omitempty"`
}

func GetRestaurants() []Restaurant {
	var restaurants []Restaurant
	DB.Model(&Restaurant{}).Find(&restaurants)

	return restaurants
}

func CreateRestaurant(restaurant Restaurant) (Restaurant, error) {
	if err := DB.Create(&restaurant).Error; err != nil {
		return Restaurant{}, err
	}

	return restaurant, nil
}

func UpdateRestaurant(id string, restaurant UpdateRestaurantInput) (Restaurant, error) {
	var updatedRestaurant Restaurant
	if err := DB.First(&updatedRestaurant, id).Error; err != nil {
		return Restaurant{}, err
	}
	if err := DB.Model(&updatedRestaurant).Updates(restaurant).Error; err != nil {
		return Restaurant{}, err
	}
	return updatedRestaurant, nil
}

func GetRestaurant(id string) (Restaurant, error) {
	var restaurant Restaurant
	if err := DB.Preload("Users").First(&restaurant, id).Error; err != nil {
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
