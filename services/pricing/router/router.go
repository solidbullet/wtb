package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/pricing/handler"
)

func SetupRouter(h *handler.PricingHandler, jwtSecret []byte) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/pricing")
	{
		api.GET("/recharge-plans", h.ListRechargePlans)
		api.GET("/promotions", h.ListPromotions)

		auth := api.Group("")
		auth.Use(jwt.AuthMiddleware(string(jwtSecret)))
		{
			auth.POST("/calculate", h.CalculateOrderPrice)
			auth.GET("/dish/:id", h.GetDishPrice)
			auth.POST("/admin/rule", h.CreatePriceRule)
			auth.POST("/admin/promotion", h.CreatePromotion)
			auth.POST("/admin/combo", h.CreateCombo)
			auth.POST("/admin/recharge-plan", h.CreateRechargePlan)
			auth.PUT("/admin/recharge-plan/:id", h.UpdateRechargePlan)
			auth.DELETE("/admin/recharge-plan/:id", h.DeleteRechargePlan)
		}
	}

	return r
}
