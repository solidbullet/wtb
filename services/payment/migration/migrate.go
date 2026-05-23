package migration
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/payment/model")
func AutoMigrate(db *gorm.DB) error { return db.AutoMigrate(&model.PaymentOrder{}, &model.PaymentRecord{}, &model.RefundRecord{}, &model.RechargeOrder{}) }
