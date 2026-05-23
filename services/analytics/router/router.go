package router
import ("github.com/gin-gonic/gin"; "github.com/wtb-ordering/pkg/jwt"; "github.com/wtb-ordering/services/analytics/handler")
func SetupRouter(h *handler.AnalyticsHandler, jwtSecret []byte) *gin.Engine {
	r := gin.Default(); r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	api := r.Group("/api/analytics"); auth := api.Group(""); auth.Use(jwt.AuthMiddleware(string(jwtSecret)))
	{ auth.GET("/dashboard", h.Dashboard); auth.GET("/revenue", h.Revenue); auth.GET("/dishes", h.Dishes); auth.GET("/members", h.Members); auth.GET("/points", h.Points); auth.GET("/coupons", h.Coupons); auth.GET("/activities", h.Activities); auth.POST("/export", h.Export) }
	return r }
