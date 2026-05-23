package model
import "time"
type RechargeOrder struct {
	ID            uint      `gorm:"primaryKey;column:id" json:"id"`
	UserID        uint      `gorm:"column:user_id" json:"user_id"`
	Amount        int       `gorm:"column:amount" json:"amount"`
	GiftedAmount  int       `gorm:"default:0;column:gifted_amount" json:"gifted_amount"`
	DiscountRate  float64   `gorm:"type:numeric(3,2);default:1.00;column:discount_rate" json:"discount_rate"`
	FinalAmount   int       `gorm:"column:final_amount" json:"final_amount"`
	Status        string    `gorm:"size:20;default:'pending';column:status" json:"status"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}
func (RechargeOrder) TableName() string { return "recharge_orders" }
