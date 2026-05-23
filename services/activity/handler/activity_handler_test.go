package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/activity/migration"
	"github.com/wtb-ordering/services/activity/repository"
	"github.com/wtb-ordering/services/activity/service"
)

func setupTestHandler() (*ActivityHandler, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	dsn := "host=/tmp user=admin dbname=wtb_activity sslmode=disable TimeZone=Asia/Shanghai"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	migration.AutoMigrate(db)
	jwt.Init("test-secret")
	annRepo := repository.NewAnnouncementRepo(db)
	actRepo := repository.NewActivityRepo(db)
	regRepo := repository.NewRegistrationRepo(db)
	return NewActivityHandler(service.NewActivityService(annRepo, actRepo, regRepo)), db
}

func TestListActivities(t *testing.T) {
	h, db := setupTestHandler()
	db.Exec("DELETE FROM activities")
	db.Exec("INSERT INTO activities (title, status) VALUES ('测试活动', 'published')")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/activity/list", nil)
	h.ListActivities(c)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	db.Exec("DELETE FROM activities")
}
