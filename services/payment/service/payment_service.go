package service
import ("errors"; "fmt"; "time"; "github.com/wtb-ordering/services/payment/model"; "github.com/wtb-ordering/services/payment/repository")
type PaymentService struct {
	orderRepo    *repository.PaymentOrderRepo; recordRepo *repository.PaymentRecordRepo; refundRepo *repository.RefundRecordRepo; rechargeRepo *repository.RechargeOrderRepo }
func NewPaymentService(orderRepo *repository.PaymentOrderRepo, recordRepo *repository.PaymentRecordRepo, refundRepo *repository.RefundRecordRepo, rechargeRepo *repository.RechargeOrderRepo) *PaymentService {
	return &PaymentService{orderRepo: orderRepo, recordRepo: recordRepo, refundRepo: refundRepo, rechargeRepo: rechargeRepo} }
func (s *PaymentService) CreatePayment(orderNo string, userID uint, amount int, channel string) (*model.PaymentOrder, error) {
	outTradeNo := fmt.Sprintf("PAY%s", time.Now().Format("20060102150405")); po := &model.PaymentOrder{OrderNo: orderNo, OutTradeNo: outTradeNo, UserID: userID, Amount: amount, Channel: channel, Status: "pending"}; return po, s.orderRepo.Create(po) }
func (s *PaymentService) BalancePay(outTradeNo string) error {
	po, err := s.orderRepo.FindByOutTradeNo(outTradeNo); if err != nil { return errors.New("支付单不存在") }; if po.Status != "pending" { return errors.New("支付单状态错误") }
	if err := s.orderRepo.UpdateStatus(po.ID, "paid"); err != nil { return err }; now := time.Now(); return s.recordRepo.Create(&model.PaymentRecord{PaymentOrderID: po.ID, Channel: "balance", Amount: po.Amount, PaidAt: &now}) }
func (s *PaymentService) WxCallback(outTradeNo, transactionID string) error {
	po, err := s.orderRepo.FindByOutTradeNo(outTradeNo); if err != nil { return err }; if err := s.orderRepo.UpdateStatus(po.ID, "paid"); err != nil { return err }; now := time.Now(); return s.recordRepo.Create(&model.PaymentRecord{PaymentOrderID: po.ID, Channel: "wxpay", Amount: po.Amount, TransactionID: transactionID, PaidAt: &now}) }
func (s *PaymentService) Refund(outTradeNo string, amount int, reason string) error {
	po, err := s.orderRepo.FindByOutTradeNo(outTradeNo); if err != nil { return err }; refundNo := fmt.Sprintf("REF%s", time.Now().Format("20060102150405")); return s.refundRepo.Create(&model.RefundRecord{PaymentOrderID: po.ID, RefundNo: refundNo, Amount: amount, Reason: reason}) }
func (s *PaymentService) Recharge(userID uint, amount int) (*model.RechargeOrder, error) {
	gifted := amount / 5; finalAmount := amount + gifted; ro := &model.RechargeOrder{UserID: userID, Amount: amount, GiftedAmount: gifted, FinalAmount: finalAmount, Status: "pending"}; return ro, s.rechargeRepo.Create(ro) }
func (s *PaymentService) Query(outTradeNo string) (*model.PaymentOrder, error) { return s.orderRepo.FindByOutTradeNo(outTradeNo) }
