package service

import (
	"errors"
	"fmt"

	"github.com/wtb-ordering/services/seat/model"
	"github.com/wtb-ordering/services/seat/repository"
)

type SeatService struct {
	areaRepo    *repository.AreaRepo
	seatRepo    *repository.SeatRepo
	logRepo     *repository.SeatStatusLogRepo
}

func NewSeatService(areaRepo *repository.AreaRepo, seatRepo *repository.SeatRepo, logRepo *repository.SeatStatusLogRepo) *SeatService {
	return &SeatService{areaRepo: areaRepo, seatRepo: seatRepo, logRepo: logRepo}
}

func (s *SeatService) ListAreas() ([]model.Area, error) {
	return s.areaRepo.List()
}

func (s *SeatService) CreateArea(name string, sortOrder int) (*model.Area, error) {
	if name == "" {
		return nil, errors.New("区域名不能为空")
	}
	area := &model.Area{Name: name, SortOrder: sortOrder}
	return area, s.areaRepo.Create(area)
}

func (s *SeatService) ListSeats(areaID uint) ([]model.Seat, error) {
	return s.seatRepo.ListByArea(areaID)
}

func (s *SeatService) GetSeat(id uint) (*model.Seat, []model.SeatStatusLog, error) {
	seat, err := s.seatRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	logs, err := s.logRepo.ListBySeatID(id)
	return seat, logs, err
}

func (s *SeatService) GenerateQrcodeBatch(seatIDs []uint) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	for _, id := range seatIDs {
		seat, err := s.seatRepo.FindByID(id)
		if err != nil {
			continue
		}
		qrcodeURL := fmt.Sprintf("https://wtb.lqqnw.cn/seat?seat_id=%s", seat.Name)
		s.seatRepo.UpdateQrcode(seat.ID, qrcodeURL)
		result = append(result, map[string]interface{}{
			"seat_id":    seat.ID,
			"qrcode_url": qrcodeURL,
		})
	}
	return result, nil
}

func (s *SeatService) ScanQrcode(code string) (*model.Seat, error) {
	seat, err := s.seatRepo.FindByQrcode(code)
	if err != nil {
		return nil, errors.New("无效的二维码")
	}
	if seat.Status == "occupied" {
		return nil, errors.New("座位已被占用")
	}
	return seat, nil
}

func (s *SeatService) GetSeatInternal(id uint) (*model.Seat, error) {
	return s.seatRepo.FindByID(id)
}
