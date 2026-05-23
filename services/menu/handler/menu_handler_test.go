package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/menu/migration"
	"github.com/wtb-ordering/services/menu/repository"
	"github.com/wtb-ordering/services/menu/service"
)

func setupTestHandler() (*MenuHandler, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	dsn := "host=/tmp user=admin dbname=wtb_menu sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	migration.AutoMigrate(db)
	jwt.Init("test-secret")

	catRepo := repository.NewCategoryRepo(db)
	dishRepo := repository.NewDishRepo(db)
	priceRepo := repository.NewDishPriceRepo(db)
	stockRepo := repository.NewDishStockRepo(db)
	menuSvc := service.NewMenuService(catRepo, dishRepo, priceRepo, stockRepo)
	return NewMenuHandler(menuSvc), db
}

func TestCategoryTree(t *testing.T) {
	h, db := setupTestHandler()
	db.Exec("DELETE FROM categories")
	db.Exec("INSERT INTO categories (name, parent_id, sort_order, status) VALUES ('主食', 0, 1, 1)")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/menu/categories", nil)
	h.GetCategories(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 200 {
		t.Fatalf("expected 200, got %v", resp["code"])
	}

	db.Exec("DELETE FROM categories")
}

func TestSearchDishes(t *testing.T) {
	h, db := setupTestHandler()
	db.Exec("DELETE FROM dishes")
	db.Exec("INSERT INTO dishes (category_id, name, status) VALUES (1, '红烧肉饭', 1)")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/menu/search?q=红烧肉", nil)
	h.SearchDishes(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	db.Exec("DELETE FROM dishes")
}
