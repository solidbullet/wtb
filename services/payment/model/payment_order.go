package model
import "time"
type PaymentOrder struct {
	ID           uint      `gorm:"primaryKey;column:id" json:"id"`
	OrderNo      string    `gorm:"size:32;column:order_no" json:"order_no"`
	OutTradeNo   string    `gorm:"size:32;uniqueIndex;column:out_trade_no" json:"out_trade_no"`
	UserID       uint      `gorm:"column:user_id" json:"user_id"`
	Amount       int       `gorm:"column:amount" json:"amount"`
	Channel      string    `gorm:"size:20;column:channel" json:"channel"`
	Status       string    `gorm:"size:20;default:'pending';column:status" json:"status"`
	WxPrepayID   string    `gorm:"size:64;default:'';column:wx_prepay_id" json:"wx_prepay_id"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}
func (PaymentOrder) TableName() string { return "payment_orders" }
