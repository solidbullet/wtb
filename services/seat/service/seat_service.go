package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wtb-ordering/internal/wechat"
	"github.com/wtb-ordering/services/seat/model"
	"github.com/wtb-ordering/services/seat/repository"
)

type SeatService struct {
	areaRepo    *repository.AreaRepo
	seatRepo    *repository.SeatRepo
	logRepo     *repository.SeatStatusLogRepo
	wxClient    *wechat.Client
}

func NewSeatService(areaRepo *repository.AreaRepo, seatRepo *repository.SeatRepo, logRepo *repository.SeatStatusLogRepo, wxClient *wechat.Client) *SeatService {
	return &SeatService{areaRepo: areaRepo, seatRepo: seatRepo, logRepo: logRepo, wxClient: wxClient}
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

func (s *SeatService) CreateSeat(areaID uint, name, seatType string, capacity int) (*model.Seat, error) {
	if name == "" {
		return nil, errors.New("座位名不能为空")
	}
	if areaID == 0 {
		return nil, errors.New("请选择区域")
	}
	if seatType == "" {
		seatType = "normal"
	}
	if capacity <= 0 {
		capacity = 4
	}
	seat := &model.Seat{
		AreaID:   areaID,
		Name:     name,
		Type:     seatType,
		Capacity: capacity,
		Status:   "available",
	}
	return seat, s.seatRepo.Create(seat)
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
	if s.wxClient == nil {
		return nil, errors.New("微信客户端未初始化")
	}

	accessToken, err := s.wxClient.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("获取微信access_token失败: %w", err)
	}

	var result []map[string]interface{}
	for _, id := range seatIDs {
		seat, err := s.seatRepo.FindByID(id)
		if err != nil {
			continue
		}

		scene := fmt.Sprintf("seat_id=%s", seat.Name)
		qrcodePath := fmt.Sprintf("pages/order/menu?seat_id=%s", seat.Name)

		envVersion := os.Getenv("WX_ENV_VERSION")
		if envVersion == "" {
			envVersion = "release"
		}
		imgData, err := s.wxClient.GetWXACodeUnlimited(accessToken, scene, "pages/order/menu", false, envVersion)
		if err != nil {
			// 单张失败不影响其他
			result = append(result, map[string]interface{}{
				"seat_id":    seat.ID,
				"seat_name":  seat.Name,
				"qrcode_url": qrcodePath,
				"wxa_path":   qrcodePath,
				"has_image":  false,
				"error":      err.Error(),
			})
			continue
		}

		imageDir := os.Getenv("IMAGE_PATH")
		if imageDir == "" {
			imageDir = "uploads"
		}
		seatsDir := filepath.Join(imageDir, "seats")
		os.MkdirAll(seatsDir, 0755)
		filename := fmt.Sprintf("seat_%d.png", seat.ID)
		dst := filepath.Join(seatsDir, filename)
		if err := os.WriteFile(dst, imgData, 0644); err != nil {
			result = append(result, map[string]interface{}{
				"seat_id":    seat.ID,
				"seat_name":  seat.Name,
				"qrcode_url": qrcodePath,
				"wxa_path":   qrcodePath,
				"has_image":  false,
				"error":      "保存图片失败: " + err.Error(),
			})
			continue
		}

		imgURL := fmt.Sprintf("/images/seats/%s", filename)
		s.seatRepo.UpdateQrcode(seat.ID, imgURL)

		result = append(result, map[string]interface{}{
			"seat_id":    seat.ID,
			"seat_name":  seat.Name,
			"qrcode_url": imgURL,
			"wxa_path":   qrcodePath,
			"has_image":  true,
		})
	}
	return result, nil
}

func (s *SeatService) ScanQrcode(code string) (*model.Seat, error) {
	seat, err := s.seatRepo.FindByQrcode(code)
	if err != nil {
		return nil, errors.New("无效的二维码")
	}
	// 允许多人同桌扫码点餐，不再因 occupied 状态拒绝
	return seat, nil
}

func (s *SeatService) GetSeatInternal(id uint) (*model.Seat, error) {
	return s.seatRepo.FindByID(id)
}
