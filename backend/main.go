package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/wtb-ordering/internal/wechat"

	activityhandler "github.com/wtb-ordering/services/activity/handler"
	activitymigration "github.com/wtb-ordering/services/activity/migration"
	activityrepo "github.com/wtb-ordering/services/activity/repository"
	activityservice "github.com/wtb-ordering/services/activity/service"

	adminhandler "github.com/wtb-ordering/services/admin/handler"

	analyticshandler "github.com/wtb-ordering/services/analytics/handler"

	menuhandler "github.com/wtb-ordering/services/menu/handler"
	menumigration "github.com/wtb-ordering/services/menu/migration"
	menurepo "github.com/wtb-ordering/services/menu/repository"
	menuservice "github.com/wtb-ordering/services/menu/service"

	orderhandler "github.com/wtb-ordering/services/order/handler"
	ordermigration "github.com/wtb-ordering/services/order/migration"
	orderrepo "github.com/wtb-ordering/services/order/repository"
	orderservice "github.com/wtb-ordering/services/order/service"

	paymenthandler "github.com/wtb-ordering/services/payment/handler"
	paymentmigration "github.com/wtb-ordering/services/payment/migration"
	paymentrepo "github.com/wtb-ordering/services/payment/repository"
	paymentservice "github.com/wtb-ordering/services/payment/service"

	pointshandler "github.com/wtb-ordering/services/points/handler"
	pointsmigration "github.com/wtb-ordering/services/points/migration"
	pointsrepo "github.com/wtb-ordering/services/points/repository"
	pointsservice "github.com/wtb-ordering/services/points/service"

	pricinghandler "github.com/wtb-ordering/services/pricing/handler"
	pricingmigration "github.com/wtb-ordering/services/pricing/migration"
	pricingrepo "github.com/wtb-ordering/services/pricing/repository"
	pricingservice "github.com/wtb-ordering/services/pricing/service"

	seathandler "github.com/wtb-ordering/services/seat/handler"
	seatmigration "github.com/wtb-ordering/services/seat/migration"
	seatrepo "github.com/wtb-ordering/services/seat/repository"
	seatservice "github.com/wtb-ordering/services/seat/service"

	userhandler "github.com/wtb-ordering/services/user/handler"
	usermigration "github.com/wtb-ordering/services/user/migration"
	userrepo "github.com/wtb-ordering/services/user/repository"
	userservice "github.com/wtb-ordering/services/user/service"
)

func openDB(dsn, dbname string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn+" dbname="+dbname), &gorm.Config{})
	if err != nil {
		log.Fatalf("%s db error: %v", dbname, err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("%s sql.DB error: %v", dbname, err)
	}
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	return db
}

func main() {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}
	gin.SetMode(ginMode)

	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		dbDSN = "host=/tmp user=admin sslmode=disable TimeZone=Asia/Shanghai"
	}
	log.Printf("database DSN: %s", dbDSN)

	userDB := openDB(dbDSN, "wtb_user")
	menuDB := openDB(dbDSN, "wtb_menu")
	orderDB := openDB(dbDSN, "wtb_order")
	pricingDB := openDB(dbDSN, "wtb_pricing")
	activityDB := openDB(dbDSN, "wtb_activity")
	pointsDB := openDB(dbDSN, "wtb_points")
	paymentDB := openDB(dbDSN, "wtb_payment")
	seatDB := openDB(dbDSN, "wtb_seat")
	log.Println("all databases connected")

	if os.Getenv("AUTO_MIGRATE") != "false" {
		log.Println("running database migrations...")
		if err := usermigration.AutoMigrate(userDB); err != nil {
			log.Fatalf("user migration error: %v", err)
		}
		if err := menumigration.AutoMigrate(menuDB); err != nil {
			log.Fatalf("menu migration error: %v", err)
		}
		if err := ordermigration.AutoMigrate(orderDB); err != nil {
			log.Fatalf("order migration error: %v", err)
		}
		if err := pricingmigration.AutoMigrate(pricingDB); err != nil {
			log.Fatalf("pricing migration error: %v", err)
		}
		if err := activitymigration.AutoMigrate(activityDB); err != nil {
			log.Fatalf("activity migration error: %v", err)
		}
		if err := pointsmigration.AutoMigrate(pointsDB); err != nil {
			log.Fatalf("points migration error: %v", err)
		}
		if err := paymentmigration.AutoMigrate(paymentDB); err != nil {
			log.Fatalf("payment migration error: %v", err)
		}
		if err := seatmigration.AutoMigrate(seatDB); err != nil {
			log.Fatalf("seat migration error: %v", err)
		}
		log.Println("all migrations completed")
	}

	// User Service
	ur := userrepo.NewUserRepo(userDB)
	urr := userrepo.NewRechargeRepo(userDB)
	ubr := userrepo.NewBalanceLogRepo(userDB)
	ucr := userrepo.NewConsumptionRepo(userDB)
	upr := userrepo.NewPetRepo(userDB)
	wc := wechat.NewClient(wechat.Config{
		AppID:     os.Getenv("WX_APPID"),
		AppSecret: os.Getenv("WX_APPSECRET"),
	})
	userSvc := userservice.NewUserService(ur, urr, ubr, ucr, upr, wc, "")

	// Menu Service
	mcr := menurepo.NewCategoryRepo(menuDB)
	mdr := menurepo.NewDishRepo(menuDB)
	mpr := menurepo.NewDishPriceRepo(menuDB)
	msr := menurepo.NewDishStockRepo(menuDB)
	menuSvc := menuservice.NewMenuService(mcr, mdr, mpr, msr)
	menuHandler := menuhandler.NewMenuHandler(menuSvc)

	// Order Service
	or := orderrepo.NewOrderRepo(orderDB)
	oir := orderrepo.NewOrderItemRepo(orderDB)
	osl := orderrepo.NewOrderStatusLogRepo(orderDB)
	cr := orderrepo.NewCartRepo(orderDB)
	cartSvc := orderservice.NewCartService(cr)
	orderSvc := orderservice.NewOrderService(or, oir, osl, cartSvc)
	orderHandler := orderhandler.NewOrderHandler(orderSvc, cartSvc)
	userHandler := userhandler.NewUserHandler(userSvc, or)

	// Pricing Service
	prr := pricingrepo.NewPriceRuleRepo(pricingDB)
	ppr := pricingrepo.NewPromotionRepo(pricingDB)
	pcr := pricingrepo.NewComboRepo(pricingDB)
	prpr := pricingrepo.NewRechargePlanRepo(pricingDB)
	pricingSvc := pricingservice.NewPricingService(prr, ppr, pcr, prpr)
	pricingHandler := pricinghandler.NewPricingHandler(pricingSvc)

	// Activity Service
	ar := activityrepo.NewAnnouncementRepo(activityDB)
	acr := activityrepo.NewActivityRepo(activityDB)
	areg := activityrepo.NewRegistrationRepo(activityDB)
	activitySvc := activityservice.NewActivityService(ar, acr, areg)
	activityHandler := activityhandler.NewActivityHandler(activitySvc)

	// Points Service
	pur := pointsrepo.NewUserPointsRepo(pointsDB)
	plr := pointsrepo.NewPointsLogRepo(pointsDB)
	pgr := pointsrepo.NewExchangeGoodsRepo(pointsDB)
	por := pointsrepo.NewExchangeOrderRepo(pointsDB)
	pointsSvc := pointsservice.NewPointsService(pur, plr, pgr, por)
	pointsHandler := pointshandler.NewPointsHandler(pointsSvc)

	// Payment Service
	pmor := paymentrepo.NewPaymentOrderRepo(paymentDB)
	pmrr := paymentrepo.NewPaymentRecordRepo(paymentDB)
	pmref := paymentrepo.NewRefundRecordRepo(paymentDB)
	pmrech := paymentrepo.NewRechargeOrderRepo(paymentDB)
	paymentSvc := paymentservice.NewPaymentService(pmor, pmrr, pmref, pmrech)
	paymentHandler := paymenthandler.NewPaymentHandler(paymentSvc)

	// Seat Service
	sear := seatrepo.NewAreaRepo(seatDB)
	ser := seatrepo.NewSeatRepo(seatDB)
	ssl := seatrepo.NewSeatStatusLogRepo(seatDB)
	seatSvc := seatservice.NewSeatService(sear, ser, ssl, wc)
	seatHandler := seathandler.NewSeatHandler(seatSvc)

	// Admin & Analytics
	adminHandler := adminhandler.NewAdminHandler()
	analyticsHandler := analyticshandler.NewAnalyticsHandler()

	handlers := &Handlers{
		User:      userHandler,
		Menu:      menuHandler,
		Order:     orderHandler,
		Pricing:   pricingHandler,
		Activity:  activityHandler,
		Points:    pointsHandler,
		Payment:   paymentHandler,
		Seat:      seatHandler,
		Admin:     adminHandler,
		Analytics: analyticsHandler,
	}

	r := setupRouter(handlers, ur, userDB, orderDB)

	addr := ":8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("backend starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	var sqlDBs []*sql.DB
	dbList := []*gorm.DB{userDB, menuDB, orderDB, pricingDB, activityDB, pointsDB, paymentDB, seatDB}
	for _, db := range dbList {
		if sdb, err := db.DB(); err == nil {
			sqlDBs = append(sqlDBs, sdb)
		}
	}
	for _, sdb := range sqlDBs {
		sdb.Close()
	}
	log.Println("server exited")
}
