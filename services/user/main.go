package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/wtb-ordering/internal/wechat"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/user/config"
	"github.com/wtb-ordering/services/user/handler"
	"github.com/wtb-ordering/services/user/migration"
	"github.com/wtb-ordering/services/user/repository"
	"github.com/wtb-ordering/services/user/router"
	"github.com/wtb-ordering/services/user/service"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := migration.AutoMigrate(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	jwt.Init(cfg.JWTSecret)

	wc := wechat.NewClient(wechat.Config{
		AppID:     cfg.Wechat.AppID,
		AppSecret: cfg.Wechat.AppSecret,
		MchID:     cfg.Wechat.MchID,
		APIv3Key:  cfg.Wechat.APIv3Key,
	})

	userRepo := repository.NewUserRepo(db)
	rechargeRepo := repository.NewRechargeRepo(db)
	balanceRepo := repository.NewBalanceLogRepo(db)
	consumptionRepo := repository.NewConsumptionRepo(db)
	petRepo := repository.NewPetRepo(db)

	userSvc := service.NewUserService(userRepo, rechargeRepo, balanceRepo, consumptionRepo, petRepo, wc, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userSvc)

	r := router.SetupRouter(userHandler, []byte(cfg.JWTSecret))
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("user-service starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
