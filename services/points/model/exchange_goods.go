package model
type ExchangeGoods struct {
	ID           uint   `gorm:"primaryKey;column:id" json:"id"`
	Name         string `gorm:"size:100;column:name" json:"name"`
	Image        string `gorm:"size:500;default:'';column:image" json:"image"`
	PointsPrice  int    `gorm:"column:points_price" json:"points_price"`
	Stock        int    `gorm:"default:0;column:stock" json:"stock"`
	Type         string `gorm:"size:20;default:'physical';column:type" json:"type"`
	Status       int16  `gorm:"default:1;column:status" json:"status"`
}
func (ExchangeGoods) TableName() string { return "exchange_goods" }
