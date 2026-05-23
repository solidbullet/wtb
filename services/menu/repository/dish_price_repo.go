package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/menu/model"
)

type DishPriceRepo struct {
	db *gorm.DB
}

func NewDishPriceRepo(db *gorm.DB) *DishPriceRepo {
	return &DishPriceRepo{db: db}
}

func (r *DishPriceRepo) Create(dp *model.DishPrice) error {
	return r.db.Create(dp).Error
}

func (r *DishPriceRepo) ListByDishID(dishID uint) ([]model.DishPrice, error) {
	var prices []model.DishPrice
	err := r.db.Where("dish_id = ?", dishID).Find(&prices).Error
	return prices, err
}

func (r *DishPriceRepo) DeleteByDishID(dishID uint) error {
	return r.db.Where("dish_id = ?", dishID).Delete(&model.DishPrice{}).Error
}
