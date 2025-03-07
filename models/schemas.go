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
		&OrderStep{},
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
	Avatar   string    `json:"avatar"`

	DishGroup []DishGroup `json:"dishGroups" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE"`
	Steps     []OrderStep `json:"steps" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE"`
}

type File struct {
	Model
	ID           uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey"`
	RestaurantID uuid.UUID  `json:"restaurantId" gorm:"index;type:char(36);not null"`
	Name         string     `json:"name" gorm:"not null" binding:"required"`
	MineType     string     `json:"mineType" gorm:"not null" binding:"required"`
	Restaurant   Restaurant `json:"-" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE"`
}

// Need unique name in each restaurant
type DishGroup struct {
	Model
	RestaurantID uuid.UUID  `json:"restaurantId" gorm:"uniqueIndex:idx_name_restaurant;type:char(36);not null"`
	Name         string     `json:"name" gorm:"uniqueIndex:idx_name_restaurant;type:char(255);not null" binding:"required"`
	Restaurant   Restaurant `json:"-" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE"`

	Dishes []Dish `json:"dishes" gorm:"foreignKey:DishGroupID;references:ID;constraint:OnDelete:CASCADE"`
}

type Dish struct {
	Model
	DishGroupID uint      `json:"dishGroupId" gorm:"index;not null" binding:"required"`
	Name        string    `json:"name" gorm:"not null" binding:"required"`
	Description string    `json:"description"` // optional
	Price       float64   `json:"price" gorm:"not null" binding:"required"`
	Status      uint      `json:"status" gorm:"TINYINT;not null;default 1"` // 0: inactive, 1: active
	DishGroup   DishGroup `json:"-" gorm:"foreignKey:DishGroupID;references:ID;constraint:OnDelete:CASCADE"`

	Images   []DishImage `json:"-" gorm:"foreignKey:DishID;references:ID;constraint:OnDelete:CASCADE"`
	ImageIds []uuid.UUID `json:"imageIds" gorm:"-"`
}

type DishImage struct {
	Model
	DishID uint      `json:"dishId" gorm:"index;not null" binding:"required"`
	FileID uuid.UUID `json:"fileId" gorm:"index;not null" binding:"required"`
	Dish   Dish      `json:"-" gorm:"foreignKey:DishID;references:ID;constraint:OnDelete:CASCADE"`
	File   File      `json:"-" gorm:"foreignKey:FileID;references:ID"`
}

type Discount struct {
	Model
	RestaurantID uuid.UUID  `json:"restaurantId" gorm:"index;type:char(36);not null"`
	Percent      float64    `json:"percent" gorm:"not null" binding:"required"`
	Quantity     uint       `json:"quantity" gorm:"not null" binding:"required"`
	Restaurant   Restaurant `json:"-" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE"`
}

type DiscountDish struct { // apply multiple discount for a dish
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
	TableID uint        `json:"tableId" gorm:"index;" binding:"required"`
	Status  OrderStatus `json:"status" gorm:"not null;default:'InProgress'"`
	Table   Table       `json:"-" gorm:"foreignKey:TableID;references:ID;constraint:OnDelete:SET NULL"`

	OrderDetails []OrderDetail `json:"orderDetails" gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE"`
}

type OrderDetail struct {
	Model
	OrderID          uint    `json:"orderId" gorm:"index;not null" binding:"required"`
	DishID           uint    `json:"dishId" gorm:"index;not null" binding:"required"`
	Quantity         uint    `json:"quantity" gorm:"not null" binding:"required"`
	Step             uint    `json:"step" gorm:"type:TINYINT;not null;default:0"`                  // 0,1,2,3
	DiscountPercent  float64 `json:"discountPercent" gorm:"not null;default 0" binding:"required"` // value of current discount
	Order            Order   `json:"-" gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE"`
	Dish             Dish    `json:"dish" gorm:"foreignKey:DishID;references:ID;constraint:OnDelete:CASCADE"`
	Note             string  `json:"note"`
	GroupOrderNumber uint    `json:"groupOrderNumber" gorm:"not null;default:1" binding:"required"`
}

type OrderStep struct {
	Model
	RestaurantID uuid.UUID  `json:"restaurantId" gorm:"index;type:char(36);not null"`
	Name         string     `json:"name" gorm:"not null" binding:"required"`
	Step         uint       `json:"step" gorm:"not null" binding:"required"`
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

func (d *Dish) AfterFind(tx *gorm.DB) (err error) {
	ImageIds := []uuid.UUID{}
	for _, image := range d.Images {
		ImageIds = append(ImageIds, image.FileID)
	}
	d.ImageIds = ImageIds

	return
}
