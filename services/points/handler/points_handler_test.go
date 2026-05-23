package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/points/migration"
	"github.com/wtb-ordering/services/points/repository"
	"github.com/wtb-ordering/services/points/service"
)

func setupTestHandler() (*PointsHandler, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	dsn := "host=/tmp user=admin dbname=wtb_points sslmode=disable TimeZone=Asia/Shanghai"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	migration.AutoMigrate(db)
	jwt.Init("test-secret")
	upRepo := repository.NewUserPointsRepo(db)
	logRepo := repository.NewPointsLogRepo(db)
	goodsRepo := repository.NewExchangeGoodsRepo(db)
	orderRepo := repository.NewExchangeOrderRepo(db)
	return NewPointsHandler(service.NewPointsService(upRepo, logRepo, goodsRepo, orderRepo)), db
}

func TestListGoods(t *testing.T) {
	h, db := setupTestHandler()
	db.Exec("DELETE FROM exchange_goods")
	db.Exec("INSERT INTO exchange_goods (name, points_price, stock, status) VALUES ('狗粮', 200, 50, 1)")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/points/goods", nil)
	h.ListGoods(c)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	db.Exec("DELETE FROM exchange_goods")
}
