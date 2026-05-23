package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/seat/model"
)

type SeatRepo struct {
	db *gorm.DB
}

func NewSeatRepo(db *gorm.DB) *SeatRepo {
	return &SeatRepo{db: db}
}

func (r *SeatRepo) Create(seat *model.Seat) error {
	return r.db.Create(seat).Error
}

func (r *SeatRepo) ListByArea(areaID uint) ([]model.Seat, error) {
	var seats []model.Seat
	err := r.db.Where("area_id = ?", areaID).Find(&seats).Error
	return seats, err
}

func (r *SeatRepo) FindByID(id uint) (*model.Seat, error) {
	var seat model.Seat
	err := r.db.First(&seat, id).Error
	if err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *SeatRepo) UpdateStatus(seatID uint, status string) error {
	return r.db.Model(&model.Seat{}).Where("id = ?", seatID).Update("status", status).Error
}

func (r *SeatRepo) UpdateQrcode(seatID uint, qrcodeURL string) error {
	return r.db.Model(&model.Seat{}).Where("id = ?", seatID).Update("qrcode_url", qrcodeURL).Error
}

func (r *SeatRepo) FindByQrcode(code string) (*model.Seat, error) {
	var seat model.Seat
	err := r.db.Where("qrcode_url = ? OR id::text = ? OR name = ?", code, code, code).First(&seat).Error
	if err != nil {
		return nil, err
	}
	return &seat, nil
}
