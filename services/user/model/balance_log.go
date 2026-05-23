package model

import "time"

type BalanceLog struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	UserID    uint      `gorm:"column:user_id" json:"user_id"`
	Type      string    `gorm:"size:20;column:type" json:"type"`
	Amount    int       `gorm:"column:amount" json:"amount"`
	OrderNo   string    `gorm:"size:64;default:'';column:order_no" json:"order_no"`
	Remark    string    `gorm:"size:255;default:'';column:remark" json:"remark"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (BalanceLog) TableName() string { return "balance_logs" }
