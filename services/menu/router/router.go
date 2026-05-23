package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/menu/handler"
)

func SetupRouter(h *handler.MenuHandler, jwtSecret []byte) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/menu")
	api.GET("/categories", h.GetCategories)
	api.GET("/dishes", h.ListDishes)
	api.GET("/dish/:id", h.GetDish)
	api.GET("/search", h.SearchDishes)
	api.POST("/dishes/batch", h.BatchDishes)
	{
		auth := api.Group("")
		auth.Use(jwt.AuthMiddleware(string(jwtSecret)))
		{
			auth.POST("/admin/category", h.CreateCategory)
			auth.PUT("/admin/category/:id", h.UpdateCategory)
			auth.DELETE("/admin/category/:id", h.DeleteCategory)
			auth.POST("/admin/dish", h.CreateDish)
			auth.PUT("/admin/dish/:id", h.UpdateDish)
			auth.DELETE("/admin/dish/:id", h.DeleteDish)
			auth.POST("/admin/stock", h.SetStock)
		}
	}

	return r
}
