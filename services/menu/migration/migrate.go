package migration

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/menu/model"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.Category{}, &model.Dish{}, &model.DishPrice{}, &model.DishStock{})
}
