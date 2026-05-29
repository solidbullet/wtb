package main

import (
	"net/http"
	"os"
	"strings"

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
	"gorm.io/gorm"
)

var allowedOrigins = map[string]bool{
	"http://localhost:5173": true,
	"http://localhost:3000": true,
	"https://wtb.lqqnw.cn": true,
}

func CORSMiddleware() gin.HandlerFunc {
	env := os.Getenv("ENV")
	isProd := env == "production"

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if isProd {
			if origin == "" || !allowedOrigins[origin] {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		} else if origin == "" {
			origin = "*"
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

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

func setupRouter(h *Handlers, userRepo *repository.UserRepo, userDB, orderDB *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.GET("/health", func(c *gin.Context) {
		healthy := true
		sqlDB, err := userDB.DB()
		if err != nil || sqlDB.Ping() != nil {
			healthy = false
		}
		if healthy {
			c.JSON(200, gin.H{"status": "ok"})
		} else {
			c.JSON(503, gin.H{"status": "unhealthy"})
		}
	})

	imagePath := os.Getenv("IMAGE_PATH")
	if imagePath == "" {
		imagePath = "uploads"
	}
	r.Static("/images", imagePath)

	auth := middleware.OpenIDAuth(userRepo)

	// ========== 公开接口 ==========

	r.POST("/api/user/wx-login", h.User.WxLogin)

	r.GET("/api/menu/categories", h.Menu.GetCategories)
	r.GET("/api/menu/dishes", h.Menu.ListDishes)
	r.GET("/api/menu/dish/:id", h.Menu.GetDish)
	r.GET("/api/menu/search", h.Menu.SearchDishes)
	r.POST("/api/menu/dishes/batch", h.Menu.BatchDishes)

	r.GET("/api/activity/announcements", h.Activity.ListAnnouncements)
	r.GET("/api/activity/list", h.Activity.ListActivities)

	r.GET("/api/pricing/recharge-plans", h.Pricing.ListRechargePlans)
	r.GET("/api/pricing/promotions", h.Pricing.ListPromotions)

	r.POST("/api/order/cart/add", h.Order.CartAdd)
	r.GET("/api/order/cart/list", h.Order.CartList)
	r.PUT("/api/order/cart/update", h.Order.CartUpdate)
	r.POST("/api/order/cart/remove", h.Order.CartRemove)

	r.POST("/api/pay/callback/wx", h.Payment.WxCallback)

	r.POST("/admin/login", h.Admin.AdminLogin)

	r.GET("/api/seat/scan", h.Seat.ScanQrcode)

	seatInternal := r.Group("/api/seat/internal")
	{
		seatInternal.GET("/:id", h.Seat.GetSeatInternal)
	}

	// ========== 需要鉴权的接口 ==========

	userAuth := r.Group("/api/user")
	userAuth.Use(auth)
	{
		userAuth.GET("/profile", h.User.GetProfile)
		userAuth.GET("/member-info", h.User.GetMemberInfo)
		userAuth.GET("/consumption", h.User.GetConsumption)
		userAuth.GET("/consumption/summary", h.User.GetConsumptionSummary)
		userAuth.POST("/recharge", h.User.Recharge)
		userAuth.GET("/recharge-plans", h.User.ListRechargePlans)
		userAuth.GET("/recharge-records", h.User.GetRechargeRecords)
		userAuth.GET("/upgrade-info", h.User.GetUpgradeInfo)
		userAuth.GET("/pets", h.User.ListPets)
		userAuth.POST("/pets", h.User.AddPet)
		userAuth.PUT("/pets/:id", h.User.UpdatePet)
		userAuth.DELETE("/pets/:id", h.User.DeletePet)
		userAuth.POST("/bind-phone", h.User.BindPhone)
		userAuth.POST("/avatar", h.User.UpdateAvatar)
	}

	adminUser := r.Group("/api/admin/users")
	adminUser.Use(auth)
	{
		adminUser.GET("", h.User.ListUsers)
		adminUser.GET("/:id", h.User.GetUserDetail)
	}

	userInternal := r.Group("/api/user/internal")
	{
		userInternal.POST("/balance/deduct", h.User.DeductBalance)
		userInternal.POST("/balance/refund", h.User.RefundBalance)
		userInternal.GET("/:id", h.User.GetUserInternal)
	}

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
		menuAdmin.POST("/upload", h.Menu.UploadImage)
	}

	orderAuth := r.Group("/api/order")
	orderAuth.Use(auth)
	{
		orderAuth.POST("/create", h.Order.CreateOrder)
		orderAuth.GET("/:id/status", h.Order.GetOrderStatus)
		orderAuth.GET("/list", h.Order.ListOrders)
		orderAuth.PUT("/admin/status", h.Order.UpdateStatus)
		orderAuth.GET("/admin/today-paid", h.Order.TodayPaidOrders)
	}

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

	activityAuth := r.Group("/api/activity")
	activityAuth.Use(auth)
	{
		activityAuth.POST("/:id/register", h.Activity.Register)
		activityAuth.GET("/my-registrations", h.Activity.MyRegistrations)
		activityAuth.PUT("/:id/cancel", h.Activity.CancelRegistration)
		activityAuth.POST("/admin/announcement", h.Activity.CreateAnnouncement)
		activityAuth.POST("/admin/activity", h.Activity.CreateActivity)
	}

	seatAuth := r.Group("/api/seat")
	seatAuth.Use(auth)
	{
		seatAuth.GET("/areas", h.Seat.ListAreas)
		seatAuth.POST("/areas", h.Seat.CreateArea)
		seatAuth.GET("/list", h.Seat.ListSeats)
		seatAuth.POST("/create", h.Seat.CreateSeat)
		seatAuth.GET("/:id", h.Seat.GetSeat)
		seatAuth.POST("/qrcode/batch", h.Seat.GenerateQrcodeBatch)
	}

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

	petsAdmin := r.Group("/api/pets")
	petsAdmin.Use(auth)
	{
		petsAdmin.GET("", h.User.AdminListPets)
		petsAdmin.PUT("/:id", h.User.AdminUpdatePet)
		petsAdmin.DELETE("/:id", h.User.AdminDeletePet)
	}

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/admin/") {
			path := strings.TrimPrefix(c.Request.URL.Path, "/api/admin")
			c.Request.URL.Path = "/api" + path
			c.Request.RequestURI = "/api" + path
			r.HandleContext(c)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
	})

	return r
}
