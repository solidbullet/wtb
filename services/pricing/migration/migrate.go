package migration

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/pricing/model"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.PriceRule{}, &model.Promotion{}, &model.Combo{}, &model.RechargePlan{})
}
