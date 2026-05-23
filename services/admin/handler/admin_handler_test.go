package handler
import ("net/http/httptest"; "testing"; "github.com/gin-gonic/gin")
func TestProxy(t *testing.T) {
	gin.SetMode(gin.TestMode); h := NewAdminHandler(); w := httptest.NewRecorder(); c, _ := gin.CreateTestContext(w); h.Proxy(c); if w.Code != 200 { t.Fatalf("expected 200, got %d", w.Code) } }
