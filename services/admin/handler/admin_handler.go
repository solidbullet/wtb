package handler

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/pkg/response"
)

var serviceMap = map[string]string{
	"menu":     "http://localhost:8083",
	"order":    "http://localhost:8084",
	"activity": "http://localhost:8087",
	"points":   "http://localhost:8086",
	"pricing":  "http://localhost:8088",
	"user":     "http://localhost:8081",
	"seat":     "http://localhost:8082",
	"payment":  "http://localhost:8085",
}

var adminPassword = "admin123"

type AdminHandler struct {
	client *http.Client
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// AdminLogin POST /api/admin/login (public)
func (h *AdminHandler) AdminLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	if req.Username != "admin" || req.Password != adminPassword {
		response.Error(c, 40002, "用户名或密码错误")
		return
	}
	token, err := jwt.GenerateToken("admin", "admin_openid", 99)
	if err != nil {
		response.Error(c, 50001, "token生成失败")
		return
	}
	response.Success(c, gin.H{"token": token})
}

func (h *AdminHandler) Proxy(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		response.Error(c, 40001, "path required")
		return
	}

	parts := strings.SplitN(strings.TrimPrefix(path, "/"), "/", 2)
	svcKey := parts[0]
	targetBase, ok := serviceMap[svcKey]
	if !ok {
		response.Error(c, 40004, "unknown service: "+svcKey)
		return
	}

	targetURL := targetBase + "/api/" + svcKey
	if len(parts) > 1 {
		targetURL += "/" + parts[1]
	}
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		response.Error(c, 50001, "proxy error: "+err.Error())
		return
	}

	for k, v := range c.Request.Header {
		req.Header[k] = v
	}

	resp, err := h.client.Do(req)
	if err != nil {
		response.Error(c, 50002, "service unreachable: "+err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
