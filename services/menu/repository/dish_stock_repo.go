package repository

import (
	"time"
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/menu/model"
)

type DishStockRepo struct {
	db *gorm.DB
}

func NewDishStockRepo(db *gorm.DB) *DishStockRepo {
	return &DishStockRepo{db: db}
}

func (r *DishStockRepo) Create(stock *model.DishStock) error {
	return r.db.Create(stock).Error
}

func (r *DishStockRepo) FindByDishAndDate(dishID uint, date time.Time) (*model.DishStock, error) {
	var stock model.DishStock
	err := r.db.Where("dish_id = ? AND date = ?", dishID, date).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *DishStockRepo) UpdateSoldCount(dishID uint, date time.Time, count int) error {
	return r.db.Model(&model.DishStock{}).Where("dish_id = ? AND date = ?", dishID, date).
		UpdateColumn("sold_count", gorm.Expr("sold_count + ?", count)).Error
}

func (r *DishStockRepo) DeleteByDishID(dishID uint) error {
	return r.db.Where("dish_id = ?", dishID).Delete(&model.DishStock{}).Error
}
