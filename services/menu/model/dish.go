package model

import "time"

type Dish struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	CategoryID  uint      `gorm:"column:category_id" json:"category_id"`
	Name        string    `gorm:"size:100;column:name" json:"name"`
	Subtitle    string    `gorm:"size:200;default:'';column:subtitle" json:"subtitle"`
	Description string    `gorm:"type:text;column:description" json:"description"`
	Images      string    `gorm:"type:text;column:images" json:"images"`
	Tags        string    `gorm:"size:200;default:'';column:tags" json:"tags"`
	Status      int16     `gorm:"default:1;column:status" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Dish) TableName() string { return "dishes" }
