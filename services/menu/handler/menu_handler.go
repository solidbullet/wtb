package handler

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wtb-ordering/pkg/response"
	"github.com/wtb-ordering/services/menu/model"
	"github.com/wtb-ordering/services/menu/service"
)

type MenuHandler struct {
	svc *service.MenuService
}

func NewMenuHandler(svc *service.MenuService) *MenuHandler {
	return &MenuHandler{svc: svc}
}

// GetCategories GET /api/menu/categories
func (h *MenuHandler) GetCategories(c *gin.Context) {
	tree, err := h.svc.BuildCategoryTree()
	if err != nil {
		response.Error(c, 50001, "获取分类失败")
		return
	}
	response.Success(c, tree)
}

// ListDishes GET /api/menu/dishes
func (h *MenuHandler) ListDishes(c *gin.Context) {
	catID, _ := strconv.Atoi(c.Query("category_id"))
	tags := c.Query("tags")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	var list []model.Dish
	var total int64
	var err error
	if tags != "" {
		list, total, err = h.svc.ListDishesByTags(tags, page, pageSize)
	} else {
		list, total, err = h.svc.ListDishes(uint(catID), page, pageSize)
	}
	if err != nil {
		response.Error(c, 50001, "获取菜品列表失败")
		return
	}
	// 补充价格信息
	var result []gin.H
	for _, d := range list {
		prices, _ := h.svc.GetDishPrices(d.ID)
		price := 0
		memberPrice := 0
		for _, p := range prices {
			if p.PriceType == "standard" {
				price = p.Price
			} else if p.PriceType == "member" {
				memberPrice = p.Price
			}
		}
		if price == 0 && len(prices) > 0 {
			price = prices[0].Price
		}
		stockInfo, _ := h.svc.GetDishStock(d.ID)
		stock := 0
		if stockInfo != nil {
			if stockInfo.DailyLimit == -1 {
				stock = -1
			} else {
				stock = stockInfo.DailyLimit - stockInfo.SoldCount
			}
		}
		result = append(result, gin.H{
			"id":           d.ID,
			"category_id":  d.CategoryID,
			"name":         d.Name,
			"subtitle":     d.Subtitle,
			"description":  d.Description,
			"images":       d.Images,
			"tags":         d.Tags,
			"status":       d.Status,
			"price":        price,
			"member_price": memberPrice,
			"stock":        stock,
			"created_at":   d.CreatedAt,
			"updated_at":   d.UpdatedAt,
		})
	}
	response.SuccessPage(c, total, page, pageSize, result)
}

// GetDish GET /api/menu/dish/:id
func (h *MenuHandler) GetDish(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	dish, prices, stock, err := h.svc.GetDish(uint(id))
	if err != nil {
		response.Error(c, 50001, "菜品不存在")
		return
	}
	price := 0
	memberPrice := 0
	for _, p := range prices {
		if p.PriceType == "standard" {
			price = p.Price
		} else if p.PriceType == "member" {
			memberPrice = p.Price
		}
	}
	if price == 0 && len(prices) > 0 {
		price = prices[0].Price
	}
	response.Success(c, gin.H{
		"id":           dish.ID,
		"category_id":  dish.CategoryID,
		"name":         dish.Name,
		"subtitle":     dish.Subtitle,
		"description":  dish.Description,
		"images":       dish.Images,
		"tags":         dish.Tags,
		"price":        price,
		"member_price": memberPrice,
		"prices":       prices,
		"stock":        stock,
		"status":       dish.Status,
		"created_at":   dish.CreatedAt,
	})
}

// SearchDishes GET /api/menu/search
func (h *MenuHandler) SearchDishes(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.Success(c, []model.Dish{})
		return
	}
	list, err := h.svc.SearchDishes(q)
	if err != nil {
		response.Error(c, 50001, "搜索失败")
		return
	}
	response.Success(c, list)
}

// BatchDishes POST /api/menu/dishes/batch
func (h *MenuHandler) BatchDishes(c *gin.Context) {
	var req struct {
		DishIDs []uint `json:"dish_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	list, err := h.svc.BatchDishes(req.DishIDs)
	if err != nil {
		response.Error(c, 50001, "批量查询失败")
		return
	}
	var result []gin.H
	for _, d := range list {
		prices, _ := h.svc.GetDishPrices(d.ID)
		price := 0
		memberPrice := 0
		for _, p := range prices {
			if p.PriceType == "standard" && price == 0 {
				price = p.Price
			}
			if p.PriceType == "member" {
				memberPrice = p.Price
			}
		}
		if price == 0 && len(prices) > 0 {
			price = prices[0].Price
		}
		result = append(result, gin.H{
			"id":           d.ID,
			"name":         d.Name,
			"price":        price,
			"member_price": memberPrice,
			"images":       d.Images,
		})
	}
	response.Success(c, result)
}

// CreateCategory POST /api/menu/admin/category
func (h *MenuHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name      string `json:"name"`
		ParentID  uint   `json:"parent_id"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	cat, err := h.svc.CreateCategory(req.Name, req.ParentID, req.SortOrder)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, cat)
}

// UpdateCategory PUT /api/menu/admin/category/:id
func (h *MenuHandler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Name      string `json:"name"`
		ParentID  uint   `json:"parent_id"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	if err := h.svc.UpdateCategory(uint(id), req.Name, req.ParentID, req.SortOrder); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id})
}

// DeleteCategory DELETE /api/menu/admin/category/:id
func (h *MenuHandler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteCategory(uint(id)); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id})
}

// CreateDish POST /api/menu/admin/dish
func (h *MenuHandler) CreateDish(c *gin.Context) {
	var req struct {
		CategoryID  uint              `json:"category_id"`
		Name        string            `json:"name"`
		Subtitle    string            `json:"subtitle"`
		Description string            `json:"description"`
		Images      string            `json:"images"`
		Tags        string            `json:"tags"`
		Price       int               `json:"price"`
		Stock       int               `json:"stock"`
		Prices      []model.DishPrice `json:"prices"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	dish, err := h.svc.CreateDish(req.CategoryID, req.Name, req.Subtitle, req.Description, req.Images, req.Tags, req.Prices)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	// 如果前端传了 price 但没有传 prices 数组，自动创建标准价格
	if req.Price > 0 && len(req.Prices) == 0 {
		h.svc.AddDishPrice(dish.ID, "standard", req.Price)
	}
	// 如果前端传了 stock，自动创建库存
	if req.Stock > 0 {
		h.svc.AddDishStock(dish.ID, req.Stock)
	}
	response.Success(c, dish)
}

// UpdateDish PUT /api/menu/admin/dish/:id
func (h *MenuHandler) UpdateDish(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		CategoryID  uint   `json:"category_id"`
		Name        string `json:"name"`
		Subtitle    string `json:"subtitle"`
		Description string `json:"description"`
		Images      string `json:"images"`
		Tags        string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	if err := h.svc.UpdateDish(uint(id), req.CategoryID, req.Name, req.Subtitle, req.Description, req.Images, req.Tags); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id})
}

// DeleteDish DELETE /api/menu/admin/dish/:id
func (h *MenuHandler) DeleteDish(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	// 先查询菜品，获取图片路径
	dish, _, _, err := h.svc.GetDish(uint(id))
	imagePath := ""
	if err == nil && dish != nil && dish.Images != "" {
		imagePath = dish.Images
	}

	if err := h.svc.DeleteDish(uint(id)); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	// 删除对应的图片文件
	if imagePath != "" {
		filename := strings.TrimPrefix(imagePath, "/images/")
		if filename != imagePath {
			filePath := filepath.Join("..", "miniprogram", "images", filename)
			os.Remove(filePath)
		}
	}

	response.Success(c, gin.H{"id": id})
}

// UploadImage POST /api/menu/admin/upload
func (h *MenuHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, 40001, "请选择图片文件")
		return
	}
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".png"
	}
	filename := fmt.Sprintf("dish_%d_%d%s", time.Now().Unix(), rand.Intn(10000), ext)
	dst := filepath.Join("..", "miniprogram", "images", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.Error(c, 50001, "保存图片失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"url":      "/images/" + filename,
		"filename": filename,
	})
}

// SetStock POST /api/menu/admin/stock
func (h *MenuHandler) SetStock(c *gin.Context) {
	var req struct {
		DishID     uint      `json:"dish_id"`
		Date       time.Time `json:"date"`
		DailyLimit int       `json:"daily_limit"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误")
		return
	}
	if err := h.svc.SetStock(req.DishID, req.Date, req.DailyLimit); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	response.Success(c, gin.H{"dish_id": req.DishID})
}
