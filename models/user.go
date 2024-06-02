package models

type CreateUserInput struct {
	Name         string `json:"name" binding:"required" example:"Nguyen Van A"`
	Email        string `json:"email" binding:"required,email" example:"restaurant@gmail.com"`
	Password     string `json:"password" binding:"required,min=8" example:"12345678"`
	RestaurantID uint   `json:"restaurantId" binding:"required"` // Foreign key
}

// omitempty is optional
type UpdateUserInput struct {
	Name     string `json:"name" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=8"`
}

func GetUsers() []User {
	var users []User
	DB.Model(&User{}).Preload("Restaurant").Find(&users)
	return users
}

func CreateUser(user User) (User, error) {
	if err := DB.Create(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func UpdateUser(id string, user UpdateUserInput) (User, error) {
	var updatedUser User
	if err := DB.First(&updatedUser, id).Error; err != nil {
		return User{}, err
	}
	if err := DB.Model(&updatedUser).Updates(user).Error; err != nil {
		return User{}, err
	}
	return updatedUser, nil
}

func GetUser(id string) (User, error) {
	var user User
	if err := DB.Preload("Restaurant").First(&user, id).Error; err != nil {
		return User{}, err
	}

	return user, nil
}

func DeleteUser(id string) error {
	var user User
	if err := DB.Delete(&user, id).Error; err != nil {
		return err
	}

	return nil
}
