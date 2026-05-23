package repository
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/points/model")
type PointsLogRepo struct{ db *gorm.DB }
func NewPointsLogRepo(db *gorm.DB) *PointsLogRepo { return &PointsLogRepo{db: db} }
func (r *PointsLogRepo) Create(log *model.PointsLog) error { return r.db.Create(log).Error }
func (r *PointsLogRepo) ListByUserID(userID uint, page, pageSize int) ([]model.PointsLog, int64, error) {
	var logs []model.PointsLog; var total int64; r.db.Model(&model.PointsLog{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Offset((page-1)*pageSize).Limit(pageSize).Find(&logs).Error; return logs, total, err }
