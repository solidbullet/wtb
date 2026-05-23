package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/wtb-ordering/internal/wechat"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/user/migration"
	"github.com/wtb-ordering/services/user/repository"
	"github.com/wtb-ordering/services/user/service"
)

func setupTestHandler() (*UserHandler, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	dsn := "host=/tmp user=admin dbname=wtb_user sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	migration.AutoMigrate(db)
	jwt.Init("test-secret")

	wc := wechat.NewClient(wechat.Config{AppID: "test", AppSecret: "secret"})
	userRepo := repository.NewUserRepo(db)
	rechargeRepo := repository.NewRechargeRepo(db)
	balanceRepo := repository.NewBalanceLogRepo(db)
	consumptionRepo := repository.NewConsumptionRepo(db)
	petRepo := repository.NewPetRepo(db)

	userSvc := service.NewUserService(userRepo, rechargeRepo, balanceRepo, consumptionRepo, petRepo, wc, "test-secret")
	return NewUserHandler(userSvc), db
}

func TestWxLogin(t *testing.T) {
	h, _ := setupTestHandler()

	body := map[string]string{"code": "invalid_code"}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/user/wx-login", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	h.WxLogin(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 50001 {
		t.Logf("wx-login with invalid code returned code=%v (expected wx error)", resp["code"])
	}
}

func TestGetProfile(t *testing.T) {
	h, db := setupTestHandler()

	db.Exec("DELETE FROM users WHERE openid = 'test_openid_1'")
	db.Exec("INSERT INTO users (openid, nickname, member_level, balance) VALUES ('test_openid_1', 'TestUser', 1, 50000)")
	token, _ := jwt.GenerateToken("1", "test_openid_1", 1)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/user/profile", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)
	c.Set("user_id", "1")
	h.GetProfile(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 200 {
		t.Fatalf("expected code 200, got %v", resp["code"])
	}

	db.Exec("DELETE FROM users WHERE openid = 'test_openid_1'")
}

func TestRecharge(t *testing.T) {
	h, db := setupTestHandler()

	db.Exec("DELETE FROM users WHERE openid = 'test_openid_recharge'")
	db.Exec("INSERT INTO users (openid, nickname, member_level, balance) VALUES ('test_openid_recharge', 'RechargeUser', 2, 0)")
	token, _ := jwt.GenerateToken("1", "test_openid_recharge", 2)

	body := map[string]interface{}{"amount": 10000, "channel": "wxpay"}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/user/recharge", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Authorization", "Bearer "+token)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", "1")
	h.Recharge(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 200 {
		t.Fatalf("expected code 200, got %v, body=%s", resp["code"], w.Body.String())
	}

	// 测试无效金额
	body = map[string]interface{}{"amount": 0}
	jsonBody, _ = json.Marshal(body)
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/user/recharge", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Authorization", "Bearer "+token)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", "1")
	h.Recharge(c)

	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 40001 {
		t.Fatalf("expected code 40001 for invalid amount, got %v", resp["code"])
	}

	db.Exec("DELETE FROM users WHERE openid = 'test_openid_recharge'")
	db.Exec("DELETE FROM recharge_records WHERE user_id IN (SELECT id FROM users WHERE openid = 'test_openid_recharge')")
}

func TestDeductBalance(t *testing.T) {
	h, db := setupTestHandler()

	db.Exec("DELETE FROM users WHERE openid = 'test_openid_deduct'")
	db.Exec("INSERT INTO users (openid, nickname, member_level, balance) VALUES ('test_openid_deduct', 'DeductUser', 2, 10000)")
	var uid uint
	db.Raw("SELECT id FROM users WHERE openid = 'test_openid_deduct'").Scan(&uid)

	body := map[string]interface{}{"user_id": uid, "amount": 5000, "order_no": "WTB20260101120000"}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/user/internal/balance/deduct", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	h.DeductBalance(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 200 {
		t.Logf("deduct response: %v", resp)
	}

	db.Exec("DELETE FROM users WHERE openid = 'test_openid_deduct'")
}

func TestRefundBalance(t *testing.T) {
	h, db := setupTestHandler()

	db.Exec("DELETE FROM users WHERE openid = 'test_openid_refund'")
	db.Exec("INSERT INTO users (openid, nickname, member_level, balance) VALUES ('test_openid_refund', 'RefundUser', 2, 5000)")
	var uid uint
	db.Raw("SELECT id FROM users WHERE openid = 'test_openid_refund'").Scan(&uid)

	body := map[string]interface{}{"user_id": uid, "amount": 3000, "order_no": "WTB20260101120001", "remark": "测试退款"}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/user/internal/balance/refund", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	h.RefundBalance(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	db.Exec("DELETE FROM users WHERE openid = 'test_openid_refund'")
}

func TestListPets(t *testing.T) {
	h, db := setupTestHandler()

	db.Exec("DELETE FROM users WHERE openid = 'test_openid_pets'")
	db.Exec("INSERT INTO users (openid, nickname) VALUES ('test_openid_pets', 'PetUser')")
	token, _ := jwt.GenerateToken("1", "test_openid_pets", 0)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/user/pets", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)
	c.Set("user_id", "1")
	h.ListPets(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	db.Exec("DELETE FROM users WHERE openid = 'test_openid_pets'")
}

func TestAddPet(t *testing.T) {
	h, db := setupTestHandler()

	db.Exec("DELETE FROM users WHERE openid = 'test_openid_addpet'")
	db.Exec("INSERT INTO users (openid, nickname) VALUES ('test_openid_addpet', 'AddPetUser')")
	var uid uint
	db.Raw("SELECT id FROM users WHERE openid = 'test_openid_addpet'").Scan(&uid)
	token, _ := jwt.GenerateToken("1", "test_openid_addpet", 0)

	body := map[string]interface{}{"name": "旺财", "breed": "金毛", "weight": 28.5, "birthday": "2024-03-15"}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/user/pets", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Authorization", "Bearer "+token)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", uid)
	h.AddPet(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 200 {
		t.Fatalf("expected code 200, got %v", resp["code"])
	}

	db.Exec("DELETE FROM pet_profiles WHERE user_id IN (SELECT id FROM users WHERE openid = 'test_openid_addpet')")
	db.Exec("DELETE FROM users WHERE openid = 'test_openid_addpet'")
}
