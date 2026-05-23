package migration
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/points/model")
func AutoMigrate(db *gorm.DB) error { return db.AutoMigrate(&model.PointsRule{}, &model.UserPoints{}, &model.PointsLog{}, &model.ExchangeGoods{}, &model.ExchangeOrder{}) }
