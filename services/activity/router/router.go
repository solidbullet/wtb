package router
import ("github.com/gin-gonic/gin"; "github.com/wtb-ordering/pkg/jwt"; "github.com/wtb-ordering/services/activity/handler")
func SetupRouter(h *handler.ActivityHandler, jwtSecret []byte) *gin.Engine {
	r := gin.Default(); r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	api := r.Group("/api/activity")
	api.GET("/announcements", h.ListAnnouncements)
	api.GET("/list", h.ListActivities)
	auth := api.Group(""); auth.Use(jwt.AuthMiddleware(string(jwtSecret)))
	{ auth.POST("/:id/register", h.Register); auth.GET("/my-registrations", h.MyRegistrations); auth.PUT("/:id/cancel", h.CancelRegistration); auth.POST("/admin/announcement", h.CreateAnnouncement); auth.POST("/admin/activity", h.CreateActivity) }
	return r }
