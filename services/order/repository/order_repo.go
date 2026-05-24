package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/order/model"
)

type OrderRepo struct{ db *gorm.DB }

func NewOrderRepo(db *gorm.DB) *OrderRepo { return &OrderRepo{db: db} }

func (r *OrderRepo) DB() *gorm.DB { return r.db }

func (r *OrderRepo) Create(o *model.Order) error {
	return r.db.Create(o).Error
}

func (r *OrderRepo) FindByID(id uint) (*model.Order, error) {
	var o model.Order
	err := r.db.Preload("Items").First(&o, id).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepo) FindByOrderNo(no string) (*model.Order, error) {
	var o model.Order
	err := r.db.Where("order_no = ?", no).First(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepo) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepo) ListByUser(userID uint, status string, page, pageSize int) ([]model.Order, int64, error) {
	var os []model.Order
	var total int64
	db := r.db.Model(&model.Order{}).Where("user_id = ?", userID)
	if status != "" {
		db = db.Where("status = ?", status)
	}
	db.Count(&total)
	query := r.db.Preload("Items").Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&os).Error
	return os, total, err
}

func (r *OrderRepo) ListAll(status string, page, pageSize int) ([]model.Order, int64, error) {
	var os []model.Order
	var total int64
	db := r.db.Model(&model.Order{})
	if status != "" {
		db = db.Where("status = ?", status)
	}
	db.Count(&total)
	query := r.db.Preload("Items")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&os).Error
	return os, total, err
}

func (r *OrderRepo) ListTodayPaidOrders() ([]model.Order, error) {
	var os []model.Order
	err := r.db.Preload("Items").
		Where("status IN ? AND DATE(created_at) = CURRENT_DATE", []string{"paid", "completed"}).
		Order("created_at desc").Find(&os).Error
	return os, err
}
