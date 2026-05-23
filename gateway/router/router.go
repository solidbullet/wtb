package router
import ("net/http/httputil"; "net/url"; "strings"; "github.com/gin-gonic/gin"; "github.com/wtb-ordering/pkg/jwt")
var routes = map[string]string{
	"/api/user/": "http://localhost:8081", "/api/seat/": "http://localhost:8082", "/api/menu/": "http://localhost:8083",
	"/api/order/": "http://localhost:8084", "/api/pay/": "http://localhost:8085", "/api/points/": "http://localhost:8086",
	"/api/activity/": "http://localhost:8087", "/api/pricing/": "http://localhost:8088", "/api/analytics/": "http://localhost:8089",
	"/api/admin/": "http://localhost:8090" }
var publicPaths = []string{
	"/api/user/wx-login",
	"/api/pay/callback/wx",
	"/api/activity/announcements",
	"/api/activity/list",
	"/api/menu/categories",
	"/api/menu/dishes",
	"/api/menu/dish/",
	"/api/menu/search",
	"/api/menu/dishes/batch",
	"/api/pricing/recharge-plans",
	"/api/pricing/promotions",
	"/api/order/cart/",
	"/api/order/create",
	"/api/seat/scan",
}
func isPublic(path string) bool { for _, p := range publicPaths { if strings.HasPrefix(path, p) { return true } }; return false }
func SetupRouter(jwtSecret []byte) *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.Use(func(c *gin.Context) {
		path := c.Request.URL.Path
		if !isPublic(path) {
			auth := c.GetHeader("Authorization")
			if auth == "" { c.AbortWithStatusJSON(401, gin.H{"code": 40101, "message": "未登录"}); return }
			tokenStr := strings.TrimPrefix(auth, "Bearer "); jwt.Init(string(jwtSecret))
			if _, err := jwt.ParseToken(tokenStr); err != nil { c.AbortWithStatusJSON(401, gin.H{"code": 40102, "message": "Token 过期或无效"}); return }
		}
		for prefix, target := range routes {
			if strings.HasPrefix(path, prefix) {
				remote, _ := url.Parse(target); proxy := httputil.NewSingleHostReverseProxy(remote); proxy.ServeHTTP(c.Writer, c.Request); c.Abort(); return
			} }
		c.Next() })
	return r }
