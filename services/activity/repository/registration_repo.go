package repository
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/activity/model")
type RegistrationRepo struct{ db *gorm.DB }
func NewRegistrationRepo(db *gorm.DB) *RegistrationRepo { return &RegistrationRepo{db: db} }
func (r *RegistrationRepo) Create(reg *model.ActivityRegistration) error { return r.db.Create(reg).Error }
func (r *RegistrationRepo) FindByUserAndActivity(userID, activityID uint) (*model.ActivityRegistration, error) { var reg model.ActivityRegistration; err := r.db.Where("user_id = ? AND activity_id = ?", userID, activityID).First(&reg).Error; if err != nil { return nil, err }; return &reg, nil }
func (r *RegistrationRepo) ListByUser(userID uint) ([]model.ActivityRegistration, error) { var rs []model.ActivityRegistration; err := r.db.Where("user_id = ?", userID).Find(&rs).Error; return rs, err }
func (r *RegistrationRepo) Cancel(id uint) error { return r.db.Model(&model.ActivityRegistration{}).Where("id = ?", id).Update("status", "cancelled").Error }
