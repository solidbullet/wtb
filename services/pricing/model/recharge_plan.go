package model

import "time"

type RechargePlan struct {
	ID           uint      `gorm:"primaryKey;column:id" json:"id"`
	Name         string    `gorm:"size:100;column:name" json:"name"`
	Amount       int       `gorm:"column:amount" json:"amount"`
	FinalAmount  int       `gorm:"column:final_amount" json:"final_amount"`
	GiftAmount   int       `gorm:"column:gift_amount" json:"gift_amount"`
	SortOrder    int       `gorm:"default:0;column:sort_order" json:"sort_order"`
	Status       int16     `gorm:"default:1;column:status" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (RechargePlan) TableName() string { return "recharge_plans" }
