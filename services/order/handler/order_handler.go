package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/pkg/response"
	"github.com/wtb-ordering/services/order/model"
	"github.com/wtb-ordering/services/order/service"
)

type OrderHandler struct {
	orderSvc *service.OrderService
	cartSvc  *service.CartService
}

func NewOrderHandler(orderSvc *service.OrderService, cartSvc *service.CartService) *OrderHandler {
	return &OrderHandler{orderSvc: orderSvc, cartSvc: cartSvc}
}

func (h *OrderHandler) CartAdd(c *gin.Context) {
	var req struct {
		SeatID    string `json:"seat_id"`
		DishID    uint   `json:"dish_id"`
		DishName  string `json:"dish_name"`
		Quantity  int    `json:"quantity"`
		UnitPrice int    `json:"unit_price"`
		Remark    string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	item := model.CartItem{
		DishID:    req.DishID,
		DishName:  req.DishName,
		Quantity:  req.Quantity,
		UnitPrice: req.UnitPrice,
		Remark:    req.Remark,
	}
	if err := h.cartSvc.Add(req.SeatID, item); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, gin.H{"seat_id": req.SeatID})
}

func (h *OrderHandler) CartList(c *gin.Context) {
	seatID := c.Query("seat_id")
	items, err := h.cartSvc.List(seatID)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *OrderHandler) CartUpdate(c *gin.Context) {
	var req struct {
		SeatID   string `json:"seat_id"`
		DishID   uint   `json:"dish_id"`
		Quantity int    `json:"quantity"`
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	if err := h.cartSvc.Update(req.SeatID, req.DishID, req.Quantity, req.Remark); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, gin.H{"seat_id": req.SeatID})
}

func (h *OrderHandler) CartRemove(c *gin.Context) {
	var req struct {
		SeatID string `json:"seat_id"`
		DishID uint   `json:"dish_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	if err := h.cartSvc.Remove(req.SeatID, req.DishID); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, gin.H{"seat_id": req.SeatID})
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req struct {
		SeatID string `json:"seat_id"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	uid, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	order, err := h.orderSvc.CreateOrder(req.SeatID, uint(uid), req.Remark)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, order)
}

func (h *OrderHandler) GetOrderStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	order, logs, err := h.orderSvc.GetOrderStatus(uint(id))
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, gin.H{"order": order, "status_logs": logs})
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	status := c.Query("status")

	var list []model.Order
	var total int64
	var err error
	if uid == 0 {
		list, total, err = h.orderSvc.ListAllOrders(status, page, pageSize)
	} else {
		list, total, err = h.orderSvc.ListOrders(uint(uid), status, page, pageSize)
	}
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.SuccessPage(c, total, page, pageSize, list)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Operator string `json:"operator"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	if err := h.orderSvc.UpdateStatus(uint(id), req.From, req.To, req.Operator); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id})
}

func (h *OrderHandler) TodayPaidOrders(c *gin.Context) {
	list, err := h.orderSvc.ListTodayPaidOrders()
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, list)
}
