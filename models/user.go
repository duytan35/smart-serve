package models

import (
	"errors"
	"fmt"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// omitempty is optional
type UpdateUserInput struct {
	Name     string `json:"name" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=8"`
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

func UpdateUser(id string, user UpdateUserInput) (User, error) {
	var updatedUser User
	if err := DB.First(&updatedUser, id).Error; err != nil {
		return User{}, err
	}
	if err := DB.Model(&updatedUser).Updates(user).Error; err != nil {
		return User{}, err
	}
	fmt.Println(updatedUser)
	return updatedUser, nil
}

func GetUser(id string) (User, error) {
	var user User
	if err := DB.First(&user, id).Error; err != nil {
		return User{}, err
	}

	return user, nil
}

func DeleteUser(id string) error {
	var user User
	if err := DB.First(&user, id).Error; err != nil {
		return errors.New("User not found")
	}

	if err := DB.Delete(&user, id).Error; err != nil {
		return errors.New("Internal server error")
	}

	return nil
}
