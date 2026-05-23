package model
type OrderItem struct {
	ID         uint   `gorm:"primaryKey;column:id" json:"id"`
	OrderID    uint   `gorm:"column:order_id" json:"order_id"`
	DishID     uint   `gorm:"column:dish_id" json:"dish_id"`
	DishName   string `gorm:"size:100;column:dish_name" json:"dish_name"`
	Quantity   int    `gorm:"column:quantity" json:"quantity"`
	UnitPrice  int    `gorm:"column:unit_price" json:"unit_price"`
}
func (OrderItem) TableName() string { return "order_items" }
