package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/seat/migration"
	"github.com/wtb-ordering/services/seat/repository"
	"github.com/wtb-ordering/services/seat/service"
)

func setupTestHandler() (*SeatHandler, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	dsn := "host=/tmp user=admin dbname=wtb_seat sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	migration.AutoMigrate(db)
	jwt.Init("test-secret")

	areaRepo := repository.NewAreaRepo(db)
	seatRepo := repository.NewSeatRepo(db)
	logRepo := repository.NewSeatStatusLogRepo(db)
	seatSvc := service.NewSeatService(areaRepo, seatRepo, logRepo)
	return NewSeatHandler(seatSvc), db
}

func TestCreateAndListArea(t *testing.T) {
	h, db := setupTestHandler()
	db.Exec("DELETE FROM areas")

	body := map[string]interface{}{"name": "室内A区", "sort_order": 1}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/seat/areas", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	h.CreateArea(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 200 {
		t.Fatalf("expected code 200, got %v", resp["code"])
	}

	// list
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/seat/areas", nil)
	h.ListAreas(c)
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 200 {
		t.Fatalf("list areas failed: %v", resp)
	}

	db.Exec("DELETE FROM areas")
}

func TestScanQrcode(t *testing.T) {
	h, db := setupTestHandler()
	db.Exec("DELETE FROM seats")
	db.Exec("INSERT INTO seats (area_id, name, type, status) VALUES (1, 'A1', 'normal', 'available')")
	var seatID uint
	db.Raw("SELECT id FROM seats WHERE name = 'A1'").Scan(&seatID)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/seat/scan?code="+fmt.Sprintf("%d", seatID), nil)
	h.ScanQrcode(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	db.Exec("DELETE FROM seats")
}
