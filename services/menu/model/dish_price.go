package model

import "time"

type DishPrice struct {
	ID        uint       `gorm:"primaryKey;column:id" json:"id"`
	DishID    uint       `gorm:"column:dish_id" json:"dish_id"`
	PriceType string     `gorm:"size:20;column:price_type" json:"price_type"`
	Price     int        `gorm:"column:price" json:"price"`
	StartTime *time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime   *time.Time `gorm:"column:end_time" json:"end_time"`
}

func (DishPrice) TableName() string { return "dish_prices" }
