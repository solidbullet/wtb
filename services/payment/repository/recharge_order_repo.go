package repository
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/payment/model")
type RechargeOrderRepo struct{ db *gorm.DB }
func NewRechargeOrderRepo(db *gorm.DB) *RechargeOrderRepo { return &RechargeOrderRepo{db: db} }
func (r *RechargeOrderRepo) Create(o *model.RechargeOrder) error { return r.db.Create(o).Error }
func (r *RechargeOrderRepo) UpdateStatus(id uint, status string) error { return r.db.Model(&model.RechargeOrder{}).Where("id = ?", id).Update("status", status).Error }
