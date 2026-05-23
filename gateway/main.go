package main
import ("fmt"; "log"; "github.com/wtb-ordering/gateway/config"; "github.com/wtb-ordering/gateway/router")
func main() {
	cfg := config.Load(); r := router.SetupRouter([]byte(cfg.JWTSecret))
	log.Printf("gateway on :%d", cfg.Port); r.Run(fmt.Sprintf(":%d", cfg.Port)) }
