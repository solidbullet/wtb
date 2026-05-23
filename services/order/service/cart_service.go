package service

import (
	"encoding/json"
	"fmt"
	"sync"
)

type CartItem struct {
	UserID    uint   `json:"user_id"`
	DishID    uint   `json:"dish_id"`
	DishName  string `json:"dish_name"`
	Quantity  int    `json:"quantity"`
	UnitPrice int    `json:"unit_price"`
	Remark    string `json:"remark"`
}

type CartService struct {
	mu  sync.RWMutex
	mem map[string]map[string]string // seatID -> dishID -> json
}

func NewCartService() *CartService {
	return &CartService{mem: make(map[string]map[string]string)}
}

func (s *CartService) Add(seatID string, item CartItem) error {
	field := fmt.Sprintf("%d", item.DishID)
	if item.Quantity <= 0 {
		return s.Remove(seatID, item.DishID)
	}
	data, _ := json.Marshal(item)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.mem[seatID] == nil {
		s.mem[seatID] = make(map[string]string)
	}
	s.mem[seatID][field] = string(data)
	return nil
}

func (s *CartService) List(seatID string) ([]CartItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var items []CartItem
	for _, v := range s.mem[seatID] {
		var item CartItem
		json.Unmarshal([]byte(v), &item)
		items = append(items, item)
	}
	return items, nil
}

func (s *CartService) Update(seatID string, dishID uint, quantity int, remark string) error {
	if quantity <= 0 {
		return s.Remove(seatID, dishID)
	}
	field := fmt.Sprintf("%d", dishID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.mem[seatID] == nil {
		return fmt.Errorf("cart not found")
	}
	data := s.mem[seatID][field]
	var item CartItem
	json.Unmarshal([]byte(data), &item)
	item.Quantity = quantity
	item.Remark = remark
	data2, _ := json.Marshal(item)
	s.mem[seatID][field] = string(data2)
	return nil
}

func (s *CartService) Remove(seatID string, dishID uint) error {
	field := fmt.Sprintf("%d", dishID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.mem[seatID] != nil {
		delete(s.mem[seatID], field)
	}
	return nil
}

func (s *CartService) Clear(seatID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.mem, seatID)
	return nil
}
