package handler

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/pkg/response"
)

func getAdminPassword() string {
	pw := os.Getenv("ADMIN_PASSWORD")
	if pw != "" {
		return pw
	}
	return "1234"
}

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) AdminLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	if req.Username != "admin" || req.Password != getAdminPassword() {
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
