package model
import "time"
type PaymentRecord struct {
	ID             uint      `gorm:"primaryKey;column:id" json:"id"`
	PaymentOrderID uint      `gorm:"column:payment_order_id" json:"payment_order_id"`
	Channel        string    `gorm:"size:20;column:channel" json:"channel"`
	Amount         int       `gorm:"column:amount" json:"amount"`
	TransactionID  string    `gorm:"size:64;default:'';column:transaction_id" json:"transaction_id"`
	PaidAt         *time.Time `gorm:"column:paid_at" json:"paid_at"`
}
func (PaymentRecord) TableName() string { return "payment_records" }
