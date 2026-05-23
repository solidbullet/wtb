package handler
import ("bytes"; "encoding/json"; "net/http/httptest"; "testing"; "github.com/gin-gonic/gin"; "gorm.io/driver/postgres"; "gorm.io/gorm"; "github.com/wtb-ordering/pkg/jwt"; "github.com/wtb-ordering/services/order/migration"; "github.com/wtb-ordering/services/order/repository"; "github.com/wtb-ordering/services/order/service")
func setupTestHandler() (*OrderHandler, *gorm.DB) {
	gin.SetMode(gin.TestMode); dsn := "host=/tmp user=admin dbname=wtb_order sslmode=disable TimeZone=Asia/Shanghai"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{}); migration.AutoMigrate(db); jwt.Init("test-secret")
	orderRepo := repository.NewOrderRepo(db); itemRepo := repository.NewOrderItemRepo(db); logRepo := repository.NewOrderStatusLogRepo(db); cartSvc := service.NewCartService("localhost:6379")
	orderSvc := service.NewOrderService(orderRepo, itemRepo, logRepo, cartSvc); return NewOrderHandler(orderSvc, cartSvc), db }
func TestCreateOrder(t *testing.T) {
	h, db := setupTestHandler(); db.Exec("DELETE FROM orders")
	body := map[string]interface{}{"seat_id": 1, "remark": ""}; jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder(); c, _ := gin.CreateTestContext(w); c.Request = httptest.NewRequest("POST", "/api/order/create", bytes.NewBuffer(jsonBody)); c.Request.Header.Set("Content-Type", "application/json"); c.Set("user_id", "1")
	h.CreateOrder(c); if w.Code != 200 { t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String()) }; db.Exec("DELETE FROM orders") }
