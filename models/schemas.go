package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Migrate() {
	DB.AutoMigrate(&Restaurant{})
}

type Restaurant struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" gorm:"uniqueIndex;size:255" binding:"required,email"`
	Address  string `json:"address" binding:"required"`
	Password string `json:"-" binding:"required,min=8"`
}

// apply only create
func (u *Restaurant) BeforeSave(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}
