package config
import ("os"; "strconv")
type Config struct { Port int; DSN string; JWTSecret string }
func Load() *Config { return &Config{Port: getEnvInt("PORT", 8084), DSN: getEnv("DATABASE_DSN", "host=/tmp user=admin dbname=wtb_order sslmode=disable TimeZone=Asia/Shanghai"), JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-in-production")} }
func getEnv(key, defaultVal string) string { if val := os.Getenv(key); val != "" { return val }; return defaultVal }
func getEnvInt(key string, defaultVal int) int { if val := os.Getenv(key); val != "" { if i, err := strconv.Atoi(val); err == nil { return i } }; return defaultVal }
