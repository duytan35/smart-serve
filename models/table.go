package models

import (
	"smart-serve/utils"

	"github.com/google/uuid"
)

type TableInput struct {
	Name  string `json:"name" binding:"required" example:"Bàn 1"`
	Seats uint   `json:"seats" example:"4"`
}

func GetTables(restaurantId string) []Table {
	var tables []Table
	DB.Model(&Table{}).Where("restaurant_id = ?", restaurantId).Find(&tables)

	return tables
}

func CreateTable(table Table) (Table, error) {
	if err := DB.Create(&table).Error; err != nil {
		return Table{}, err
	}

	return table, nil
}

func UpdateTable(id string, restaurantId uuid.UUID, table TableInput) (Table, error) {
	var updatedTable Table
	if err := DB.Where("id = ? AND restaurant_id = ?", id, restaurantId).
		First(&updatedTable).Error; err != nil {
		return Table{}, err
	}

	if err := DB.Model(&updatedTable).Updates(utils.ToMap(table)).Error; err != nil {
		return Table{}, err
	}

	return updatedTable, nil
}

func GetTable(id string, restaurantId string) (Table, error) {
	var dish Table
	if err := DB.Where("id = ? AND restaurant_id = ?", id, restaurantId).First(&dish).Error; err != nil {
		return Table{}, err
	}

	return dish, nil
}

func GetTableById(id string) (Table, error) {
	var dish Table
	if err := DB.Where("id = ?", id).First(&dish).Error; err != nil {
		return Table{}, err
	}

	return dish, nil
}

func DeleteTable(id string, restaurantId uuid.UUID) error {
	var table Table

	if err := DB.Where("id = ? AND restaurant_id = ?", id, restaurantId).
		First(&table).Error; err != nil {
		return err
	}

	if err := DB.Delete(&table, id).Error; err != nil {
		return err
	}

	return nil
}
