package handler
import ("github.com/gin-gonic/gin"; "github.com/wtb-ordering/pkg/response")
type AnalyticsHandler struct{}
func NewAnalyticsHandler() *AnalyticsHandler { return &AnalyticsHandler{} }
func (h *AnalyticsHandler) Dashboard(c *gin.Context) { response.Success(c, gin.H{"today_revenue": 0, "today_orders": 0, "active_users": 0}) }
func (h *AnalyticsHandler) Revenue(c *gin.Context) { response.Success(c, gin.H{"total": 0, "list": []gin.H{}}) }
func (h *AnalyticsHandler) Dishes(c *gin.Context) { response.Success(c, gin.H{"total": 0, "list": []gin.H{}}) }
func (h *AnalyticsHandler) Members(c *gin.Context) { response.Success(c, gin.H{"total": 0, "list": []gin.H{}}) }
func (h *AnalyticsHandler) Points(c *gin.Context) { response.Success(c, gin.H{"total": 0, "list": []gin.H{}}) }
func (h *AnalyticsHandler) Coupons(c *gin.Context) { response.Success(c, gin.H{"total": 0, "list": []gin.H{}}) }
func (h *AnalyticsHandler) Activities(c *gin.Context) { response.Success(c, gin.H{"total": 0, "list": []gin.H{}}) }
func (h *AnalyticsHandler) Export(c *gin.Context) { response.Success(c, gin.H{"url": ""}) }
