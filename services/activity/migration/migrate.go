package migration
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/activity/model")
func AutoMigrate(db *gorm.DB) error { return db.AutoMigrate(&model.Announcement{}, &model.Activity{}, &model.ActivityRegistration{}) }
