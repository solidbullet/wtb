package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/order/model"
)

type OrderStatusLogRepo struct{ db *gorm.DB }

func NewOrderStatusLogRepo(db *gorm.DB) *OrderStatusLogRepo { return &OrderStatusLogRepo{db: db} }

func (r *OrderStatusLogRepo) DB() *gorm.DB { return r.db }

func (r *OrderStatusLogRepo) Create(log *model.OrderStatusLog) error {
	return r.db.Create(log).Error
}

func (r *OrderStatusLogRepo) ListByOrder(orderID uint) ([]model.OrderStatusLog, error) {
	var logs []model.OrderStatusLog
	err := r.db.Where("order_id = ?", orderID).Order("created_at asc").Find(&logs).Error
	return logs, err
}
