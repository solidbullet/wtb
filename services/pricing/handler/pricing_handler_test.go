package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/pricing/migration"
	"github.com/wtb-ordering/services/pricing/repository"
	"github.com/wtb-ordering/services/pricing/service"
)

func setupTestHandler() (*PricingHandler, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	dsn := "host=/tmp user=admin dbname=wtb_pricing sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	migration.AutoMigrate(db)
	jwt.Init("test-secret")

	ruleRepo := repository.NewPriceRuleRepo(db)
	promoRepo := repository.NewPromotionRepo(db)
	comboRepo := repository.NewComboRepo(db)
	rechargeRepo := repository.NewRechargePlanRepo(db)
	pricingSvc := service.NewPricingService(ruleRepo, promoRepo, comboRepo, rechargeRepo)
	return NewPricingHandler(pricingSvc), db
}

func TestCalculateOrderPrice(t *testing.T) {
	h, db := setupTestHandler()
	db.Exec("DELETE FROM price_rules")
	db.Exec("INSERT INTO price_rules (dish_id, rule_type, price, status) VALUES (1, 'normal', 3800, 1)")
	db.Exec("INSERT INTO price_rules (dish_id, rule_type, price, status) VALUES (1, 'member', 3200, 1)")

	body := map[string]interface{}{
		"user_level": 1,
		"items": []map[string]interface{}{
			{"dish_id": 1, "quantity": 2},
		},
	}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/pricing/calculate", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	h.CalculateOrderPrice(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 200 {
		t.Fatalf("expected 200, got %v", resp["code"])
	}

	db.Exec("DELETE FROM price_rules")
}
