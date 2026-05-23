package model

import "time"

type Promotion struct {
	ID         uint      `gorm:"primaryKey;column:id" json:"id"`
	Name       string    `gorm:"size:100;column:name" json:"name"`
	Type       string    `gorm:"size:30;column:type" json:"type"`
	ConfigJSON string    `gorm:"type:text;column:config_json" json:"config_json"`
	StartTime  time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime    time.Time `gorm:"column:end_time" json:"end_time"`
	Status     int16     `gorm:"default:1;column:status" json:"status"`
}

func (Promotion) TableName() string { return "promotions" }
