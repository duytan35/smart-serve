package models

func Migrate() {
	DB.AutoMigrate(&Restaurant{}, &User{})
}

type Restaurant struct {
	ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Users []User `json:"users" gorm:"foreignKey:RestaurantID"`
}

type User struct {
	ID           uint        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string      `json:"name" binding:"required"`
	Email        string      `json:"email" binding:"required,email"`
	Password     string      `json:"-" binding:"required"`
	RestaurantID uint        `json:"restaurantId" gorm:"not null"` // Foreign key
	Restaurant   *Restaurant `json:"restaurant"`
}
