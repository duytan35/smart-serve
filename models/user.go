package models

import "fmt"

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" binding:"required,omitempty"`
	Email    string `json:"email" binding:"required,email,omitempty"`
	Password string `json:"password" binding:"required,min=8,omitempty"`
}

func GetUsers() []User {
	var users []User
	DB.Find(&users)
	return users
}

func CreateUser(user User) User {
	DB.Create(&user)
	return user
}

// error - should use patch, create new struct for update, error id return
func UpdateUser(id string, user User) (User, error) {
	var oldUser User
	if err := DB.First(&oldUser, id).Error; err != nil {
		return User{}, err
	}
	if err := DB.Model(&oldUser).Updates(user).Error; err != nil {
		return User{}, err
	}
	fmt.Println(user)
	return user, nil
}
