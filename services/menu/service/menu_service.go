package service

import (
	"errors"
	"time"

	"github.com/wtb-ordering/services/menu/model"
	"github.com/wtb-ordering/services/menu/repository"
)

type MenuService struct {
	catRepo    *repository.CategoryRepo
	dishRepo   *repository.DishRepo
	priceRepo  *repository.DishPriceRepo
	stockRepo  *repository.DishStockRepo
}

func NewMenuService(catRepo *repository.CategoryRepo, dishRepo *repository.DishRepo, priceRepo *repository.DishPriceRepo, stockRepo *repository.DishStockRepo) *MenuService {
	return &MenuService{catRepo: catRepo, dishRepo: dishRepo, priceRepo: priceRepo, stockRepo: stockRepo}
}

func (s *MenuService) BuildCategoryTree() ([]map[string]interface{}, error) {
	cats, err := s.catRepo.ListAll()
	if err != nil {
		return nil, err
	}
	var roots []map[string]interface{}
	childrenMap := make(map[uint][]map[string]interface{})
	for _, cat := range cats {
		node := map[string]interface{}{
			"id":         cat.ID,
			"name":       cat.Name,
			"parent_id":  cat.ParentID,
			"sort_order": cat.SortOrder,
			"children":   []map[string]interface{}{},
		}
		if cat.ParentID == 0 {
			roots = append(roots, node)
		} else {
			childrenMap[cat.ParentID] = append(childrenMap[cat.ParentID], node)
		}
	}
	for i := range roots {
		roots[i]["children"] = childrenMap[roots[i]["id"].(uint)]
	}
	return roots, nil
}

func (s *MenuService) ListDishes(categoryID uint, page, pageSize int) ([]model.Dish, int64, error) {
	return s.dishRepo.ListByCategory(categoryID, page, pageSize)
}

func (s *MenuService) ListDishesByTags(tag string, page, pageSize int) ([]model.Dish, int64, error) {
	return s.dishRepo.ListByTags(tag, page, pageSize)
}

func (s *MenuService) GetDish(id uint) (*model.Dish, []model.DishPrice, *model.DishStock, error) {
	dish, err := s.dishRepo.FindByID(id)
	if err != nil {
		return nil, nil, nil, err
	}
	prices, _ := s.priceRepo.ListByDishID(id)
	stock, _ := s.stockRepo.FindByDishAndDate(id, time.Now())
	return dish, prices, stock, nil
}

func (s *MenuService) GetDishPrices(dishID uint) ([]model.DishPrice, error) {
	return s.priceRepo.ListByDishID(dishID)
}

func (s *MenuService) GetDishStock(dishID uint) (*model.DishStock, error) {
	return s.stockRepo.FindByDishAndDate(dishID, time.Now())
}

func (s *MenuService) AddDishPrice(dishID uint, priceType string, price int) error {
	return s.priceRepo.Create(&model.DishPrice{DishID: dishID, PriceType: priceType, Price: price})
}

func (s *MenuService) AddDishStock(dishID uint, dailyLimit int) error {
	return s.stockRepo.Create(&model.DishStock{DishID: dishID, Date: time.Now(), DailyLimit: dailyLimit})
}

func (s *MenuService) SearchDishes(q string) ([]model.Dish, error) {
	return s.dishRepo.Search(q)
}

func (s *MenuService) BatchDishes(ids []uint) ([]model.Dish, error) {
	return s.dishRepo.BatchByIDs(ids)
}

func (s *MenuService) CreateCategory(name string, parentID uint, sortOrder int) (*model.Category, error) {
	if name == "" {
		return nil, errors.New("分类名不能为空")
	}
	cat := &model.Category{Name: name, ParentID: parentID, SortOrder: sortOrder}
	return cat, s.catRepo.Create(cat)
}

func (s *MenuService) UpdateCategory(id uint, name string, parentID uint, sortOrder int) error {
	return s.catRepo.Update(id, &model.Category{Name: name, ParentID: parentID, SortOrder: sortOrder})
}

func (s *MenuService) DeleteCategory(id uint) error {
	return s.catRepo.Delete(id)
}

func (s *MenuService) CreateDish(categoryID uint, name, subtitle, description, images, tags string, prices []model.DishPrice) (*model.Dish, error) {
	if name == "" {
		return nil, errors.New("菜品名不能为空")
	}
	dish := &model.Dish{CategoryID: categoryID, Name: name, Subtitle: subtitle, Description: description, Images: images, Tags: tags}
	if err := s.dishRepo.Create(dish); err != nil {
		return nil, err
	}
	for _, p := range prices {
		p.DishID = dish.ID
		s.priceRepo.Create(&p)
	}
	return dish, nil
}

func (s *MenuService) UpdateDish(id uint, categoryID uint, name, subtitle, description, images, tags string) error {
	return s.dishRepo.Update(id, &model.Dish{CategoryID: categoryID, Name: name, Subtitle: subtitle, Description: description, Images: images, Tags: tags})
}

func (s *MenuService) DeleteDish(id uint) error {
	if err := s.priceRepo.DeleteByDishID(id); err != nil {
		return err
	}
	if err := s.stockRepo.DeleteByDishID(id); err != nil {
		return err
	}
	return s.dishRepo.Delete(id)
}

func (s *MenuService) SetStock(dishID uint, date time.Time, dailyLimit int) error {
	stock := &model.DishStock{DishID: dishID, Date: date, DailyLimit: dailyLimit, SoldCount: 0}
	return s.stockRepo.Create(stock)
}
