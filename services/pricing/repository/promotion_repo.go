package repository

import (
	"time"
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/pricing/model"
)

type PromotionRepo struct {
	db *gorm.DB
}

func NewPromotionRepo(db *gorm.DB) *PromotionRepo {
	return &PromotionRepo{db: db}
}

func (r *PromotionRepo) Create(promo *model.Promotion) error {
	return r.db.Create(promo).Error
}

func (r *PromotionRepo) ListActive() ([]model.Promotion, error) {
	var promos []model.Promotion
	now := time.Now()
	err := r.db.Where("status = ? AND start_time <= ? AND end_time >= ?", 1, now, now).Find(&promos).Error
	return promos, err
}
