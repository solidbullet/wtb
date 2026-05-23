package router
import ("net/http/httptest"; "testing"; "github.com/gin-gonic/gin")
func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode); r := SetupRouter([]byte("test")); w := httptest.NewRecorder(); req := httptest.NewRequest("GET", "/health", nil); r.ServeHTTP(w, req); if w.Code != 200 { t.Fatalf("expected 200, got %d", w.Code) } }
