package model

type Combo struct {
	ID       uint   `gorm:"primaryKey;column:id" json:"id"`
	Name     string `gorm:"size:100;column:name" json:"name"`
	Price    int    `gorm:"column:price" json:"price"`
	DishList string `gorm:"type:text;column:dish_list" json:"dish_list"`
	Status   int16  `gorm:"default:1;column:status" json:"status"`
}

func (Combo) TableName() string { return "combos" }
