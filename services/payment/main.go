package main
import ("fmt"; "log"; "gorm.io/driver/postgres"; "gorm.io/gorm"; "github.com/wtb-ordering/pkg/jwt"; "github.com/wtb-ordering/services/payment/config"; "github.com/wtb-ordering/services/payment/handler"; "github.com/wtb-ordering/services/payment/migration"; "github.com/wtb-ordering/services/payment/repository"; "github.com/wtb-ordering/services/payment/router"; "github.com/wtb-ordering/services/payment/service")
func main() {
	cfg := config.Load(); db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{}); if err != nil { log.Fatalf("db error: %v", err) }
	migration.AutoMigrate(db); jwt.Init(cfg.JWTSecret); orderRepo := repository.NewPaymentOrderRepo(db); recordRepo := repository.NewPaymentRecordRepo(db); refundRepo := repository.NewRefundRecordRepo(db); rechargeRepo := repository.NewRechargeOrderRepo(db)
	svc := service.NewPaymentService(orderRepo, recordRepo, refundRepo, rechargeRepo); h := handler.NewPaymentHandler(svc); r := router.SetupRouter(h, []byte(cfg.JWTSecret))
	log.Printf("payment-service on :%d", cfg.Port); r.Run(fmt.Sprintf(":%d", cfg.Port)) }
