package service
import ("errors"; "fmt"; "time"; "github.com/wtb-ordering/services/order/model"; "github.com/wtb-ordering/services/order/repository")
type OrderService struct {
	orderRepo    *repository.OrderRepo; itemRepo *repository.OrderItemRepo; logRepo *repository.OrderStatusLogRepo; cartSvc *CartService }
func NewOrderService(orderRepo *repository.OrderRepo, itemRepo *repository.OrderItemRepo, logRepo *repository.OrderStatusLogRepo, cartSvc *CartService) *OrderService {
	return &OrderService{orderRepo: orderRepo, itemRepo: itemRepo, logRepo: logRepo, cartSvc: cartSvc} }
func (s *OrderService) CreateOrder(seatID string, userID uint, remark string) (*model.Order, error) {
	items, err := s.cartSvc.List(seatID); if err != nil || len(items) == 0 { return nil, errors.New("购物车为空") }
	orderNo := fmt.Sprintf("WTB%s", time.Now().Format("20060102150405")); total := 0
	for _, it := range items { total += it.UnitPrice * it.Quantity }
	order := &model.Order{OrderNo: orderNo, SeatID: seatID, UserID: userID, Status: "pending", TotalAmount: total, PayAmount: total, Remark: remark}
	if err := s.orderRepo.Create(order); err != nil { return nil, err }
	for _, it := range items { s.itemRepo.Create(&model.OrderItem{OrderID: order.ID, DishID: it.DishID, DishName: it.DishName, Quantity: it.Quantity, UnitPrice: it.UnitPrice}) }
	s.logRepo.Create(&model.OrderStatusLog{OrderID: order.ID, ToStatus: "pending"}); s.cartSvc.Clear(seatID); return order, nil }
func (s *OrderService) GetOrderStatus(orderID uint) (*model.Order, []model.OrderStatusLog, error) {
	order, err := s.orderRepo.FindByID(orderID); if err != nil { return nil, nil, err }; logs, _ := s.logRepo.ListByOrder(orderID); return order, logs, nil }
func (s *OrderService) ListOrders(userID uint, page, pageSize int) ([]model.Order, int64, error) { return s.orderRepo.ListByUser(userID, page, pageSize) }
func (s *OrderService) ListAllOrders(page, pageSize int) ([]model.Order, int64, error) { return s.orderRepo.ListAll(page, pageSize) }
func (s *OrderService) UpdateStatus(orderID uint, from, to, operator string) error {
	if err := s.orderRepo.UpdateStatus(orderID, to); err != nil { return err }; return s.logRepo.Create(&model.OrderStatusLog{OrderID: orderID, FromStatus: from, ToStatus: to, Operator: operator}) }
