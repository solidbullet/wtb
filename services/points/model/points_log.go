package model
import "time"
type PointsLog struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	UserID    uint      `gorm:"column:user_id" json:"user_id"`
	Type      string    `gorm:"size:20;column:type" json:"type"`
	Points    int       `gorm:"column:points" json:"points"`
	SourceID  string    `gorm:"size:64;default:'';column:source_id" json:"source_id"`
	Remark    string    `gorm:"size:200;default:'';column:remark" json:"remark"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}
func (PointsLog) TableName() string { return "points_logs" }
