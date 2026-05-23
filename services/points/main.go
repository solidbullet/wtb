package main
import ("fmt"; "log"; "gorm.io/driver/postgres"; "gorm.io/gorm"; "github.com/wtb-ordering/pkg/jwt"; "github.com/wtb-ordering/services/points/config"; "github.com/wtb-ordering/services/points/handler"; "github.com/wtb-ordering/services/points/migration"; "github.com/wtb-ordering/services/points/repository"; "github.com/wtb-ordering/services/points/router"; "github.com/wtb-ordering/services/points/service")
func main() {
	cfg := config.Load(); db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{}); if err != nil { log.Fatalf("db error: %v", err) }
	migration.AutoMigrate(db); jwt.Init(cfg.JWTSecret); upRepo := repository.NewUserPointsRepo(db); logRepo := repository.NewPointsLogRepo(db); goodsRepo := repository.NewExchangeGoodsRepo(db); orderRepo := repository.NewExchangeOrderRepo(db)
	svc := service.NewPointsService(upRepo, logRepo, goodsRepo, orderRepo); h := handler.NewPointsHandler(svc); r := router.SetupRouter(h, []byte(cfg.JWTSecret))
	log.Printf("points-service on :%d", cfg.Port); r.Run(fmt.Sprintf(":%d", cfg.Port)) }
