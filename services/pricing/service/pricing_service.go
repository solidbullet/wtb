package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/wtb-ordering/services/pricing/model"
	"github.com/wtb-ordering/services/pricing/repository"
)

type PricingService struct {
	ruleRepo      *repository.PriceRuleRepo
	promoRepo     *repository.PromotionRepo
	comboRepo     *repository.ComboRepo
	rechargeRepo  *repository.RechargePlanRepo
}

func NewPricingService(ruleRepo *repository.PriceRuleRepo, promoRepo *repository.PromotionRepo, comboRepo *repository.ComboRepo, rechargeRepo *repository.RechargePlanRepo) *PricingService {
	return &PricingService{
		ruleRepo:     ruleRepo,
		promoRepo:    promoRepo,
		comboRepo:    comboRepo,
		rechargeRepo: rechargeRepo,
	}
}

func (s *PricingService) CalculateOrderPrice(userLevel int, items []struct {
	DishID   uint `json:"dish_id"`
	Quantity int  `json:"quantity"`
}) (map[string]interface{}, error) {
	var resultItems []map[string]interface{}
	var totalAmount int

	for _, item := range items {
		price, priceType, err := s.getDishPrice(item.DishID, userLevel)
		if err != nil {
			return nil, fmt.Errorf("菜品ID不存在")
		}
		subtotal := price * item.Quantity
		totalAmount += subtotal
		resultItems = append(resultItems, map[string]interface{}{
			"dish_id":    item.DishID,
			"quantity":   item.Quantity,
			"unit_price": price,
			"subtotal":   subtotal,
			"price_type": priceType,
		})
	}

	discountAmount := 0
	var appliedPromos []map[string]interface{}
	promos, _ := s.promoRepo.ListActive()
	for _, promo := range promos {
		discount := s.applyPromotion(promo, totalAmount)
		if discount > 0 {
			discountAmount += discount
			appliedPromos = append(appliedPromos, map[string]interface{}{
				"promotion_id": promo.ID,
				"name":         promo.Name,
				"discount":     discount,
			})
		}
	}

	finalAmount := totalAmount - discountAmount
	if finalAmount < 0 {
		finalAmount = 0
	}

	return map[string]interface{}{
		"items":               resultItems,
		"total_amount":        totalAmount,
		"discount_amount":     discountAmount,
		"final_amount":        finalAmount,
		"applied_promotions":  appliedPromos,
	}, nil
}

func (s *PricingService) getDishPrice(dishID uint, userLevel int) (int, string, error) {
	rules, err := s.ruleRepo.ListByDishID(dishID)
	if err != nil || len(rules) == 0 {
		// fallback: 调用 menu-service 获取基础价（简化）
		return 0, "", errors.New("菜品价格未配置")
	}

	now := time.Now()
	for _, r := range rules {
		if r.RuleType == "time_slot" && r.StartTime != nil && r.EndTime != nil {
			st, _ := time.Parse("15:04", *r.StartTime)
			et, _ := time.Parse("15:04", *r.EndTime)
			cur, _ := time.Parse("15:04", now.Format("15:04"))
			if cur.After(st) || cur.Equal(st) && cur.Before(et) {
				return r.Price, "time_slot", nil
			}
		}
	}

	// 会员价优先
	if userLevel >= 1 {
		for _, r := range rules {
			if r.RuleType == "member" {
				return r.Price, "member", nil
			}
		}
	}

	// 普通价
	for _, r := range rules {
		if r.RuleType == "normal" {
			return r.Price, "normal", nil
		}
	}
	return 0, "", errors.New("未找到有效价格")
}

func (s *PricingService) applyPromotion(promo model.Promotion, totalAmount int) int {
	var config map[string]interface{}
	json.Unmarshal([]byte(promo.ConfigJSON), &config)

	switch promo.Type {
	case "full_reduction":
		threshold, _ := config["threshold"].(float64)
		reduce, _ := config["reduce"].(float64)
		if totalAmount >= int(threshold) {
			return int(reduce)
		}
	case "discount":
		rate, _ := config["rate"].(float64)
		maxReduce, _ := config["max_reduce"].(float64)
		discount := int(float64(totalAmount) * (1 - rate))
		if maxReduce > 0 && discount > int(maxReduce) {
			discount = int(maxReduce)
		}
		return discount
	}
	return 0
}

func (s *PricingService) GetDishPrice(dishID uint, userLevel int) (map[string]interface{}, error) {
	price, priceType, err := s.getDishPrice(dishID, userLevel)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"dish_id":        dishID,
		"price":          price,
		"price_type":     priceType,
		"original_price": price,
	}, nil
}

func (s *PricingService) ListPromotions() ([]model.Promotion, error) {
	return s.promoRepo.ListActive()
}

func (s *PricingService) CreatePriceRule(rule *model.PriceRule) error {
	return s.ruleRepo.Create(rule)
}

func (s *PricingService) CreatePromotion(promo *model.Promotion) error {
	return s.promoRepo.Create(promo)
}

func (s *PricingService) CreateCombo(combo *model.Combo) error {
	return s.comboRepo.Create(combo)
}

func (s *PricingService) ListRechargePlans() ([]model.RechargePlan, error) {
	return s.rechargeRepo.ListActive()
}

func (s *PricingService) CreateRechargePlan(plan *model.RechargePlan) error {
	return s.rechargeRepo.Create(plan)
}

func (s *PricingService) UpdateRechargePlan(id uint, plan *model.RechargePlan) error {
	return s.rechargeRepo.Update(id, plan)
}

func (s *PricingService) DeleteRechargePlan(id uint) error {
	return s.rechargeRepo.Delete(id)
}
