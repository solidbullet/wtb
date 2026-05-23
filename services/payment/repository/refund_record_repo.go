package repository
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/payment/model")
type RefundRecordRepo struct{ db *gorm.DB }
func NewRefundRecordRepo(db *gorm.DB) *RefundRecordRepo { return &RefundRecordRepo{db: db} }
func (r *RefundRecordRepo) Create(rec *model.RefundRecord) error { return r.db.Create(rec).Error }
