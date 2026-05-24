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

func (r *UserRepo) DB() *gorm.DB { return r.db }

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

func (r *UserRepo) FindByIDs(ids []uint) ([]model.User, error) {
	var users []model.User
	if len(ids) == 0 {
		return users, nil
	}
	err := r.db.Where("id IN ?", ids).Find(&users).Error
	return users, err
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

func (r *UserRepo) ListAll(keyword string, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	db := r.db.Model(&model.User{})
	if keyword != "" {
		db = db.Where("nickname LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	db.Count(&total)
	err := db.Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}
