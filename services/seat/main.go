package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/services/seat/config"
	"github.com/wtb-ordering/services/seat/handler"
	"github.com/wtb-ordering/services/seat/migration"
	"github.com/wtb-ordering/services/seat/repository"
	"github.com/wtb-ordering/services/seat/router"
	"github.com/wtb-ordering/services/seat/service"
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

	areaRepo := repository.NewAreaRepo(db)
	seatRepo := repository.NewSeatRepo(db)
	logRepo := repository.NewSeatStatusLogRepo(db)

	seatSvc := service.NewSeatService(areaRepo, seatRepo, logRepo)
	seatHandler := handler.NewSeatHandler(seatSvc)

	r := router.SetupRouter(seatHandler, []byte(cfg.JWTSecret))
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("seat-service starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
