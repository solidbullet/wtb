package main
import ("fmt"; "log"; "github.com/wtb-ordering/pkg/jwt"; "github.com/wtb-ordering/services/admin/config"; "github.com/wtb-ordering/services/admin/handler"; "github.com/wtb-ordering/services/admin/router")
func main() {
	cfg := config.Load(); jwt.Init(cfg.JWTSecret); h := handler.NewAdminHandler(); r := router.SetupRouter(h, []byte(cfg.JWTSecret))
	log.Printf("admin-bff on :%d", cfg.Port); r.Run(fmt.Sprintf(":%d", cfg.Port)) }
