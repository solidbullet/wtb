package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/backend/middleware"
	activityhandler "github.com/wtb-ordering/services/activity/handler"
	adminhandler "github.com/wtb-ordering/services/admin/handler"
	analyticshandler "github.com/wtb-ordering/services/analytics/handler"
	menuhandler "github.com/wtb-ordering/services/menu/handler"
	orderhandler "github.com/wtb-ordering/services/order/handler"
	paymenthandler "github.com/wtb-ordering/services/payment/handler"
	pointshandler "github.com/wtb-ordering/services/points/handler"
	pricinghandler "github.com/wtb-ordering/services/pricing/handler"
	seathandler "github.com/wtb-ordering/services/seat/handler"
	userhandler "github.com/wtb-ordering/services/user/handler"
	"github.com/wtb-ordering/services/user/repository"
)

type Handlers struct {
	User      *userhandler.UserHandler
	Menu      *menuhandler.MenuHandler
	Order     *orderhandler.OrderHandler
	Pricing   *pricinghandler.PricingHandler
	Activity  *activityhandler.ActivityHandler
	Points    *pointshandler.PointsHandler
	Payment   *paymenthandler.PaymentHandler
	Seat      *seathandler.SeatHandler
	Admin     *adminhandler.AdminHandler
	Analytics *analyticshandler.AnalyticsHandler
}

func setupRouter(h *Handlers, userRepo *repository.UserRepo) *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 静态图片服务（菜品图片、头像等）
	r.Static("/images", "../miniprogram/images")

	// 鉴权中间件
	auth := middleware.OpenIDAuth(userRepo)

	// ========== 公开接口 ==========

	// user
	r.POST("/api/user/wx-login", h.User.WxLogin)

	// menu
	r.GET("/api/menu/categories", h.Menu.GetCategories)
	r.GET("/api/menu/dishes", h.Menu.ListDishes)
	r.GET("/api/menu/dish/:id", h.Menu.GetDish)
	r.GET("/api/menu/search", h.Menu.SearchDishes)
	r.POST("/api/menu/dishes/batch", h.Menu.BatchDishes)

	// activity
	r.GET("/api/activity/announcements", h.Activity.ListAnnouncements)
	r.GET("/api/activity/list", h.Activity.ListActivities)

	// pricing
	r.GET("/api/pricing/recharge-plans", h.Pricing.ListRechargePlans)
	r.GET("/api/pricing/promotions", h.Pricing.ListPromotions)

	// order（购物车公开，方便测试）
	r.POST("/api/order/cart/add", h.Order.CartAdd)
	r.GET("/api/order/cart/list", h.Order.CartList)
	r.PUT("/api/order/cart/update", h.Order.CartUpdate)
	r.POST("/api/order/cart/remove", h.Order.CartRemove)

	// payment callback
	r.POST("/api/pay/callback/wx", h.Payment.WxCallback)

	// admin login
	r.POST("/admin/login", h.Admin.AdminLogin)

	// seat scan (public)
	r.GET("/api/seat/scan", h.Seat.ScanQrcode)

	// seat internal
	seatInternal := r.Group("/api/seat/internal")
	{
		seatInternal.GET("/:id", h.Seat.GetSeatInternal)
	}

	// ========== 需要鉴权的接口 ==========

	// user
	userAuth := r.Group("/api/user")
	userAuth.Use(auth)
	{
		userAuth.GET("/profile", h.User.GetProfile)
		userAuth.GET("/member-info", h.User.GetMemberInfo)
		userAuth.GET("/consumption", h.User.GetConsumption)
		userAuth.GET("/consumption/summary", h.User.GetConsumptionSummary)
		userAuth.POST("/recharge", h.User.Recharge)
		userAuth.GET("/recharge-plans", h.User.ListRechargePlans)
		userAuth.GET("/upgrade-info", h.User.GetUpgradeInfo)
		userAuth.GET("/pets", h.User.ListPets)
		userAuth.POST("/pets", h.User.AddPet)
		userAuth.POST("/bind-phone", h.User.BindPhone)
		userAuth.POST("/avatar", h.User.UpdateAvatar)
	}

	// user internal
	userInternal := r.Group("/api/user/internal")
	{
		userInternal.POST("/balance/deduct", h.User.DeductBalance)
		userInternal.POST("/balance/refund", h.User.RefundBalance)
		userInternal.GET("/:id", h.User.GetUserInternal)
	}

	// menu admin
	menuAdmin := r.Group("/api/menu/admin")
	menuAdmin.Use(auth)
	{
		menuAdmin.POST("/category", h.Menu.CreateCategory)
		menuAdmin.PUT("/category/:id", h.Menu.UpdateCategory)
		menuAdmin.DELETE("/category/:id", h.Menu.DeleteCategory)
		menuAdmin.POST("/dish", h.Menu.CreateDish)
		menuAdmin.PUT("/dish/:id", h.Menu.UpdateDish)
		menuAdmin.DELETE("/dish/:id", h.Menu.DeleteDish)
		menuAdmin.POST("/stock", h.Menu.SetStock)
	}

	// order auth
	orderAuth := r.Group("/api/order")
	orderAuth.Use(auth)
	{
		orderAuth.POST("/create", h.Order.CreateOrder)
		orderAuth.GET("/:id/status", h.Order.GetOrderStatus)
		orderAuth.GET("/list", h.Order.ListOrders)
		orderAuth.PUT("/admin/status", h.Order.UpdateStatus)
	}

	// payment auth
	payAuth := r.Group("/api/pay")
	payAuth.Use(auth)
	{
		payAuth.POST("/create", h.Payment.Create)
		payAuth.POST("/wx/prepay", h.Payment.WxPrepay)
		payAuth.POST("/balance", h.Payment.BalancePay)
		payAuth.POST("/recharge", h.Payment.Recharge)
		payAuth.POST("/refund", h.Payment.Refund)
		payAuth.GET("/query/:outTradeNo", h.Payment.Query)
	}

	// points
	pointsAuth := r.Group("/api/points")
	pointsAuth.Use(auth)
	{
		pointsAuth.GET("/account", h.Points.GetAccount)
		pointsAuth.GET("/logs", h.Points.GetLogs)
		pointsAuth.GET("/goods", h.Points.ListGoods)
		pointsAuth.POST("/exchange", h.Points.Exchange)
	}
	pointsAdmin := r.Group("/api/points/admin")
	pointsAdmin.Use(auth)
	{
		pointsAdmin.POST("/goods", h.Points.CreateGoods)
	}
	pointsInternal := r.Group("/api/points/internal")
	{
		pointsInternal.POST("/grant", h.Points.GrantPoints)
	}

	// pricing auth
	pricingAuth := r.Group("/api/pricing")
	pricingAuth.Use(auth)
	{
		pricingAuth.POST("/calculate", h.Pricing.CalculateOrderPrice)
		pricingAuth.GET("/dish/:id", h.Pricing.GetDishPrice)
	}
	pricingAdmin := r.Group("/api/pricing/admin")
	pricingAdmin.Use(auth)
	{
		pricingAdmin.POST("/rule", h.Pricing.CreatePriceRule)
		pricingAdmin.POST("/promotion", h.Pricing.CreatePromotion)
		pricingAdmin.POST("/combo", h.Pricing.CreateCombo)
		pricingAdmin.POST("/recharge-plan", h.Pricing.CreateRechargePlan)
		pricingAdmin.PUT("/recharge-plan/:id", h.Pricing.UpdateRechargePlan)
		pricingAdmin.DELETE("/recharge-plan/:id", h.Pricing.DeleteRechargePlan)
	}

	// activity auth
	activityAuth := r.Group("/api/activity")
	activityAuth.Use(auth)
	{
		activityAuth.POST("/:id/register", h.Activity.Register)
		activityAuth.GET("/my-registrations", h.Activity.MyRegistrations)
		activityAuth.PUT("/:id/cancel", h.Activity.CancelRegistration)
		activityAuth.POST("/admin/announcement", h.Activity.CreateAnnouncement)
		activityAuth.POST("/admin/activity", h.Activity.CreateActivity)
	}

	// seat
	seatAuth := r.Group("/api/seat")
	seatAuth.Use(auth)
	{
		seatAuth.GET("/areas", h.Seat.ListAreas)
		seatAuth.POST("/areas", h.Seat.CreateArea)
		seatAuth.GET("/list", h.Seat.ListSeats)
		seatAuth.GET("/:id", h.Seat.GetSeat)
		seatAuth.POST("/qrcode/batch", h.Seat.GenerateQrcodeBatch)
	}

	// analytics
	analyticsAuth := r.Group("/api/analytics")
	analyticsAuth.Use(auth)
	{
		analyticsAuth.GET("/dashboard", h.Analytics.Dashboard)
		analyticsAuth.GET("/revenue", h.Analytics.Revenue)
		analyticsAuth.GET("/dishes", h.Analytics.Dishes)
		analyticsAuth.GET("/members", h.Analytics.Members)
		analyticsAuth.GET("/points", h.Analytics.Points)
		analyticsAuth.GET("/coupons", h.Analytics.Coupons)
		analyticsAuth.GET("/activities", h.Analytics.Activities)
		analyticsAuth.POST("/export", h.Analytics.Export)
	}

	// admin
	adminAPI := r.Group("/api/admin")
	adminAPI.Use(auth)
	{
		adminAPI.Any("/*path", h.Admin.Proxy)
	}

	return r
}
