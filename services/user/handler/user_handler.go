package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/pkg/response"
	"github.com/wtb-ordering/services/user/model"
	"github.com/wtb-ordering/services/user/service"
	orderModel "github.com/wtb-ordering/services/order/model"
	orderrepo "github.com/wtb-ordering/services/order/repository"
)

type UserHandler struct {
	svc *service.UserService
	orderRepo *orderrepo.OrderRepo
}

func NewUserHandler(svc *service.UserService, orderRepo *orderrepo.OrderRepo) *UserHandler {
	return &UserHandler{svc: svc, orderRepo: orderRepo}
}

// WxLogin POST /api/user/wx-login
// 支持两种模式：
// 1. 传 code：普通小程序模式，后端调用微信 jscode2session 换取 openid
// 2. 传 openid：云开发模式，前端已通过云函数获取真实 openid
func (h *UserHandler) WxLogin(c *gin.Context) {
	var req struct {
		Code   string `json:"code"`
		OpenID string `json:"openid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}

	var token string
	var user *model.User
	var err error

	if req.OpenID != "" {
		// 云开发模式：直接传入 openid，查库不存在则自动创建
		token, user, err = h.svc.LoginByOpenID(req.OpenID)
	} else if req.Code != "" {
		// 普通模式：用 code 换取 openid
		token, user, err = h.svc.WxLogin(req.Code)
	} else {
		response.Error(c, 40001, "参数错误: code 或 openid 至少传一个")
		return
	}

	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// GetProfile GET /api/user/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	id, _ := strconv.ParseUint(userID, 10, 64)

	user, err := h.svc.GetProfile(uint(id))
	if err != nil {
		response.Error(c, 50001, "获取用户信息失败")
		return
	}

	response.Success(c, user)
}

// GetConsumption GET /api/user/consumption
func (h *UserHandler) GetConsumption(c *gin.Context) {
	userID := c.GetString("user_id")
	id, _ := strconv.ParseUint(userID, 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	list, total, err := h.svc.GetConsumption(uint(id), page, pageSize)
	if err != nil {
		response.Error(c, 50001, "获取消费记录失败")
		return
	}

	response.SuccessPage(c, total, page, pageSize, list)
}

// GetConsumptionSummary GET /api/user/consumption/summary
func (h *UserHandler) GetConsumptionSummary(c *gin.Context) {
	userID := c.GetString("user_id")
	id, _ := strconv.ParseUint(userID, 10, 64)

	summary, err := h.svc.GetConsumptionSummary(uint(id))
	if err != nil {
		response.Error(c, 50001, "获取消费汇总失败")
		return
	}

	response.Success(c, summary)
}

// Recharge POST /api/user/recharge
// 支持两种模式：
// 1. 传 plan_id：按档位充值/升级（推荐）
// 2. 传 amount：兼容旧接口
func (h *UserHandler) Recharge(c *gin.Context) {
	var req struct {
		PlanID  int    `json:"plan_id"`
		Amount  int    `json:"amount"`
		Channel string `json:"channel"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}

	userID := c.GetString("user_id")
	id, _ := strconv.ParseUint(userID, 10, 64)

	var record *model.RechargeRecord
	var err error
	if req.PlanID > 0 {
		record, err = h.svc.RechargeByPlan(uint(id), req.PlanID, req.Channel)
	} else if req.Amount > 0 {
		record, err = h.svc.Recharge(uint(id), req.Amount, req.Channel)
	} else {
		response.Error(c, 40001, "请传入 plan_id 或 amount")
		return
	}

	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, gin.H{
		"recharge_order_id": record.ID,
		"amount":            record.Amount,
		"gifted_amount":     record.GiftedAmount,
		"final_amount":      record.Amount + record.GiftedAmount,
		"wx_pay_params":     gin.H{},
	})
}

// ListRechargePlans GET /api/user/recharge-plans
func (h *UserHandler) ListRechargePlans(c *gin.Context) {
	plans := h.svc.GetRechargePlans()
	response.Success(c, plans)
}

// GetRechargeRecords GET /api/user/recharge-records
func (h *UserHandler) GetRechargeRecords(c *gin.Context) {
	userID := c.GetString("user_id")
	id, _ := strconv.ParseUint(userID, 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	list, total, err := h.svc.GetRechargeRecords(uint(id), page, pageSize)
	if err != nil {
		response.Error(c, 50001, "获取充值记录失败")
		return
	}

	response.SuccessPage(c, total, page, pageSize, list)
}

// GetUpgradeInfo GET /api/user/upgrade-info
func (h *UserHandler) GetUpgradeInfo(c *gin.Context) {
	userID := c.GetString("user_id")
	id, _ := strconv.ParseUint(userID, 10, 64)

	info, err := h.svc.GetUpgradeInfo(uint(id))
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, info)
}

// DeductBalance POST /api/user/balance/deduct
func (h *UserHandler) DeductBalance(c *gin.Context) {
	var req struct {
		UserID  uint   `json:"user_id"`
		Amount  int    `json:"amount"`
		OrderNo string `json:"order_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}

	result, err := h.svc.DeductBalance(req.UserID, req.Amount, req.OrderNo)
	if err != nil {
		if err.Error() == "余额不足" {
			response.Error(c, 40003, err.Error())
			return
		}
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, result)
}

// RefundBalance POST /api/user/balance/refund
func (h *UserHandler) RefundBalance(c *gin.Context) {
	var req struct {
		UserID  uint   `json:"user_id"`
		Amount  int    `json:"amount"`
		OrderNo string `json:"order_no"`
		Remark  string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}

	result, err := h.svc.RefundBalance(req.UserID, req.Amount, req.OrderNo, req.Remark)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, result)
}

// ListPets GET /api/user/pets
func (h *UserHandler) ListPets(c *gin.Context) {
	userID := c.GetString("user_id")
	id, _ := strconv.ParseUint(userID, 10, 64)

	pets, err := h.svc.ListPets(uint(id))
	if err != nil {
		response.Error(c, 50001, "获取宠物列表失败")
		return
	}

	response.Success(c, pets)
}

// AddPet POST /api/user/pets
func (h *UserHandler) AddPet(c *gin.Context) {
	var req struct {
		Name     string  `json:"name"`
		Breed    string  `json:"breed"`
		Gender   string  `json:"gender"`
		Weight   float64 `json:"weight"`
		Birthday string  `json:"birthday"`
		PhotoURL string  `json:"photo_url"`
		Notes    string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}

	userID := c.GetString("user_id")
	id, _ := strconv.ParseUint(userID, 10, 64)

	pet, err := h.svc.AddPet(uint(id), req.Name, req.Breed, req.Gender, req.Weight, req.Birthday, req.PhotoURL, req.Notes)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, pet)
}

// UpdatePet PUT /api/user/pets/:id
func (h *UserHandler) UpdatePet(c *gin.Context) {
	petID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Name     string  `json:"name"`
		Breed    string  `json:"breed"`
		Gender   string  `json:"gender"`
		Weight   float64 `json:"weight"`
		Birthday string  `json:"birthday"`
		PhotoURL string  `json:"photo_url"`
		Notes    string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	userID := c.GetString("user_id")
	uid, _ := strconv.ParseUint(userID, 10, 64)

	pet, err := h.svc.UpdatePet(uint(petID), uint(uid), req.Name, req.Breed, req.Gender, req.Weight, req.Birthday, req.PhotoURL, req.Notes)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, pet)
}

// DeletePet DELETE /api/user/pets/:id
func (h *UserHandler) DeletePet(c *gin.Context) {
	petID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := c.GetString("user_id")
	uid, _ := strconv.ParseUint(userID, 10, 64)

	if err := h.svc.DeletePet(uint(petID), uint(uid)); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, gin.H{"id": petID})
}

// AdminListPets GET /api/admin/pets
func (h *UserHandler) AdminListPets(c *gin.Context) {
	name := c.Query("name")
	phone := c.Query("phone")
	pets, err := h.svc.AdminListPets(name, phone)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, pets)
}

// AdminUpdatePet PUT /api/admin/pets/:id
func (h *UserHandler) AdminUpdatePet(c *gin.Context) {
	petID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Name     string  `json:"name"`
		Breed    string  `json:"breed"`
		Gender   string  `json:"gender"`
		Weight   float64 `json:"weight"`
		Birthday string  `json:"birthday"`
		PhotoURL string  `json:"photo_url"`
		Notes    string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	pet, err := h.svc.AdminUpdatePet(uint(petID), req.Name, req.Breed, req.Gender, req.Weight, req.Birthday, req.PhotoURL, req.Notes)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, pet)
}

// AdminDeletePet DELETE /api/admin/pets/:id
func (h *UserHandler) AdminDeletePet(c *gin.Context) {
	petID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.AdminDeletePet(uint(petID)); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, gin.H{"id": petID})
}

// GetMemberInfo GET /api/user/member-info
func (h *UserHandler) GetMemberInfo(c *gin.Context) {
	userID := c.GetString("user_id")
	id, _ := strconv.ParseUint(userID, 10, 64)

	user, err := h.svc.GetProfile(uint(id))
	if err != nil {
		response.Error(c, 50001, "获取用户信息失败")
		return
	}

	multiplier := h.svc.GetMemberMultiplier(user.MemberLevel)

	// 获取累计充值金额
	totalRecharged, _ := h.svc.GetTotalRechargedAmount(uint(id))

	// 计算升级进度
	var nextLevelName string
	var progressPercent int
	var needRecharge int // 还需充值金额

	switch user.MemberLevel {
	case 0: // 普通客户
		nextLevelName = "会员客户"
		// 进度基于充值金额（199元门槛）
		rechargeProgress := float64(totalRecharged) / float64(service.MembershipFee) * 100
		progressPercent = int(rechargeProgress)
		if progressPercent > 100 {
			progressPercent = 100
		}
		needRecharge = service.MembershipFee - totalRecharged
		if needRecharge < 0 {
			needRecharge = 0
		}
	case 1: // 会员客户
		nextLevelName = "充值客户"
		// 进度基于预充值金额（1000元门槛）
		rechargeProgress := float64(totalRecharged) / float64(service.PrechargeMinimum) * 100
		progressPercent = int(rechargeProgress)
		if progressPercent > 100 {
			progressPercent = 100
		}
		needRecharge = service.PrechargeMinimum - totalRecharged
		if needRecharge < 0 {
			needRecharge = 0
		}
	default: // 充值客户（最高级）
		nextLevelName = ""
		progressPercent = 100
	}

	response.Success(c, gin.H{
		"member_level":      user.MemberLevel,
		"member_level_name": getMemberLevelName(user.MemberLevel),
		"multiplier":        multiplier,
		"balance":           user.Balance,
		"total_recharged":   totalRecharged,
		"next_level_name":   nextLevelName,
		"progress_percent":  progressPercent,
		"need_recharge":     needRecharge,
	})
}

func getMemberLevelName(level int16) string {
	switch level {
	case 2:
		return "充值客户"
	case 1:
		return "会员客户"
	default:
		return "普通客户"
	}
}

// BindPhone POST /api/user/bind-phone
func (h *UserHandler) BindPhone(c *gin.Context) {
	userID := c.GetString("user_id")
	id, _ := strconv.ParseUint(userID, 10, 64)

	var req struct {
		Phone    string `json:"phone" binding:"required"`
		Nickname string `json:"nickname"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: phone 不能为空")
		return
	}

	if err := h.svc.BindPhone(uint(id), req.Phone, req.Nickname); err != nil {
		response.Error(c, 50001, "绑定手机号失败")
		return
	}

	response.Success(c, gin.H{"message": "绑定成功"})
}

// UpdateAvatar POST /api/user/avatar
func (h *UserHandler) UpdateAvatar(c *gin.Context) {
	userID := c.GetString("user_id")
	id, _ := strconv.ParseUint(userID, 10, 64)

	var req struct {
		AvatarURL string `json:"avatar_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: avatar_url 不能为空")
		return
	}

	if err := h.svc.UpdateAvatar(uint(id), req.AvatarURL); err != nil {
		response.Error(c, 50001, "更新头像失败")
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

// GetUserInternal GET /api/user/internal/:id
func (h *UserHandler) GetUserInternal(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	user, err := h.svc.GetUserInternal(uint(id))
	if err != nil {
		response.Error(c, 50001, "用户不存在")
		return
	}

	response.Success(c, gin.H{
		"id":           user.ID,
		"member_level": user.MemberLevel,
		"balance":      user.Balance,
	})
}

// ListUsers GET /api/admin/users
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	list, total, err := h.svc.ListUsers(keyword, page, pageSize)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.SuccessPage(c, total, page, pageSize, list)
}

// GetUserDetail GET /api/admin/user/:id
func (h *UserHandler) GetUserDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	user, pets, records, err := h.svc.GetUserDetail(uint(id))
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	// 查询用户订单
	var orders []orderModel.Order
	var orderTotal int64
	if h.orderRepo != nil {
		orders, orderTotal, _ = h.orderRepo.ListByUser(uint(id), "", 1, 50)
	}

	response.Success(c, gin.H{
		"user":           user,
		"pets":           pets,
		"orders":         orders,
		"order_total":    orderTotal,
		"recharge_records": records,
	})
}
