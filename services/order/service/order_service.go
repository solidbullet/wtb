package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"github.com/wtb-ordering/services/order/model"
	"github.com/wtb-ordering/services/order/repository"
)

type OrderService struct {
	orderRepo *repository.OrderRepo
	itemRepo  *repository.OrderItemRepo
	logRepo   *repository.OrderStatusLogRepo
	cartSvc   *CartService
}

func NewOrderService(
	orderRepo *repository.OrderRepo,
	itemRepo *repository.OrderItemRepo,
	logRepo *repository.OrderStatusLogRepo,
	cartSvc *CartService,
) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		itemRepo:  itemRepo,
		logRepo:   logRepo,
		cartSvc:   cartSvc,
	}
}

func (s *OrderService) CreateOrder(seatID string, userID uint, remark string) (*model.Order, error) {
	items, err := s.cartSvc.List(seatID)
	if err != nil || len(items) == 0 {
		return nil, errors.New("购物车为空")
	}

	orderNo := fmt.Sprintf("WTB%s", time.Now().Format("20060102150405"))
	total := 0
	for _, it := range items {
		total += it.UnitPrice * it.Quantity
	}

	order := &model.Order{
		OrderNo:    orderNo,
		SeatID:     seatID,
		UserID:     userID,
		Status:     "pending",
		TotalAmount: total,
		PayAmount:  total,
		Remark:     remark,
	}

	err = s.orderRepo.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for _, it := range items {
			oi := &model.OrderItem{
				OrderID:   order.ID,
				DishID:    it.DishID,
				DishName:  it.DishName,
				Quantity:  it.Quantity,
				UnitPrice: it.UnitPrice,
			}
			if err := tx.Create(oi).Error; err != nil {
				return err
			}
		}
		logEntry := &model.OrderStatusLog{
			OrderID:  order.ID,
			ToStatus: "pending",
		}
		return tx.Create(logEntry).Error
	})
	if err != nil {
		return nil, err
	}

	s.cartSvc.Clear(seatID)
	return order, nil
}

func (s *OrderService) GetOrderStatus(orderID uint) (*model.Order, []model.OrderStatusLog, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, nil, err
	}
	logs, err := s.logRepo.ListByOrder(orderID)
	if err != nil {
		return order, nil, err
	}
	return order, logs, nil
}

func (s *OrderService) ListOrders(userID uint, status string, page, pageSize int) ([]model.Order, int64, error) {
	return s.orderRepo.ListByUser(userID, status, page, pageSize)
}

func (s *OrderService) ListAllOrders(status string, page, pageSize int) ([]model.Order, int64, error) {
	return s.orderRepo.ListAll(status, page, pageSize)
}

func (s *OrderService) ListTodayPaidOrders() ([]model.Order, error) {
	return s.orderRepo.ListTodayPaidOrders()
}

func (s *OrderService) UpdateStatus(orderID uint, from, to, operator string) error {
	if err := s.orderRepo.UpdateStatus(orderID, to); err != nil {
		return err
	}
	return s.logRepo.Create(&model.OrderStatusLog{
		OrderID:    orderID,
		FromStatus: from,
		ToStatus:   to,
		Operator:   operator,
	})
}
