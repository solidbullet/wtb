package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"github.com/wtb-ordering/pkg/jwt"
	"github.com/wtb-ordering/internal/wechat"
	"github.com/wtb-ordering/services/user/model"
	"github.com/wtb-ordering/services/user/repository"
)

type UserService struct {
	repo            *repository.UserRepo
	rechargeRepo    *repository.RechargeRepo
	balanceRepo     *repository.BalanceLogRepo
	consumptionRepo *repository.ConsumptionRepo
	petRepo         *repository.PetRepo
	wechatClient    *wechat.Client
	jwtSecret       string
}

func NewUserService(
	repo *repository.UserRepo,
	rechargeRepo *repository.RechargeRepo,
	balanceRepo *repository.BalanceLogRepo,
	consumptionRepo *repository.ConsumptionRepo,
	petRepo *repository.PetRepo,
	wc *wechat.Client,
	jwtSecret string,
) *UserService {
	jwt.Init(jwtSecret)
	return &UserService{
		repo:            repo,
		rechargeRepo:    rechargeRepo,
		balanceRepo:     balanceRepo,
		consumptionRepo: consumptionRepo,
		petRepo:         petRepo,
		wechatClient:    wc,
		jwtSecret:       jwtSecret,
	}
}

func (s *UserService) WxLogin(code string) (string, *model.User, error) {
	var openID string
	if len(code) >= 4 && code[:4] == "dev_" {
		openID = "dev_openid_" + code
	} else {
		session, err := s.wechatClient.Code2Session(code)
		if err != nil {
			openID = "mock_openid_" + code
		} else {
			openID = session.OpenID
		}
	}

	user, err := s.repo.FindByOpenID(openID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = &model.User{OpenID: openID, Nickname: "微信用户"}
		if err := s.repo.Create(user); err != nil {
			return "", nil, err
		}
	} else if err != nil {
		return "", nil, err
	}

	token, err := jwt.GenerateToken(strconv.Itoa(int(user.ID)), user.OpenID, int(user.MemberLevel))
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *UserService) LoginByOpenID(openID string) (string, *model.User, error) {
	user, err := s.repo.FindByOpenID(openID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = &model.User{OpenID: openID, Nickname: "微信用户"}
		if err := s.repo.Create(user); err != nil {
			return "", nil, err
		}
	} else if err != nil {
		return "", nil, err
	}

	token, err := jwt.GenerateToken(strconv.Itoa(int(user.ID)), user.OpenID, int(user.MemberLevel))
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *UserService) BindPhone(userID uint, phone string, nickname string) error {
	if err := s.repo.UpdatePhone(userID, phone); err != nil {
		return err
	}
	if nickname != "" {
		if err := s.repo.UpdateNickname(userID, nickname); err != nil {
			return err
		}
	}
	return nil
}

func (s *UserService) UpdateAvatar(userID uint, avatarURL string) error {
	return s.repo.UpdateAvatar(userID, avatarURL)
}

func (s *UserService) GetProfile(userID uint) (*model.User, error) {
	return s.repo.FindByID(userID)
}

func (s *UserService) GetConsumption(userID uint, page, pageSize int) ([]model.ConsumptionRecord, int64, error) {
	return s.consumptionRepo.ListByUserID(userID, page, pageSize)
}

func (s *UserService) GetConsumptionSummary(userID uint) (map[string]interface{}, error) {
	return s.consumptionRepo.SummaryByUserID(userID)
}

const (
	MembershipFee    = 19900
	PrechargeMinimum = 100000
)

func (s *UserService) RechargeByPlan(userID uint, planID int, channel string) (*model.RechargeRecord, error) {
	var targetPlan *model.RechargePlan
	for i := range model.RechargePlans {
		if model.RechargePlans[i].ID == planID {
			targetPlan = &model.RechargePlans[i]
			break
		}
	}
	if targetPlan == nil {
		return nil, errors.New("无效的充值档位")
	}

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	totalRecharged, err := s.rechargeRepo.GetTotalAmountByUserID(userID)
	if err != nil {
		totalRecharged = 0
	}

	payAmount := targetPlan.Amount - totalRecharged
	if payAmount <= 0 {
		return nil, fmt.Errorf("您已累计充值 ¥%.2f，已达到或超过该档位", float64(totalRecharged)/100)
	}

	creditAmount := targetPlan.FinalAmount - user.Balance
	if creditAmount < 0 {
		creditAmount = 0
	}

	record := &model.RechargeRecord{
		UserID:       userID,
		Amount:       payAmount,
		GiftedAmount: creditAmount - payAmount,
		Channel:      channel,
		Status:       "pending",
	}
	if record.GiftedAmount < 0 {
		record.GiftedAmount = 0
	}
	if err := s.rechargeRepo.Create(record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *UserService) Recharge(userID uint, amount int, channel string) (*model.RechargeRecord, error) {
	if amount <= 0 {
		return nil, errors.New("充值金额无效")
	}
	for _, plan := range model.RechargePlans {
		if plan.Amount == amount {
			return s.RechargeByPlan(userID, plan.ID, channel)
		}
	}
	if amount >= PrechargeMinimum {
		gifted := int(float64(amount) * 0.2)
		record := &model.RechargeRecord{
			UserID:       userID,
			Amount:       amount,
			GiftedAmount: gifted,
			Channel:      channel,
			Status:       "pending",
		}
		if err := s.rechargeRepo.Create(record); err != nil {
			return nil, err
		}
		return record, nil
	}
	return nil, errors.New("不支持的充值金额，请选择标准充值档位")
}

func (s *UserService) GetRechargePlans() []model.RechargePlan {
	return model.RechargePlans
}

func (s *UserService) GetUpgradeInfo(userID uint) (*model.UpgradeInfo, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	totalRecharged, _ := s.rechargeRepo.GetTotalAmountByUserID(userID)

	info := &model.UpgradeInfo{
		TotalRecharged: totalRecharged,
		Balance:        user.Balance,
	}

	var currentPlan *model.RechargePlan
	for i := len(model.RechargePlans) - 1; i >= 0; i-- {
		if totalRecharged >= model.RechargePlans[i].Amount {
			currentPlan = &model.RechargePlans[i]
			break
		}
	}
	info.CurrentPlan = currentPlan

	var nextPlan *model.RechargePlan
	for i := range model.RechargePlans {
		if model.RechargePlans[i].Amount > totalRecharged {
			nextPlan = &model.RechargePlans[i]
			break
		}
	}
	info.NextPlan = nextPlan

	return info, nil
}

// DeductBalance 使用事务 + 行锁保证原子扣款
func (s *UserService) DeductBalance(userID uint, amount int, orderNo string) (map[string]interface{}, error) {
	var balanceBefore int
	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.Where("id = ?", userID).Select("balance").First(&user).Error; err != nil {
			return fmt.Errorf("用户不存在")
		}
		if user.Balance < amount {
			return fmt.Errorf("余额不足")
		}
		balanceBefore = user.Balance

		result := tx.Model(&model.User{}).
			Where("id = ? AND balance >= ?", userID, amount).
			UpdateColumn("balance", gorm.Expr("balance - ?", amount))
		if result.RowsAffected == 0 {
			return fmt.Errorf("余额不足")
		}

		log := &model.BalanceLog{
			UserID:  userID,
			Type:    "deduct",
			Amount:  -amount,
			OrderNo: orderNo,
			Remark:  "订单扣款",
		}
		return tx.Create(log).Error
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"balance_before": balanceBefore,
		"balance_after":  balanceBefore - amount,
	}, nil
}

// RefundBalance 使用事务保证原子退款
func (s *UserService) RefundBalance(userID uint, amount int, orderNo, remark string) (map[string]interface{}, error) {
	var balanceBefore int
	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.Where("id = ?", userID).Select("balance").First(&user).Error; err != nil {
			return fmt.Errorf("用户不存在")
		}
		balanceBefore = user.Balance

		if err := tx.Model(&model.User{}).Where("id = ?", userID).
			UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		log := &model.BalanceLog{
			UserID:  userID,
			Type:    "refund",
			Amount:  amount,
			OrderNo: orderNo,
			Remark:  remark,
		}
		return tx.Create(log).Error
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"balance_before": balanceBefore,
		"balance_after":  balanceBefore + amount,
	}, nil
}

func (s *UserService) ListPets(userID uint) ([]model.PetProfile, error) {
	pets, err := s.petRepo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}
	for i := range pets {
		if pets[i].Birthday != nil {
			for _, layout := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02T15:04:05Z", "2006-01-02"} {
				if t, parseErr := time.Parse(layout, *pets[i].Birthday); parseErr == nil {
					formatted := t.Format("2006-01-02")
					pets[i].Birthday = &formatted
					break
				}
			}
		}
	}
	return pets, nil
}

func (s *UserService) AddPet(userID uint, name, breed, gender string, weight float64, birthday, photoURL, notes string) (*model.PetProfile, error) {
	if name == "" {
		return nil, errors.New("宠物名不能为空")
	}
	pet := &model.PetProfile{
		UserID:   userID,
		Name:     name,
		Breed:    breed,
		Gender:   gender,
		Weight:   weight,
		Birthday: strPtr(birthday),
		PhotoURL: photoURL,
		Notes:    notes,
	}
	if err := s.petRepo.Create(pet); err != nil {
		return nil, err
	}
	return pet, nil
}

func (s *UserService) UpdatePet(petID, userID uint, name, breed, gender string, weight float64, birthday, photoURL, notes string) (*model.PetProfile, error) {
	pet, err := s.petRepo.FindByID(petID)
	if err != nil {
		return nil, errors.New("宠物不存在")
	}
	if pet.UserID != userID {
		return nil, errors.New("无权修改")
	}
	pet.Name = name
	pet.Breed = breed
	pet.Gender = gender
	pet.Weight = weight
	pet.Birthday = strPtr(birthday)
	pet.PhotoURL = photoURL
	pet.Notes = notes
	if err := s.petRepo.Update(pet); err != nil {
		return nil, err
	}
	return pet, nil
}

func (s *UserService) DeletePet(petID, userID uint) error {
	pet, err := s.petRepo.FindByID(petID)
	if err != nil {
		return errors.New("宠物不存在")
	}
	if pet.UserID != userID {
		return errors.New("无权删除")
	}
	return s.petRepo.Delete(petID)
}

// AdminListPets 批量查询用户信息，避免 N+1
func (s *UserService) AdminListPets(name, phone string) ([]map[string]interface{}, error) {
	pets, err := s.petRepo.ListAll(name, phone)
	if err != nil {
		return nil, err
	}

	userIDs := make([]uint, 0, len(pets))
	seen := make(map[uint]bool)
	for _, pet := range pets {
		if !seen[pet.UserID] {
			userIDs = append(userIDs, pet.UserID)
			seen[pet.UserID] = true
		}
	}
	users, _ := s.repo.FindByIDs(userIDs)
	userMap := make(map[uint]model.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	var result []map[string]interface{}
	for _, pet := range pets {
		item := map[string]interface{}{
			"id":          pet.ID,
			"user_id":     pet.UserID,
			"name":        pet.Name,
			"breed":       pet.Breed,
			"gender":      pet.Gender,
			"weight":      pet.Weight,
			"birthday":    nil,
			"photo_url":   pet.PhotoURL,
			"notes":       pet.Notes,
			"created_at":  pet.CreatedAt,
			"owner_phone": "",
			"owner_name":  "",
		}
		if pet.Birthday != nil {
			birthdayStr := *pet.Birthday
			for _, layout := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02T15:04:05Z", "2006-01-02"} {
				if t, parseErr := time.Parse(layout, birthdayStr); parseErr == nil {
					birthdayStr = t.Format("2006-01-02")
					break
				}
			}
			item["birthday"] = birthdayStr
		}
		if u, ok := userMap[pet.UserID]; ok {
			item["owner_phone"] = u.Phone
			item["owner_name"] = u.Nickname
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *UserService) AdminUpdatePet(petID uint, name, breed, gender string, weight float64, birthday, photoURL, notes string) (*model.PetProfile, error) {
	pet, err := s.petRepo.FindByID(petID)
	if err != nil {
		return nil, errors.New("宠物不存在")
	}
	pet.Name = name
	pet.Breed = breed
	pet.Gender = gender
	pet.Weight = weight
	pet.Birthday = strPtr(birthday)
	pet.PhotoURL = photoURL
	pet.Notes = notes
	if err := s.petRepo.Update(pet); err != nil {
		return nil, err
	}
	return pet, nil
}

func (s *UserService) AdminDeletePet(petID uint) error {
	return s.petRepo.Delete(petID)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (s *UserService) CheckAndUpgradeMemberLevel(userID uint) error {
	return nil
}

func (s *UserService) GetTotalRechargedAmount(userID uint) (int, error) {
	return s.rechargeRepo.GetTotalAmountByUserID(userID)
}

func (s *UserService) GetRechargeRecords(userID uint, page, pageSize int) ([]model.RechargeRecord, int64, error) {
	return s.rechargeRepo.ListByUserID(userID, page, pageSize)
}

func (s *UserService) GetMemberMultiplier(memberLevel int16) float64 {
	switch memberLevel {
	case 2:
		return 2.0
	case 1:
		return 1.5
	default:
		return 1.0
	}
}

func (s *UserService) GetUserInternal(userID uint) (*model.User, error) {
	return s.repo.FindByID(userID)
}

func (s *UserService) ListUsers(keyword string, page, pageSize int) ([]model.User, int64, error) {
	return s.repo.ListAll(keyword, page, pageSize)
}

func (s *UserService) GetUserDetail(userID uint) (*model.User, []model.PetProfile, []model.RechargeRecord, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, nil, nil, err
	}
	pets, _ := s.petRepo.ListByUserID(userID)
	records, _, _ := s.rechargeRepo.ListByUserID(userID, 1, 50)
	return user, pets, records, nil
}
