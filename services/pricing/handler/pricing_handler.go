package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/pkg/response"
	"github.com/wtb-ordering/services/pricing/model"
	"github.com/wtb-ordering/services/pricing/service"
)

type PricingHandler struct {
	svc *service.PricingService
}

func NewPricingHandler(svc *service.PricingService) *PricingHandler {
	return &PricingHandler{svc: svc}
}

// CalculateOrderPrice POST /api/pricing/calculate
func (h *PricingHandler) CalculateOrderPrice(c *gin.Context) {
	var req struct {
		UserLevel int `json:"user_level"`
		Items     []struct {
			DishID   uint `json:"dish_id"`
			Quantity int  `json:"quantity"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	result, err := h.svc.CalculateOrderPrice(req.UserLevel, req.Items)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, result)
}

// GetDishPrice GET /api/pricing/dish/:id
func (h *PricingHandler) GetDishPrice(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userLevel, _ := strconv.Atoi(c.DefaultQuery("user_level", "0"))
	result, err := h.svc.GetDishPrice(uint(id), userLevel)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, result)
}

// ListPromotions GET /api/pricing/promotions
func (h *PricingHandler) ListPromotions(c *gin.Context) {
	list, err := h.svc.ListPromotions()
	if err != nil {
		response.Error(c, 50001, "获取活动失败")
		return
	}
	response.Success(c, list)
}

// CreatePriceRule POST /api/pricing/admin/rule
func (h *PricingHandler) CreatePriceRule(c *gin.Context) {
	var req model.PriceRule
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	if err := h.svc.CreatePriceRule(&req); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, req)
}

// CreatePromotion POST /api/pricing/admin/promotion
func (h *PricingHandler) CreatePromotion(c *gin.Context) {
	var req model.Promotion
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	if err := h.svc.CreatePromotion(&req); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, req)
}

// CreateCombo POST /api/pricing/admin/combo
func (h *PricingHandler) CreateCombo(c *gin.Context) {
	var req model.Combo
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	if err := h.svc.CreateCombo(&req); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, req)
}

// ListRechargePlans GET /api/pricing/recharge-plans
func (h *PricingHandler) ListRechargePlans(c *gin.Context) {
	list, err := h.svc.ListRechargePlans()
	if err != nil {
		response.Error(c, 50001, "获取充值方案失败")
		return
	}
	response.Success(c, list)
}

// CreateRechargePlan POST /api/pricing/admin/recharge-plan
func (h *PricingHandler) CreateRechargePlan(c *gin.Context) {
	var req model.RechargePlan
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	if err := h.svc.CreateRechargePlan(&req); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, req)
}

// UpdateRechargePlan PUT /api/pricing/admin/recharge-plan/:id
func (h *PricingHandler) UpdateRechargePlan(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req model.RechargePlan
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	if err := h.svc.UpdateRechargePlan(uint(id), &req); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id})
}

// DeleteRechargePlan DELETE /api/pricing/admin/recharge-plan/:id
func (h *PricingHandler) DeleteRechargePlan(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteRechargePlan(uint(id)); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id})
}
