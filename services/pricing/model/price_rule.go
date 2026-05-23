package model

type PriceRule struct {
	ID        uint   `gorm:"primaryKey;column:id" json:"id"`
	DishID    uint   `gorm:"column:dish_id" json:"dish_id"`
	RuleType  string `gorm:"size:20;column:rule_type" json:"rule_type"`
	Price     int    `gorm:"column:price" json:"price"`
	StartTime *string `gorm:"column:start_time" json:"start_time"`
	EndTime   *string `gorm:"column:end_time" json:"end_time"`
	Status    int16  `gorm:"default:1;column:status" json:"status"`
}

func (PriceRule) TableName() string { return "price_rules" }
