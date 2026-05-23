package model
import "time"
type Announcement struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	Title       string    `gorm:"size:100;column:title" json:"title"`
	Content     string    `gorm:"type:text;column:content" json:"content"`
	Type        string    `gorm:"size:20;default:'text';column:type" json:"type"`
	Image       string    `gorm:"size:500;default:'';column:image" json:"image"`
	LinkType    string    `gorm:"size:20;default:'';column:link_type" json:"link_type"`
	LinkTarget  string    `gorm:"size:500;default:'';column:link_target" json:"link_target"`
	SortOrder   int       `gorm:"default:0;column:sort_order" json:"sort_order"`
	StartTime   time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime     time.Time `gorm:"column:end_time" json:"end_time"`
	Status      int16     `gorm:"default:1;column:status" json:"status"`
}
func (Announcement) TableName() string { return "announcements" }
