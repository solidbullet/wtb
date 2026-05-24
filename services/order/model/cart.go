package model

type CartItem struct {
	ID        uint   `gorm:"primaryKey" json:"-"`
	SeatID    string `gorm:"size:50;uniqueIndex:idx_seat_dish;not null" json:"seat_id"`
	DishID    uint   `gorm:"uniqueIndex:idx_seat_dish;not null" json:"dish_id"`
	DishName  string `gorm:"size:100" json:"dish_name"`
	Quantity  int    `gorm:"default:1" json:"quantity"`
	UnitPrice int    `json:"unit_price"`
	Remark    string `gorm:"size:500;default:''" json:"remark"`
}

func (CartItem) TableName() string { return "carts" }
