package main
import ("fmt"; "log"; "gorm.io/driver/postgres"; "gorm.io/gorm"; "github.com/wtb-ordering/pkg/jwt"; "github.com/wtb-ordering/services/order/config"; "github.com/wtb-ordering/services/order/handler"; "github.com/wtb-ordering/services/order/migration"; "github.com/wtb-ordering/services/order/repository"; "github.com/wtb-ordering/services/order/router"; "github.com/wtb-ordering/services/order/service")
func main() {
	cfg := config.Load(); db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{}); if err != nil { log.Fatalf("db error: %v", err) }
	migration.AutoMigrate(db); jwt.Init(cfg.JWTSecret); orderRepo := repository.NewOrderRepo(db); itemRepo := repository.NewOrderItemRepo(db); logRepo := repository.NewOrderStatusLogRepo(db); cartSvc := service.NewCartService()
	orderSvc := service.NewOrderService(orderRepo, itemRepo, logRepo, cartSvc); h := handler.NewOrderHandler(orderSvc, cartSvc); r := router.SetupRouter(h, []byte(cfg.JWTSecret))
	log.Printf("order-service on :%d", cfg.Port); r.Run(fmt.Sprintf(":%d", cfg.Port)) }
