package model

import "time"

type PetProfile struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	UserID    uint      `gorm:"column:user_id" json:"user_id"`
	Name      string    `gorm:"size:50;column:name" json:"name"`
	Breed     string    `gorm:"size:50;default:'';column:breed" json:"breed"`
	Weight    float64   `gorm:"type:numeric(5,2);default:0;column:weight" json:"weight"`
	Birthday  *string   `gorm:"type:date;column:birthday" json:"birthday"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (PetProfile) TableName() string { return "pet_profiles" }
