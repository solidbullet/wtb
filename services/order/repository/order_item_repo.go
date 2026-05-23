package repository
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/order/model")
type OrderItemRepo struct{ db *gorm.DB }
func NewOrderItemRepo(db *gorm.DB) *OrderItemRepo { return &OrderItemRepo{db: db} }
func (r *OrderItemRepo) Create(item *model.OrderItem) error { return r.db.Create(item).Error }
func (r *OrderItemRepo) ListByOrder(orderID uint) ([]model.OrderItem, error) { var items []model.OrderItem; err := r.db.Where("order_id = ?", orderID).Find(&items).Error; return items, err }
