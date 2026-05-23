package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/user/model"
)

type ConsumptionRepo struct {
	db *gorm.DB
}

func NewConsumptionRepo(db *gorm.DB) *ConsumptionRepo {
	return &ConsumptionRepo{db: db}
}

func (r *ConsumptionRepo) Create(record *model.ConsumptionRecord) error {
	return r.db.Create(record).Error
}

func (r *ConsumptionRepo) ListByUserID(userID uint, page, pageSize int) ([]model.ConsumptionRecord, int64, error) {
	var records []model.ConsumptionRecord
	var total int64
	r.db.Model(&model.ConsumptionRecord{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records).Error
	return records, total, err
}

func (r *ConsumptionRepo) SummaryByUserID(userID uint) (map[string]interface{}, error) {
	var totalAmount int
	var totalOrders int64
	var thisMonthAmount int
	var thisMonthOrders int64

	r.db.Model(&model.ConsumptionRecord{}).Where("user_id = ?", userID).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)
	r.db.Model(&model.ConsumptionRecord{}).Where("user_id = ?", userID).
		Count(&totalOrders)

	r.db.Model(&model.ConsumptionRecord{}).Where("user_id = ? AND created_at >= date_trunc('month', now())", userID).
		Select("COALESCE(SUM(amount), 0)").Scan(&thisMonthAmount)
	r.db.Model(&model.ConsumptionRecord{}).Where("user_id = ? AND created_at >= date_trunc('month', now())", userID).
		Count(&thisMonthOrders)

	avgAmount := 0
	if totalOrders > 0 {
		avgAmount = totalAmount / int(totalOrders)
	}

	return map[string]interface{}{
		"total_amount":      totalAmount,
		"total_orders":      totalOrders,
		"avg_amount":        avgAmount,
		"this_month_amount": thisMonthAmount,
		"this_month_orders": thisMonthOrders,
	}, nil
}
