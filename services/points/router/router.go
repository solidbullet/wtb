package router
import ("github.com/gin-gonic/gin"; "github.com/wtb-ordering/pkg/jwt"; "github.com/wtb-ordering/services/points/handler")
func SetupRouter(h *handler.PointsHandler, jwtSecret []byte) *gin.Engine {
	r := gin.Default(); r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	api := r.Group("/api/points"); auth := api.Group(""); auth.Use(jwt.AuthMiddleware(string(jwtSecret)))
	{ auth.GET("/account", h.GetAccount); auth.GET("/logs", h.GetLogs); auth.GET("/goods", h.ListGoods); auth.POST("/exchange", h.Exchange) }
	admin := api.Group("/admin"); admin.Use(jwt.AuthMiddleware(string(jwtSecret)))
	{ admin.POST("/goods", h.CreateGoods) }
	internal := api.Group("/internal"); internal.POST("/grant", h.GrantPoints); return r }
