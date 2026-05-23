package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/pricing/model"
)

type ComboRepo struct {
	db *gorm.DB
}

func NewComboRepo(db *gorm.DB) *ComboRepo {
	return &ComboRepo{db: db}
}

func (r *ComboRepo) Create(combo *model.Combo) error {
	return r.db.Create(combo).Error
}

func (r *ComboRepo) ListActive() ([]model.Combo, error) {
	var combos []model.Combo
	err := r.db.Where("status = ?", 1).Find(&combos).Error
	return combos, err
}
