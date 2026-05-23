package handler
import ("bytes"; "encoding/json"; "net/http/httptest"; "testing"; "github.com/gin-gonic/gin"; "gorm.io/driver/postgres"; "gorm.io/gorm"; "github.com/wtb-ordering/pkg/jwt"; "github.com/wtb-ordering/services/payment/migration"; "github.com/wtb-ordering/services/payment/repository"; "github.com/wtb-ordering/services/payment/service")
func setupTestHandler() (*PaymentHandler, *gorm.DB) {
	gin.SetMode(gin.TestMode); dsn := "host=/tmp user=admin dbname=wtb_payment sslmode=disable TimeZone=Asia/Shanghai"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{}); migration.AutoMigrate(db); jwt.Init("test-secret")
	orderRepo := repository.NewPaymentOrderRepo(db); recordRepo := repository.NewPaymentRecordRepo(db); refundRepo := repository.NewRefundRecordRepo(db); rechargeRepo := repository.NewRechargeOrderRepo(db)
	return NewPaymentHandler(service.NewPaymentService(orderRepo, recordRepo, refundRepo, rechargeRepo)), db }
func TestCreatePayment(t *testing.T) {
	h, db := setupTestHandler(); db.Exec("DELETE FROM payment_orders")
	body := map[string]interface{}{"order_no": "WTB001", "amount": 10000, "channel": "wxpay"}; jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder(); c, _ := gin.CreateTestContext(w); c.Request = httptest.NewRequest("POST", "/api/pay/create", bytes.NewBuffer(jsonBody)); c.Request.Header.Set("Content-Type", "application/json"); c.Set("user_id", "1")
	h.Create(c); if w.Code != 200 { t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String()) }; db.Exec("DELETE FROM payment_orders") }
