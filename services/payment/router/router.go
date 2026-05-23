package router
import ("github.com/gin-gonic/gin"; "github.com/wtb-ordering/pkg/jwt"; "github.com/wtb-ordering/services/payment/handler")
func SetupRouter(h *handler.PaymentHandler, jwtSecret []byte) *gin.Engine {
	r := gin.Default(); r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	api := r.Group("/api/pay"); api.POST("/callback/wx", h.WxCallback)
	auth := api.Group(""); auth.Use(jwt.AuthMiddleware(string(jwtSecret)))
	{ auth.POST("/create", h.Create); auth.POST("/wx/prepay", h.WxPrepay); auth.POST("/balance", h.BalancePay); auth.POST("/recharge", h.Recharge); auth.POST("/refund", h.Refund); auth.GET("/query/:outTradeNo", h.Query) }
	return r }
