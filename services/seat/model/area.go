package model

import "time"

type Area struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	Name      string    `gorm:"size:50;column:name" json:"name"`
	SortOrder int       `gorm:"default:0;column:sort_order" json:"sort_order"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Area) TableName() string { return "areas" }
