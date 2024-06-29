package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Migrate() {
	DB.AutoMigrate(
		&Restaurant{},
		&File{},
		&DishGroup{},
		&Dish{},
		&DishImage{},
		&Discount{},
		&DiscountDish{},
		&Table{},
		&Order{},
		&OrderDetail{},
	)
}

type Model struct {
	ID        uint      `json:"id" gorm:"primarykey,autoIncrement"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type Restaurant struct {
	Model
	ID       uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Name     string    `json:"name" gorm:"not null" binding:"required"`
	Phone    string    `json:"phone" gorm:"not null" binding:"required"`
	Email    string    `json:"email" gorm:"uniqueIndex;size:255;not null" binding:"required,email"`
	Address  string    `json:"address" gorm:"not null" binding:"required"`
	Password string    `json:"-" gorm:"not null" binding:"required,min=8"`
}

type File struct {
	Model
	ID           uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey"`
	RestaurantID uuid.UUID  `json:"restaurantId" gorm:"index;type:char(36);not null"`
	Name         string     `json:"name" gorm:"not null" binding:"required"`
	MineType     string     `json:"mineType" gorm:"not null" binding:"required"`
	Url          string     `json:"url" gorm:"not null" binding:"required"`
	Restaurant   Restaurant `json:"-" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE"`
}

// Need unique name in each restaurant
type DishGroup struct {
	Model
	RestaurantID uuid.UUID  `json:"restaurantId" gorm:"type:char(36);not null"`
	Name         string     `json:"name" gorm:"not null" binding:"required"`
	Restaurant   Restaurant `json:"-" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE"`
}

type Dish struct {
	Model
	RestaurantID uuid.UUID  `json:"restaurantId" gorm:"index;type:char(36);not null"`
	DishGroupID  uint       `json:"dishGroupId" gorm:"index;not null" binding:"required"`
	Name         string     `json:"name" gorm:"not null" binding:"required"`
	Description  string     `json:"description"` // optional
	Price        float64    `json:"price" gorm:"not null" binding:"required"`
	Status       uint       `json:"status" gorm:"type:TINYINT;not null;default:1"` // 0: inactive, 1: active
	Restaurant   Restaurant `json:"-" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE"`
	DishGroup    DishGroup  `json:"-" gorm:"foreignKey:DishGroupID;references:ID;constraint:OnDelete:CASCADE"`
}

type DishImage struct {
	Model
	DishID uint      `json:"dishId" gorm:"index;not null" binding:"required"`
	FileID uuid.UUID `json:"fileId" gorm:"index;not null" binding:"required"`
	Dish   Dish      `json:"-" gorm:"foreignKey:DishID;references:ID;constraint:OnDelete:CASCADE"`
	File   File      `json:"-" gorm:"foreignKey:FileID;references:ID;constraint:OnDelete:CASCADE"`
}

type Discount struct {
	Model
	RestaurantID uuid.UUID  `json:"restaurantId" gorm:"index;type:char(36);not null"`
	Percent      float64    `json:"percent" gorm:"not null" binding:"required"`
	Quantity     uint       `json:"quantity" gorm:"not null" binding:"required"`
	Restaurant   Restaurant `json:"-" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE"`
}

type DiscountDish struct {
	Model
	DiscountID uint      `json:"discountId" gorm:"index;not null" binding:"required"`
	DishID     uuid.UUID `json:"dishId" gorm:"index;not null" binding:"required"`
	Discount   Discount  `json:"-" gorm:"foreignKey:DiscountID;references:ID;constraint:OnDelete:CASCADE"`
	Dish       Dish      `json:"-" gorm:"foreignKey:DishID;references:ID;constraint:OnDelete:CASCADE"`
}

type Table struct {
	Model
	RestaurantID uuid.UUID  `json:"restaurantId" gorm:"index;type:char(36);not null"`
	Name         string     `json:"name" gorm:"not null" binding:"required"`
	Seats        uint       `json:"seats"` // optional
	Restaurant   Restaurant `json:"-" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE"`
}

type Order struct {
	Model
	TableID uint  `json:"tableId" gorm:"index;" binding:"required"`
	Status  uint  `json:"status" gorm:"type:TINYINT;not null;default:0"` // 0,1,2,3
	Table   Table `json:"-" gorm:"foreignKey:TableID;references:ID;constraint:OnDelete:SET NULL"`
}

type OrderDetail struct {
	Model
	OrderID         uint    `json:"orderId" gorm:"index;not null" binding:"required"`
	DishID          uint    `json:"dishId" gorm:"index;not null" binding:"required"`
	Quantity        uint    `json:"quantity" gorm:"not null" binding:"required"`
	Order           Order   `json:"-" gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE"`
	Dish            Dish    `json:"-" gorm:"foreignKey:DishID;references:ID;constraint:OnDelete:CASCADE"`
	DiscountPercent float64 `json:"discountPercent" gorm:"not null" binding:"required"` // value of current discount
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
