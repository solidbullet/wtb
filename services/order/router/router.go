package router
import ("github.com/gin-gonic/gin"; "github.com/wtb-ordering/pkg/jwt"; "github.com/wtb-ordering/services/order/handler")
func SetupRouter(h *handler.OrderHandler, jwtSecret []byte) *gin.Engine {
	r := gin.Default(); r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	api := r.Group("/api/order")
	// 购物车相关接口公开访问（测试环境无需登录）
	api.POST("/cart/add", h.CartAdd)
	api.GET("/cart/list", h.CartList)
	api.PUT("/cart/update", h.CartUpdate)
	api.POST("/cart/remove", h.CartRemove)
	api.POST("/create", h.CreateOrder)
	{
		auth := api.Group("")
		auth.Use(jwt.AuthMiddleware(string(jwtSecret)))
		{ auth.GET("/:id/status", h.GetOrderStatus); auth.GET("/list", h.ListOrders); auth.PUT("/admin/status", h.UpdateStatus) }
	}
	return r }
