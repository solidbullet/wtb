package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/seat/handler"
)

func SetupRouter(h *handler.SeatHandler, jwtSecret []byte) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/seat")
	{
		auth := api.Group("")
		auth.Use(jwt.AuthMiddleware(string(jwtSecret)))
		{
			auth.GET("/areas", h.ListAreas)
			auth.POST("/areas", h.CreateArea)
			auth.GET("/list", h.ListSeats)
			auth.GET("/:id", h.GetSeat)
			auth.POST("/qrcode/batch", h.GenerateQrcodeBatch)
			auth.GET("/scan", h.ScanQrcode)
		}

		internal := api.Group("/internal")
		{
			internal.GET("/:id", h.GetSeatInternal)
		}
	}

	return r
}
