package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/pricing/config"
	"github.com/wtb-ordering/services/pricing/handler"
	"github.com/wtb-ordering/services/pricing/migration"
	"github.com/wtb-ordering/services/pricing/repository"
	"github.com/wtb-ordering/services/pricing/router"
	"github.com/wtb-ordering/services/pricing/service"
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

	ruleRepo := repository.NewPriceRuleRepo(db)
	promoRepo := repository.NewPromotionRepo(db)
	comboRepo := repository.NewComboRepo(db)
	rechargeRepo := repository.NewRechargePlanRepo(db)

	pricingSvc := service.NewPricingService(ruleRepo, promoRepo, comboRepo, rechargeRepo)
	pricingHandler := handler.NewPricingHandler(pricingSvc)

	r := router.SetupRouter(pricingHandler, []byte(cfg.JWTSecret))
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("pricing-service starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
