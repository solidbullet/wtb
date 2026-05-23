package model

import "time"

type SeatStatusLog struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	SeatID    uint      `gorm:"column:seat_id" json:"seat_id"`
	OldStatus string    `gorm:"size:20;column:old_status" json:"old_status"`
	NewStatus string    `gorm:"size:20;column:new_status" json:"new_status"`
	OrderID   *uint     `gorm:"column:order_id" json:"order_id"`
	ChangedAt time.Time `gorm:"column:changed_at" json:"changed_at"`
}

func (SeatStatusLog) TableName() string { return "seat_status_logs" }
