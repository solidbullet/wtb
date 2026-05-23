package model
import "time"
type ExchangeOrder struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	UserID      uint      `gorm:"column:user_id" json:"user_id"`
	GoodsID     uint      `gorm:"column:goods_id" json:"goods_id"`
	PointsCost  int       `gorm:"column:points_cost" json:"points_cost"`
	Status      string    `gorm:"size:20;default:'pending';column:status" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}
func (ExchangeOrder) TableName() string { return "exchange_orders" }
