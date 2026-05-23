package main
import ("fmt"; "log"; "github.com/wtb-ordering/pkg/jwt"; "github.com/wtb-ordering/services/analytics/config"; "github.com/wtb-ordering/services/analytics/handler"; "github.com/wtb-ordering/services/analytics/router")
func main() {
	cfg := config.Load(); jwt.Init(cfg.JWTSecret); h := handler.NewAnalyticsHandler(); r := router.SetupRouter(h, []byte(cfg.JWTSecret))
	log.Printf("analytics-service on :%d", cfg.Port); r.Run(fmt.Sprintf(":%d", cfg.Port)) }
