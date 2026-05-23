package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/seat/model"
)

type AreaRepo struct {
	db *gorm.DB
}

func NewAreaRepo(db *gorm.DB) *AreaRepo {
	return &AreaRepo{db: db}
}

func (r *AreaRepo) Create(area *model.Area) error {
	return r.db.Create(area).Error
}

func (r *AreaRepo) List() ([]model.Area, error) {
	var areas []model.Area
	err := r.db.Order("sort_order asc").Find(&areas).Error
	return areas, err
}

func (r *AreaRepo) FindByID(id uint) (*model.Area, error) {
	var area model.Area
	err := r.db.First(&area, id).Error
	if err != nil {
		return nil, err
	}
	return &area, nil
}
