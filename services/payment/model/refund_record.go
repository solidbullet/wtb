package model
import "time"
type RefundRecord struct {
	ID             uint      `gorm:"primaryKey;column:id" json:"id"`
	PaymentOrderID uint      `gorm:"column:payment_order_id" json:"payment_order_id"`
	RefundNo       string    `gorm:"size:32;uniqueIndex;column:refund_no" json:"refund_no"`
	Amount         int       `gorm:"column:amount" json:"amount"`
	Reason         string    `gorm:"size:200;default:'';column:reason" json:"reason"`
	Status         string    `gorm:"size:20;default:'pending';column:status" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
}
func (RefundRecord) TableName() string { return "refund_records" }
