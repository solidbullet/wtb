package repository
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/payment/model")
type PaymentOrderRepo struct{ db *gorm.DB }
func NewPaymentOrderRepo(db *gorm.DB) *PaymentOrderRepo { return &PaymentOrderRepo{db: db} }
func (r *PaymentOrderRepo) Create(o *model.PaymentOrder) error { return r.db.Create(o).Error }
func (r *PaymentOrderRepo) FindByOutTradeNo(no string) (*model.PaymentOrder, error) { var o model.PaymentOrder; err := r.db.Where("out_trade_no = ?", no).First(&o).Error; if err != nil { return nil, err }; return &o, nil }
func (r *PaymentOrderRepo) UpdateStatus(id uint, status string) error { return r.db.Model(&model.PaymentOrder{}).Where("id = ?", id).Update("status", status).Error }
