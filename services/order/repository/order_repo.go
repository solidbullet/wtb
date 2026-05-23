package repository
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/order/model")
type OrderRepo struct{ db *gorm.DB }
func NewOrderRepo(db *gorm.DB) *OrderRepo { return &OrderRepo{db: db} }
func (r *OrderRepo) Create(o *model.Order) error { return r.db.Create(o).Error }
func (r *OrderRepo) FindByID(id uint) (*model.Order, error) { var o model.Order; err := r.db.First(&o, id).Error; if err != nil { return nil, err }; return &o, nil }
func (r *OrderRepo) FindByOrderNo(no string) (*model.Order, error) { var o model.Order; err := r.db.Where("order_no = ?", no).First(&o).Error; if err != nil { return nil, err }; return &o, nil }
func (r *OrderRepo) UpdateStatus(id uint, status string) error { return r.db.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error }
func (r *OrderRepo) ListByUser(userID uint, page, pageSize int) ([]model.Order, int64, error) { var os []model.Order; var total int64; r.db.Model(&model.Order{}).Where("user_id = ?", userID).Count(&total); err := r.db.Preload("Items").Where("user_id = ?", userID).Order("created_at desc").Offset((page-1)*pageSize).Limit(pageSize).Find(&os).Error; return os, total, err }
func (r *OrderRepo) ListAll(page, pageSize int) ([]model.Order, int64, error) { var os []model.Order; var total int64; r.db.Model(&model.Order{}).Count(&total); err := r.db.Preload("Items").Order("created_at desc").Offset((page-1)*pageSize).Limit(pageSize).Find(&os).Error; return os, total, err }
