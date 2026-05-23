package service
import ("errors"; "fmt"; "github.com/wtb-ordering/services/points/model"; "github.com/wtb-ordering/services/points/repository")
type PointsService struct {
	upRepo *repository.UserPointsRepo; logRepo *repository.PointsLogRepo; goodsRepo *repository.ExchangeGoodsRepo; orderRepo *repository.ExchangeOrderRepo }
func NewPointsService(upRepo *repository.UserPointsRepo, logRepo *repository.PointsLogRepo, goodsRepo *repository.ExchangeGoodsRepo, orderRepo *repository.ExchangeOrderRepo) *PointsService {
	return &PointsService{upRepo: upRepo, logRepo: logRepo, goodsRepo: goodsRepo, orderRepo: orderRepo} }
func (s *PointsService) GetAccount(userID uint) (*model.UserPoints, error) { up, err := s.upRepo.FindByUserID(userID); if err != nil { return &model.UserPoints{UserID: userID}, nil }; return up, nil }
func (s *PointsService) GetLogs(userID uint, page, pageSize int) ([]model.PointsLog, int64, error) { return s.logRepo.ListByUserID(userID, page, pageSize) }
func (s *PointsService) ListGoods() ([]model.ExchangeGoods, error) { return s.goodsRepo.ListActive() }
func (s *PointsService) CreateGoods(g *model.ExchangeGoods) error { return s.goodsRepo.Create(g) }
func (s *PointsService) Exchange(userID, goodsID uint) (*model.ExchangeOrder, error) {
	goods, err := s.goodsRepo.FindByID(goodsID); if err != nil { return nil, errors.New("商品不存在") }
	if goods.Stock <= 0 { return nil, errors.New("库存不足") }
	up, err := s.upRepo.FindByUserID(userID)
	if err != nil {
		up = &model.UserPoints{UserID: userID}
	}
	available := up.TotalPoints - up.UsedPoints - up.FrozenPoints
	if available < goods.PointsPrice { return nil, errors.New("积分不足") }
	if err := s.upRepo.UpdatePoints(userID, -goods.PointsPrice); err != nil { return nil, err }
	if err := s.goodsRepo.DeductStock(goodsID, 1); err != nil { return nil, err }
	order := &model.ExchangeOrder{UserID: userID, GoodsID: goodsID, PointsCost: goods.PointsPrice, Status: "pending"}
	if err := s.orderRepo.Create(order); err != nil { return nil, err }
	log := &model.PointsLog{UserID: userID, Type: "exchange", Points: -goods.PointsPrice, SourceID: fmt.Sprintf("%d", order.ID), Remark: "积分兑换"}
	s.logRepo.Create(log); return order, nil }
func (s *PointsService) GrantPoints(userID uint, amount int, sourceID, remark string) error {
	if _, err := s.upRepo.FindByUserID(userID); err != nil { s.upRepo.Create(&model.UserPoints{UserID: userID, TotalPoints: amount}) } else { s.upRepo.UpdatePoints(userID, amount) }
	return s.logRepo.Create(&model.PointsLog{UserID: userID, Type: "gain", Points: amount, SourceID: sourceID, Remark: remark}) }
