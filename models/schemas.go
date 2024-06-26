package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Migrate() {
	DB.AutoMigrate(&Restaurant{}, &File{})
}

type Model struct {
	ID        uint      `json:"id" gorm:"primarykey,autoIncrement"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type Restaurant struct {
	Model
	ID       uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Name     string    `json:"name" binding:"required"`
	Phone    string    `json:"phone" binding:"required"`
	Email    string    `json:"email" gorm:"uniqueIndex;size:255" binding:"required,email"`
	Address  string    `json:"address" binding:"required"`
	Password string    `json:"-" binding:"required,min=8"`
}

type File struct {
	Model
	ID           uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey"`
	RestaurantID uuid.UUID  `json:"restaurantId" gorm:"index;type:char(36);not null"`
	Name         string     `json:"name" binding:"required"`
	MineType     string     `json:"mineType" binding:"required"`
	Url          string     `json:"url" binding:"required"`
	Restaurant   Restaurant `json:"-" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE"`
}

// apply only create
func (r *Restaurant) BeforeSave(tx *gorm.DB) (err error) {
	r.ID = uuid.New()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	r.Password = string(hashedPassword)

	return nil
}

func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	f.ID = uuid.New()
	return nil
}
