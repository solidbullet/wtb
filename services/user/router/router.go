package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/user/handler"
)

func SetupRouter(h *handler.UserHandler, jwtSecret []byte) *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/user")
	{
		// 公开接口
		api.POST("/wx-login", h.WxLogin)

		// 需要认证
		auth := api.Group("")
		auth.Use(jwt.AuthMiddleware(string(jwtSecret)))
		{
			auth.GET("/profile", h.GetProfile)
			auth.GET("/member-info", h.GetMemberInfo)
			auth.GET("/consumption", h.GetConsumption)
			auth.GET("/consumption/summary", h.GetConsumptionSummary)
			auth.POST("/recharge", h.Recharge)
			auth.GET("/pets", h.ListPets)
			auth.POST("/pets", h.AddPet)
		}

		// 内部接口（不走网关，服务间直连）
		internal := api.Group("/internal")
		{
			internal.POST("/balance/deduct", h.DeductBalance)
			internal.POST("/balance/refund", h.RefundBalance)
			internal.GET("/:id", h.GetUserInternal)
		}
	}

	return r
}
