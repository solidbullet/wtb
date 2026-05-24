package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"github.com/wtb-ordering/services/order/model"
)

type CartRepo struct {
	db *gorm.DB
}

func NewCartRepo(db *gorm.DB) *CartRepo {
	return &CartRepo{db: db}
}

func (r *CartRepo) Upsert(seatID string, item model.CartItem) error {
	item.SeatID = seatID
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "seat_id"}, {Name: "dish_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity", "unit_price", "dish_name", "remark"}),
	}).Create(&item).Error
}

func (r *CartRepo) ListBySeat(seatID string) ([]model.CartItem, error) {
	var items []model.CartItem
	err := r.db.Where("seat_id = ?", seatID).Find(&items).Error
	return items, err
}

func (r *CartRepo) Remove(seatID string, dishID uint) error {
	return r.db.Where("seat_id = ? AND dish_id = ?", seatID, dishID).Delete(&model.CartItem{}).Error
}

func (r *CartRepo) Clear(seatID string) error {
	return r.db.Where("seat_id = ?", seatID).Delete(&model.CartItem{}).Error
}
