package jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	OpenID string `json:"openid"`
	Level  int    `json:"level"` // 0=普通 1=会员 2=充值
	jwt.RegisteredClaims
}

var jwtSecret []byte

func Init(secret string) { jwtSecret = []byte(secret) }

func GenerateToken(userID, openID string, level int) (string, error) {
	claims := Claims{
		UserID: userID, OpenID: openID, Level: level,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{},
		func(t *jwt.Token) (interface{}, error) { return jwtSecret, nil })
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// AuthMiddleware gin JWT 鉴权中间件
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"code": 40101, "message": "未登录"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		Init(secret)
		claims, err := ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"code": 40102, "message": "Token 过期或无效"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("openid", claims.OpenID)
		c.Set("level", claims.Level)
		c.Next()
	}
}
