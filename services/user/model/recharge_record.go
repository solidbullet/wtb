package model

import "time"

type RechargeRecord struct {
	ID           uint      `gorm:"primaryKey;column:id" json:"id"`
	UserID       uint      `gorm:"column:user_id" json:"user_id"`
	Amount       int       `gorm:"column:amount" json:"amount"`
	GiftedAmount int       `gorm:"default:0;column:gifted_amount" json:"gifted_amount"`
	Channel      string    `gorm:"size:20;default:'wxpay';column:channel" json:"channel"`
	Status       string    `gorm:"size:20;default:'pending';column:status" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (RechargeRecord) TableName() string { return "recharge_records" }
