package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/menu/model"
)

type DishRepo struct {
	db *gorm.DB
}

func NewDishRepo(db *gorm.DB) *DishRepo {
	return &DishRepo{db: db}
}

func (r *DishRepo) Create(dish *model.Dish) error {
	return r.db.Create(dish).Error
}

func (r *DishRepo) Update(id uint, dish *model.Dish) error {
	return r.db.Model(&model.Dish{}).Where("id = ?", id).Updates(dish).Error
}

func (r *DishRepo) Delete(id uint) error {
	return r.db.Delete(&model.Dish{}, id).Error
}

func (r *DishRepo) FindByID(id uint) (*model.Dish, error) {
	var dish model.Dish
	err := r.db.First(&dish, id).Error
	if err != nil {
		return nil, err
	}
	return &dish, nil
}

func (r *DishRepo) ListByCategory(categoryID uint, page, pageSize int) ([]model.Dish, int64, error) {
	var dishes []model.Dish
	var total int64
	query := r.db.Model(&model.Dish{}).Where("status = ?", 1)
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	query.Count(&total)
	err := query.Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&dishes).Error
	return dishes, total, err
}

func (r *DishRepo) Search(q string) ([]model.Dish, error) {
	var dishes []model.Dish
	err := r.db.Where("status = ? AND (name ILIKE ? OR description ILIKE ?)", 1, "%"+q+"%", "%"+q+"%").Limit(20).Find(&dishes).Error
	return dishes, err
}

func (r *DishRepo) BatchByIDs(ids []uint) ([]model.Dish, error) {
	var dishes []model.Dish
	err := r.db.Where("id IN ?", ids).Find(&dishes).Error
	return dishes, err
}

func (r *DishRepo) ListByTags(tag string, page, pageSize int) ([]model.Dish, int64, error) {
	var dishes []model.Dish
	var total int64
	query := r.db.Model(&model.Dish{}).Where("status = ? AND tags ILIKE ?", 1, "%"+tag+"%")
	query.Count(&total)
	err := query.Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&dishes).Error
	return dishes, total, err
}
