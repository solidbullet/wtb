package repository
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/activity/model")
type ActivityRepo struct{ db *gorm.DB }
func NewActivityRepo(db *gorm.DB) *ActivityRepo { return &ActivityRepo{db: db} }
func (r *ActivityRepo) Create(a *model.Activity) error { return r.db.Create(a).Error }
func (r *ActivityRepo) ListPublished() ([]model.Activity, error) { var as []model.Activity; err := r.db.Where("status = ?", "published").Order("created_at desc").Find(&as).Error; return as, err }
func (r *ActivityRepo) FindByID(id uint) (*model.Activity, error) { var a model.Activity; err := r.db.First(&a, id).Error; if err != nil { return nil, err }; return &a, nil }
func (r *ActivityRepo) IncrementParticipants(id uint) error { return r.db.Model(&model.Activity{}).Where("id = ?", id).UpdateColumn("current_participants", gorm.Expr("current_participants + 1")).Error }
