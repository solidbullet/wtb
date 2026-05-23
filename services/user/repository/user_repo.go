package repository

import (
	"gorm.io/gorm"
	"github.com/wtb-ordering/services/user/model"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByOpenID(openid string) (*model.User, error) {
	var user model.User
	err := r.db.Where("openid = ?", openid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) UpdateBalance(userID uint, amount int) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *UserRepo) UpdateMemberLevel(userID uint, level int16) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("member_level", level).Error
}

func (r *UserRepo) UpdateConsumption(userID uint, amount int, orderCount int) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"total_consumption": gorm.Expr("total_consumption + ?", amount),
			"total_orders":      gorm.Expr("total_orders + ?", orderCount),
		}).Error
}

func (r *UserRepo) UpdatePhone(userID uint, phone string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("phone", phone).Error
}

func (r *UserRepo) UpdateNickname(userID uint, nickname string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("nickname", nickname).Error
}

func (r *UserRepo) UpdateAvatar(userID uint, avatarURL string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("avatar_url", avatarURL).Error
}
