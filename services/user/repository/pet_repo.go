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

func (r *PetRepo) Update(pet *model.PetProfile) error {
	return r.db.Save(pet).Error
}

func (r *PetRepo) Delete(id uint) error {
	return r.db.Delete(&model.PetProfile{}, id).Error
}

// ListAll 查询所有宠物，支持按名字和主人手机号筛选
func (r *PetRepo) ListAll(name, phone string) ([]model.PetProfile, error) {
	var pets []model.PetProfile
	query := r.db.Model(&model.PetProfile{})
	if name != "" {
		query = query.Where("pet_profiles.name LIKE ?", "%"+name+"%")
	}
	if phone != "" {
		query = query.Joins("JOIN users ON users.id = pet_profiles.user_id").
			Where("users.phone LIKE ?", "%"+phone+"%")
	}
	err := query.Find(&pets).Error
	return pets, err
}
