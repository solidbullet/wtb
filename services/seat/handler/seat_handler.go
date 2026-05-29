package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/pkg/response"
	"github.com/wtb-ordering/services/seat/service"
)

type SeatHandler struct {
	svc *service.SeatService
}

func NewSeatHandler(svc *service.SeatService) *SeatHandler {
	return &SeatHandler{svc: svc}
}

// ListAreas GET /api/seat/areas
func (h *SeatHandler) ListAreas(c *gin.Context) {
	areas, err := h.svc.ListAreas()
	if err != nil {
		response.Error(c, 50001, "获取区域列表失败")
		return
	}
	response.Success(c, areas)
}

// CreateArea POST /api/seat/areas
func (h *SeatHandler) CreateArea(c *gin.Context) {
	var req struct {
		Name      string `json:"name"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	area, err := h.svc.CreateArea(req.Name, req.SortOrder)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, area)
}

// ListSeats GET /api/seat/list
func (h *SeatHandler) ListSeats(c *gin.Context) {
	areaID, _ := strconv.Atoi(c.Query("area_id"))
	seats, err := h.svc.ListSeats(uint(areaID))
	if err != nil {
		response.Error(c, 50001, "获取座位列表失败")
		return
	}
	response.Success(c, seats)
}

// CreateSeat POST /api/seat/create
func (h *SeatHandler) CreateSeat(c *gin.Context) {
	var req struct {
		AreaID   uint   `json:"area_id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		Capacity int    `json:"capacity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	seat, err := h.svc.CreateSeat(req.AreaID, req.Name, req.Type, req.Capacity)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, seat)
}

// GetSeat GET /api/seat/:id
func (h *SeatHandler) GetSeat(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	seat, logs, err := h.svc.GetSeat(uint(id))
	if err != nil {
		response.Error(c, 50001, "座位不存在")
		return
	}
	response.Success(c, gin.H{
		"id":           seat.ID,
		"area_id":      seat.AreaID,
		"name":         seat.Name,
		"type":         seat.Type,
		"capacity":     seat.Capacity,
		"qrcode_url":   seat.QrcodeURL,
		"status":       seat.Status,
		"status_logs":  logs,
		"created_at":   seat.CreatedAt,
	})
}

// GenerateQrcodeBatch POST /api/seat/qrcode/batch
func (h *SeatHandler) GenerateQrcodeBatch(c *gin.Context) {
	var req struct {
		SeatIDs []uint `json:"seat_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	result, err := h.svc.GenerateQrcodeBatch(req.SeatIDs)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, result)
}

// ScanQrcode GET /api/seat/scan
func (h *SeatHandler) ScanQrcode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.Error(c, 40001, "二维码不能为空")
		return
	}
	seat, err := h.svc.ScanQrcode(code)
	if err != nil {
		if err.Error() == "座位已被占用" {
			response.Error(c, 40004, err.Error())
			return
		}
		response.Error(c, 40001, err.Error())
		return
	}
	response.Success(c, gin.H{
		"seat_id":   seat.ID,
		"area_id":   seat.AreaID,
		"seat_name": seat.Name,
		"status":    seat.Status,
	})
}

// GetSeatInternal GET /api/seat/internal/:id
func (h *SeatHandler) GetSeatInternal(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	seat, err := h.svc.GetSeatInternal(uint(id))
	if err != nil {
		response.Error(c, 50001, "座位不存在")
		return
	}
	response.Success(c, gin.H{
		"id":      seat.ID,
		"area_id": seat.AreaID,
		"name":    seat.Name,
		"status":  seat.Status,
	})
}
