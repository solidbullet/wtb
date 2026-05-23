package model

import "time"

type Seat struct {
	ID         uint      `gorm:"primaryKey;column:id" json:"id"`
	AreaID     uint      `gorm:"column:area_id" json:"area_id"`
	Name       string    `gorm:"size:50;column:name" json:"name"`
	Type       string    `gorm:"size:20;default:'normal';column:type" json:"type"`
	Capacity   int       `gorm:"default:4;column:capacity" json:"capacity"`
	QrcodeURL  string    `gorm:"size:500;default:'';column:qrcode_url" json:"qrcode_url"`
	Status     string    `gorm:"size:20;default:'available';column:status" json:"status"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Seat) TableName() string { return "seats" }
