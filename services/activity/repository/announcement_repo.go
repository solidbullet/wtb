package repository
import ("time"; "gorm.io/gorm"; "github.com/wtb-ordering/services/activity/model")
type AnnouncementRepo struct{ db *gorm.DB }
func NewAnnouncementRepo(db *gorm.DB) *AnnouncementRepo { return &AnnouncementRepo{db: db} }
func (r *AnnouncementRepo) Create(a *model.Announcement) error { return r.db.Create(a).Error }
func (r *AnnouncementRepo) ListActive() ([]model.Announcement, error) { var as []model.Announcement; now := time.Now(); err := r.db.Where("status = ? AND start_time <= ? AND end_time >= ?", 1, now, now).Order("sort_order asc").Find(&as).Error; return as, err }
