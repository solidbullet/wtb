package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/services/user/model"
	"github.com/wtb-ordering/services/user/repository"
)

// OpenIDAuth 从请求头读取 X-OpenID，查询/创建用户，将 user_id 存入 context
func OpenIDAuth(userRepo *repository.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		openid := c.GetHeader("X-OpenID")
		if openid == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
			return
		}

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
	}
}
