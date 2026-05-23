package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/pricing/model"
)

type RechargePlanRepo struct {
	db *gorm.DB
}

func NewRechargePlanRepo(db *gorm.DB) *RechargePlanRepo {
	return &RechargePlanRepo{db: db}
}

func (r *RechargePlanRepo) Create(plan *model.RechargePlan) error {
	return r.db.Create(plan).Error
}

func (r *RechargePlanRepo) Update(id uint, plan *model.RechargePlan) error {
	return r.db.Model(&model.RechargePlan{}).Where("id = ?", id).Updates(plan).Error
}

func (r *RechargePlanRepo) Delete(id uint) error {
	return r.db.Delete(&model.RechargePlan{}, id).Error
}

func (r *RechargePlanRepo) ListActive() ([]model.RechargePlan, error) {
	var plans []model.RechargePlan
	err := r.db.Where("status = ?", 1).Order("sort_order asc, created_at desc").Find(&plans).Error
	return plans, err
}

func (r *RechargePlanRepo) ListAll() ([]model.RechargePlan, error) {
	var plans []model.RechargePlan
	err := r.db.Order("sort_order asc, created_at desc").Find(&plans).Error
	return plans, err
}
