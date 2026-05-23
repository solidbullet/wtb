package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/user/model"
)

type BalanceLogRepo struct {
	db *gorm.DB
}

func NewBalanceLogRepo(db *gorm.DB) *BalanceLogRepo {
	return &BalanceLogRepo{db: db}
}

func (r *BalanceLogRepo) Create(log *model.BalanceLog) error {
	return r.db.Create(log).Error
}

func (r *BalanceLogRepo) ListByUserID(userID uint, page, pageSize int) ([]model.BalanceLog, int64, error) {
	var logs []model.BalanceLog
	var total int64
	r.db.Model(&model.BalanceLog{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}
