package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/user/model"
)

type RechargeRepo struct {
	db *gorm.DB
}

func NewRechargeRepo(db *gorm.DB) *RechargeRepo {
	return &RechargeRepo{db: db}
}

func (r *RechargeRepo) Create(record *model.RechargeRecord) error {
	return r.db.Create(record).Error
}

func (r *RechargeRepo) ListByUserID(userID uint, page, pageSize int) ([]model.RechargeRecord, int64, error) {
	var records []model.RechargeRecord
	var total int64
	r.db.Model(&model.RechargeRecord{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records).Error
	return records, total, err
}

func (r *RechargeRepo) GetTotalAmountByUserID(userID uint) (int, error) {
	var total int
	err := r.db.Model(&model.RechargeRecord{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}
