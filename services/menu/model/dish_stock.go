package model

import "time"

type DishStock struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	DishID      uint      `gorm:"column:dish_id" json:"dish_id"`
	DailyLimit  int       `gorm:"default:-1;column:daily_limit" json:"daily_limit"`
	SoldCount   int       `gorm:"default:0;column:sold_count" json:"sold_count"`
	Date        time.Time `gorm:"type:date;column:date" json:"date"`
}

func (DishStock) TableName() string { return "dish_stocks" }
