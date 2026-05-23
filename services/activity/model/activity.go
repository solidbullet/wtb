package model
import "time"
type Activity struct {
	ID                   uint      `gorm:"primaryKey;column:id" json:"id"`
	Title                string    `gorm:"size:100;column:title" json:"title"`
	Description          string    `gorm:"type:text;column:description" json:"description"`
	Image                string    `gorm:"size:500;default:'';column:image" json:"image"`
	MaxParticipants      int       `gorm:"default:-1;column:max_participants" json:"max_participants"`
	CurrentParticipants  int       `gorm:"default:0;column:current_participants" json:"current_participants"`
	EventTime            *time.Time `gorm:"column:event_time" json:"event_time"`
	Location             string    `gorm:"size:200;default:'';column:location" json:"location"`
	Status               string    `gorm:"size:20;default:'draft';column:status" json:"status"`
	CreatedAt            time.Time `gorm:"column:created_at" json:"created_at"`
}
func (Activity) TableName() string { return "activities" }
