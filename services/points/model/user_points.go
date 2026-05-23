package model
import "time"
type UserPoints struct {
	ID            uint      `gorm:"primaryKey;column:id" json:"id"`
	UserID        uint      `gorm:"uniqueIndex;column:user_id" json:"user_id"`
	TotalPoints   int       `gorm:"default:0;column:total_points" json:"total_points"`
	UsedPoints    int       `gorm:"default:0;column:used_points" json:"used_points"`
	FrozenPoints  int       `gorm:"default:0;column:frozen_points" json:"frozen_points"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}
func (UserPoints) TableName() string { return "user_points" }
