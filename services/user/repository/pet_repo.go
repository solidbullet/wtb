package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/user/model"
)

type PetRepo struct {
	db *gorm.DB
}

func NewPetRepo(db *gorm.DB) *PetRepo {
	return &PetRepo{db: db}
}

func (r *PetRepo) Create(pet *model.PetProfile) error {
	return r.db.Create(pet).Error
}

func (r *PetRepo) ListByUserID(userID uint) ([]model.PetProfile, error) {
	var pets []model.PetProfile
	err := r.db.Where("user_id = ?", userID).Find(&pets).Error
	return pets, err
}

func (r *PetRepo) FindByID(id uint) (*model.PetProfile, error) {
	var pet model.PetProfile
	err := r.db.First(&pet, id).Error
	if err != nil {
		return nil, err
	}
	return &pet, nil
}
