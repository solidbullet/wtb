package model
import "time"
type Order struct {
	ID              uint         `gorm:"primaryKey;column:id" json:"id"`
	OrderNo         string       `gorm:"size:32;uniqueIndex;column:order_no" json:"order_no"`
	SeatID          string       `gorm:"size:50;column:seat_id" json:"seat_id"`
	UserID          uint         `gorm:"column:user_id" json:"user_id"`
	Status          string       `gorm:"size:20;default:'pending';column:status" json:"status"`
	TotalAmount     int          `gorm:"column:total_amount" json:"total_amount"`
	DiscountAmount  int          `gorm:"default:0;column:discount_amount" json:"discount_amount"`
	PayAmount       int          `gorm:"column:pay_amount" json:"pay_amount"`
	Remark          string       `gorm:"size:500;default:'';column:remark" json:"remark"`
	CreatedAt       time.Time    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time    `gorm:"column:updated_at" json:"updated_at"`
	Items           []OrderItem  `gorm:"foreignKey:OrderID" json:"items"`
}
func (Order) TableName() string { return "orders" }
