package service

import (
	"github.com/wtb-ordering/services/order/model"
	"github.com/wtb-ordering/services/order/repository"
)

type CartService struct {
	repo *repository.CartRepo
}

func NewCartService(repo *repository.CartRepo) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) Add(seatID string, item model.CartItem) error {
	if item.Quantity <= 0 {
		return s.Remove(seatID, item.DishID)
	}
	return s.repo.Upsert(seatID, item)
}

func (s *CartService) List(seatID string) ([]model.CartItem, error) {
	return s.repo.ListBySeat(seatID)
}

func (s *CartService) Update(seatID string, dishID uint, quantity int, remark string) error {
	if quantity <= 0 {
		return s.Remove(seatID, dishID)
	}
	items, err := s.repo.ListBySeat(seatID)
	if err != nil {
		return err
	}
	for _, it := range items {
		if it.DishID == dishID {
			it.Quantity = quantity
			it.Remark = remark
			return s.repo.Upsert(seatID, it)
		}
	}
	return nil
}

func (s *CartService) Remove(seatID string, dishID uint) error {
	return s.repo.Remove(seatID, dishID)
}

func (s *CartService) Clear(seatID string) error {
	return s.repo.Clear(seatID)
}
