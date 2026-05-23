package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/menu/model"
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(cat *model.Category) error {
	return r.db.Create(cat).Error
}

func (r *CategoryRepo) Update(id uint, cat *model.Category) error {
	return r.db.Model(&model.Category{}).Where("id = ?", id).Updates(cat).Error
}

func (r *CategoryRepo) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

func (r *CategoryRepo) ListAll() ([]model.Category, error) {
	var cats []model.Category
	err := r.db.Where("status = ?", 1).Order("sort_order asc").Find(&cats).Error
	return cats, err
}
