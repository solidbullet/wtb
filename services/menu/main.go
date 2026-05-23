package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/menu/config"
	"github.com/wtb-ordering/services/menu/handler"
	"github.com/wtb-ordering/services/menu/migration"
	"github.com/wtb-ordering/services/menu/repository"
	"github.com/wtb-ordering/services/menu/router"
	"github.com/wtb-ordering/services/menu/service"
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

	catRepo := repository.NewCategoryRepo(db)
	dishRepo := repository.NewDishRepo(db)
	priceRepo := repository.NewDishPriceRepo(db)
	stockRepo := repository.NewDishStockRepo(db)

	menuSvc := service.NewMenuService(catRepo, dishRepo, priceRepo, stockRepo)
	menuHandler := handler.NewMenuHandler(menuSvc)

	r := router.SetupRouter(menuHandler, []byte(cfg.JWTSecret))
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("menu-service starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
