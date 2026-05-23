package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/seat/model"
)

type SeatStatusLogRepo struct {
	db *gorm.DB
}

func NewSeatStatusLogRepo(db *gorm.DB) *SeatStatusLogRepo {
	return &SeatStatusLogRepo{db: db}
}

func (r *SeatStatusLogRepo) Create(log *model.SeatStatusLog) error {
	return r.db.Create(log).Error
}

func (r *SeatStatusLogRepo) ListBySeatID(seatID uint) ([]model.SeatStatusLog, error) {
	var logs []model.SeatStatusLog
	err := r.db.Where("seat_id = ?", seatID).Order("changed_at desc").Find(&logs).Error
	return logs, err
}
