package migration
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/order/model")
func AutoMigrate(db *gorm.DB) error { return db.AutoMigrate(&model.Order{}, &model.OrderItem{}, &model.OrderStatusLog{}) }
