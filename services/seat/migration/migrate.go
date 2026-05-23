package migration

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/seat/model"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.Area{}, &model.Seat{}, &model.SeatStatusLog{})
}
