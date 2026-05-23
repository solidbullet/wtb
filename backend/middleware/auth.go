package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/user/model"
	"github.com/wtb-ordering/services/user/repository"
)

// OpenIDAuth 支持两种鉴权模式：
// 1. 小程序端：通过 X-OpenID 头鉴权（自动查询/创建用户）
// 2. 后台管理：通过 Authorization: Bearer <JWT> 头鉴权
func OpenIDAuth(userRepo *repository.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 模式1：小程序用户（X-OpenID）
		openid := c.GetHeader("X-OpenID")
		if openid != "" {
			user, err := userRepo.FindByOpenID(openid)
			if err != nil || user == nil {
				// 自动创建用户
				user = &model.User{OpenID: openid, Nickname: "微信用户"}
				if createErr := userRepo.Create(user); createErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "用户创建失败"})
					return
				}
			}
			c.Set("user_id", strconv.Itoa(int(user.ID)))
			c.Set("openid", user.OpenID)
			c.Set("level", int(user.MemberLevel))
			c.Next()
			return
		}

		// 模式2：后台管理（Bearer Token）
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwt.ParseToken(tokenStr)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("openid", claims.OpenID)
				c.Set("level", claims.Level)
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
	}
}
