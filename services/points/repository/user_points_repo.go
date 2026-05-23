package repository
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/points/model")
type UserPointsRepo struct{ db *gorm.DB }
func NewUserPointsRepo(db *gorm.DB) *UserPointsRepo { return &UserPointsRepo{db: db} }
func (r *UserPointsRepo) FindByUserID(userID uint) (*model.UserPoints, error) {
	var up model.UserPoints; err := r.db.Where("user_id = ?", userID).First(&up).Error; if err != nil { return nil, err }; return &up, nil }
func (r *UserPointsRepo) Create(up *model.UserPoints) error { return r.db.Create(up).Error }
func (r *UserPointsRepo) UpdatePoints(userID uint, points int) error {
	return r.db.Model(&model.UserPoints{}).Where("user_id = ?", userID).UpdateColumn("total_points", gorm.Expr("total_points + ?", points)).Error }
