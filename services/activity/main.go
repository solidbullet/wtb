package main
import ("fmt"; "log"; "gorm.io/driver/postgres"; "gorm.io/gorm"; "github.com/wtb-ordering/pkg/jwt"; "github.com/wtb-ordering/services/activity/config"; "github.com/wtb-ordering/services/activity/handler"; "github.com/wtb-ordering/services/activity/migration"; "github.com/wtb-ordering/services/activity/repository"; "github.com/wtb-ordering/services/activity/router"; "github.com/wtb-ordering/services/activity/service")
func main() {
	cfg := config.Load(); db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{}); if err != nil { log.Fatalf("db error: %v", err) }
	migration.AutoMigrate(db); jwt.Init(cfg.JWTSecret); annRepo := repository.NewAnnouncementRepo(db); actRepo := repository.NewActivityRepo(db); regRepo := repository.NewRegistrationRepo(db)
	svc := service.NewActivityService(annRepo, actRepo, regRepo); h := handler.NewActivityHandler(svc); r := router.SetupRouter(h, []byte(cfg.JWTSecret))
	log.Printf("activity-service on :%d", cfg.Port); r.Run(fmt.Sprintf(":%d", cfg.Port)) }
