package model
import "time"
type ActivityRegistration struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	ActivityID  uint      `gorm:"column:activity_id" json:"activity_id"`
	UserID      uint      `gorm:"column:user_id" json:"user_id"`
	Name        string    `gorm:"size:50;default:'';column:name" json:"name"`
	Phone       string    `gorm:"size:20;default:'';column:phone" json:"phone"`
	Remark      string    `gorm:"size:200;default:'';column:remark" json:"remark"`
	Status      string    `gorm:"size:20;default:'registered';column:status" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}
func (ActivityRegistration) TableName() string { return "activity_registrations" }
