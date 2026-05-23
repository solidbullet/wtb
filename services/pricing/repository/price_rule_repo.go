package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/pricing/model"
)

type PriceRuleRepo struct {
	db *gorm.DB
}

func NewPriceRuleRepo(db *gorm.DB) *PriceRuleRepo {
	return &PriceRuleRepo{db: db}
}

func (r *PriceRuleRepo) Create(rule *model.PriceRule) error {
	return r.db.Create(rule).Error
}

func (r *PriceRuleRepo) ListByDishID(dishID uint) ([]model.PriceRule, error) {
	var rules []model.PriceRule
	err := r.db.Where("dish_id = ? AND status = ?", dishID, 1).Find(&rules).Error
	return rules, err
}
