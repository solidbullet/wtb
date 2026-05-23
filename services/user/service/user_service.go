package service

import (
	"errors"
	"fmt"
	"strconv"

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
			// 微信配置未设置或接口调用失败，fallback：用 code 生成模拟 openid
			// 上线前请确保配置了正确的 AppID 和 AppSecret
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

// LoginByOpenID 云开发模式：直接传入 openid，查库不存在则自动创建用户
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

// 会员门槛（兼容旧常量）
const (
	MembershipFee    = 19900  // 199元：开通会员资格，余额不变
	PrechargeMinimum = 100000 // 1000元：预充值门槛，到账1200元（1000+200赠送）
)

// RechargeByPlan 按档位充值/升级
// planID: 目标档位ID
// 支持补差价升级：已充值金额可抵扣目标档位门槛
func (s *UserService) RechargeByPlan(userID uint, planID int, channel string) (*model.RechargeRecord, error) {
	// 查找目标档位
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

	// 获取用户信息和累计充值金额
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	totalRecharged, err := s.rechargeRepo.GetTotalAmountByUserID(userID)
	if err != nil {
		totalRecharged = 0
	}

	// 计算需要支付的差价
	payAmount := targetPlan.Amount - totalRecharged
	if payAmount <= 0 {
		return nil, fmt.Errorf("您已累计充值 ¥%.2f，已达到或超过该档位", float64(totalRecharged)/100)
	}

	// 计算实际到账金额（目标档位到账额 - 当前余额）
	creditAmount := targetPlan.FinalAmount - user.Balance
	if creditAmount < 0 {
		creditAmount = 0
	}

	// 创建充值记录（记录实际支付的差价和实际到账）
	// ⚠️ 支付接口开发中：记录状态为 pending，等微信支付回调后再改为 success
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

	// TODO: 调用微信支付接口，获取支付参数返回给前端
	// 前端支付成功后，微信支付回调更新记录状态为 success，并执行以下逻辑：

	// 更新余额（到账金额）
	// if creditAmount > 0 {
	// 	if err := s.repo.UpdateBalance(userID, creditAmount); err != nil {
	// 		fmt.Printf("update balance failed: %v\n", err)
	// 	}
	// }
	// 更新会员等级
	// if targetPlan.MemberLevel > user.MemberLevel {
	// 	if err := s.repo.UpdateMemberLevel(userID, targetPlan.MemberLevel); err != nil {
	// 		fmt.Printf("upgrade member level failed: %v\n", err)
	// 	}
	// }

	return record, nil
}

// Recharge 兼容旧接口（按金额充值，走plan匹配逻辑）
func (s *UserService) Recharge(userID uint, amount int, channel string) (*model.RechargeRecord, error) {
	if amount <= 0 {
		return nil, errors.New("充值金额无效")
	}
	// 尝试匹配档位
	for _, plan := range model.RechargePlans {
		if plan.Amount == amount {
			return s.RechargeByPlan(userID, plan.ID, channel)
		}
	}
	// 未匹配到标准档位，走旧逻辑（仅兼容1000元及以上）
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
		// TODO: 微信支付回调后更新余额和会员等级
		return record, nil
	}
	return nil, errors.New("不支持的充值金额，请选择标准充值档位")
}

// GetRechargePlans 获取充值档位列表
func (s *UserService) GetRechargePlans() []model.RechargePlan {
	return model.RechargePlans
}

// GetUpgradeInfo 获取用户当前升级信息
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

	// 查找当前已达到的最高档位
	var currentPlan *model.RechargePlan
	for i := len(model.RechargePlans) - 1; i >= 0; i-- {
		if totalRecharged >= model.RechargePlans[i].Amount {
			currentPlan = &model.RechargePlans[i]
			break
		}
	}
	info.CurrentPlan = currentPlan

	// 查找下一个可升级档位
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

func (s *UserService) DeductBalance(userID uint, amount int, orderNo string) (map[string]interface{}, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if user.Balance < amount {
		return nil, fmt.Errorf("余额不足")
	}

	if err := s.repo.UpdateBalance(userID, -amount); err != nil {
		return nil, err
	}

	log := &model.BalanceLog{
		UserID:  userID,
		Type:    "deduct",
		Amount:  -amount,
		OrderNo: orderNo,
		Remark:  "订单扣款",
	}
	if err := s.balanceRepo.Create(log); err != nil {
		return nil, err
	}

	newUser, _ := s.repo.FindByID(userID)
	balanceBefore := user.Balance
	balanceAfter := 0
	if newUser != nil {
		balanceAfter = newUser.Balance
	}

	return map[string]interface{}{
		"balance_before": balanceBefore,
		"balance_after":  balanceAfter,
	}, nil
}

func (s *UserService) RefundBalance(userID uint, amount int, orderNo, remark string) (map[string]interface{}, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	if err := s.repo.UpdateBalance(userID, amount); err != nil {
		return nil, err
	}

	log := &model.BalanceLog{
		UserID:  userID,
		Type:    "refund",
		Amount:  amount,
		OrderNo: orderNo,
		Remark:  remark,
	}
	if err := s.balanceRepo.Create(log); err != nil {
		return nil, err
	}

	balanceBefore := user.Balance
	newUser, _ := s.repo.FindByID(userID)
	balanceAfter := 0
	if newUser != nil {
		balanceAfter = newUser.Balance
	}

	return map[string]interface{}{
		"balance_before": balanceBefore,
		"balance_after":  balanceAfter,
	}, nil
}

func (s *UserService) ListPets(userID uint) ([]model.PetProfile, error) {
	return s.petRepo.ListByUserID(userID)
}

func (s *UserService) AddPet(userID uint, name, breed string, weight float64, birthday *string) (*model.PetProfile, error) {
	if name == "" {
		return nil, errors.New("宠物名不能为空")
	}
	pet := &model.PetProfile{
		UserID:   userID,
		Name:     name,
		Breed:    breed,
		Weight:   weight,
		Birthday: birthday,
	}
	if err := s.petRepo.Create(pet); err != nil {
		return nil, err
	}
	return pet, nil
}

func (s *UserService) CheckAndUpgradeMemberLevel(userID uint) error {
	// 当前规则：仅通过充值成为会员，消费不自动升级
	return nil
}

// GetTotalRechargedAmount 获取用户累计充值金额（已完成的充值记录）
func (s *UserService) GetTotalRechargedAmount(userID uint) (int, error) {
	return s.rechargeRepo.GetTotalAmountByUserID(userID)
}

// GetRechargeRecords 获取用户充值记录列表
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
