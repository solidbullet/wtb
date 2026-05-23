package migration

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/user/model"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.RechargeRecord{},
		&model.BalanceLog{},
		&model.ConsumptionRecord{},
		&model.PetProfile{},
	)
}
