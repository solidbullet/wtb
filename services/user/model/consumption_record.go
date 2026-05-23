package model

import "time"

type ConsumptionRecord struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	UserID    uint      `gorm:"column:user_id" json:"user_id"`
	OrderID   uint      `gorm:"column:order_id" json:"order_id"`
	Amount    int       `gorm:"column:amount" json:"amount"`
	DishCount int       `gorm:"default:0;column:dish_count" json:"dish_count"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (ConsumptionRecord) TableName() string { return "consumption_records" }
