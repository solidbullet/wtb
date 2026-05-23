package model
import "time"
type OrderStatusLog struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	OrderID     uint      `gorm:"column:order_id" json:"order_id"`
	FromStatus  string    `gorm:"size:20;column:from_status" json:"from_status"`
	ToStatus    string    `gorm:"size:20;column:to_status" json:"to_status"`
	Operator    string    `gorm:"size:50;default:'';column:operator" json:"operator"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}
func (OrderStatusLog) TableName() string { return "order_status_logs" }
