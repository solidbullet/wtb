package repository
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/points/model")
type ExchangeOrderRepo struct{ db *gorm.DB }
func NewExchangeOrderRepo(db *gorm.DB) *ExchangeOrderRepo { return &ExchangeOrderRepo{db: db} }
func (r *ExchangeOrderRepo) Create(o *model.ExchangeOrder) error { return r.db.Create(o).Error }
func (r *ExchangeOrderRepo) ListByUserID(userID uint) ([]model.ExchangeOrder, error) { var os []model.ExchangeOrder; err := r.db.Where("user_id = ?", userID).Find(&os).Error; return os, err }
