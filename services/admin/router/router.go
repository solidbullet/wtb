package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/admin/handler"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func SetupRouter(h *handler.AdminHandler, jwtSecret []byte) *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// Public admin login (outside /api/admin to avoid wildcard conflict)
	r.POST("/admin/login", h.AdminLogin)

	// Auth proxy - all /api/admin/* routes forwarded to downstream services
	api := r.Group("/api/admin")
	api.Use(jwt.AuthMiddleware(string(jwtSecret)))
	{ api.Any("/*path", h.Proxy) }

	return r
}
