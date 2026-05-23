package repository
import ("gorm.io/gorm"; "github.com/wtb-ordering/services/points/model")
type ExchangeGoodsRepo struct{ db *gorm.DB }
func NewExchangeGoodsRepo(db *gorm.DB) *ExchangeGoodsRepo { return &ExchangeGoodsRepo{db: db} }
func (r *ExchangeGoodsRepo) Create(g *model.ExchangeGoods) error { return r.db.Create(g).Error }
func (r *ExchangeGoodsRepo) ListActive() ([]model.ExchangeGoods, error) { var gs []model.ExchangeGoods; err := r.db.Where("status = ?", 1).Find(&gs).Error; return gs, err }
func (r *ExchangeGoodsRepo) FindByID(id uint) (*model.ExchangeGoods, error) { var g model.ExchangeGoods; err := r.db.First(&g, id).Error; if err != nil { return nil, err }; return &g, nil }
func (r *ExchangeGoodsRepo) DeductStock(id uint, count int) error { return r.db.Model(&model.ExchangeGoods{}).Where("id = ?", id).UpdateColumn("stock", gorm.Expr("stock - ?", count)).Error }
