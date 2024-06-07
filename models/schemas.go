package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Migrate() {
	DB.AutoMigrate(&Restaurant{}, &User{})
}

type Restaurant struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" gorm:"unique" binding:"required"`
	Email   string `json:"email" gorm:"unique" binding:"required,email"`
	Address string `json:"address" binding:"required"`
	Users   []User `json:"users" gorm:"foreignKey:RestaurantID"`
}

type User struct {
	ID           uint        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string      `json:"name" binding:"required"`
	Email        string      `json:"email" gorm:"unique" binding:"required,email"`
	Password     string      `json:"-" binding:"required"`
	RestaurantID uint        `json:"restaurantId" gorm:"not null"` // Foreign key
	Restaurant   *Restaurant `json:"restaurant"`
}

// apply only create user
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}
