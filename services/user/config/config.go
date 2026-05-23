package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port      int
	DSN       string
	JWTSecret string
	Wechat    WechatConfig
}

type WechatConfig struct {
	AppID     string
	AppSecret string
	MchID     string
	APIv3Key  string
}

func Load() *Config {
	return &Config{
		Port:      getEnvInt("PORT", 8081),
		DSN:       getEnv("DATABASE_DSN", "host=/tmp user=admin dbname=wtb_user sslmode=disable TimeZone=Asia/Shanghai"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		Wechat: WechatConfig{
			AppID:     getEnv("WECHAT_APPID", ""),
			AppSecret: getEnv("WECHAT_SECRET", ""),
			MchID:     getEnv("WECHAT_MCHID", ""),
			APIv3Key:  getEnv("WECHAT_APIV3_KEY", ""),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
