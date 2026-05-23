package repository
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/payment/model")
type PaymentRecordRepo struct{ db *gorm.DB }
func NewPaymentRecordRepo(db *gorm.DB) *PaymentRecordRepo { return &PaymentRecordRepo{db: db} }
func (r *PaymentRecordRepo) Create(rec *model.PaymentRecord) error { return r.db.Create(rec).Error }
