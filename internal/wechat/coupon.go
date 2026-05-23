package wechat

import "fmt"

// Coupon 商家券响应
type Coupon struct {
	CouponID  string `json:"coupon_id"`
	StockID   string `json:"stock_id"`
	Status    string `json:"status"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// CreateStockRequest 创建卡券批次请求
type CreateStockRequest struct {
	StockName   string `json:"stock_name"`
	TotalCount  int    `json:"total_count"`
	BeginTime   string `json:"begin_time"`
	EndTime     string `json:"end_time"`
	MaxAmount   int    `json:"max_amount"`   // 面额（分）
	MinAmount   int    `json:"min_amount"`   // 门槛（分）
}

// CreateStockResponse 创建卡券批次响应
type CreateStockResponse struct {
	StockID string `json:"stock_id"`
}

// SendCouponResponse 发放卡券响应
type SendCouponResponse struct {
	CouponID string `json:"coupon_id"`
}

// CreateStock 创建卡券批次（Mock）
func (c *Client) CreateStock(req CreateStockRequest) (*CreateStockResponse, error) {
	if c.config.MchID == "" {
		return nil, fmt.Errorf("商户号未配置")
	}
	return &CreateStockResponse{StockID: "mock_stock_" + req.StockName}, nil
}

// SendCoupon 发放卡券（Mock）
func (c *Client) SendCoupon(openid, stockID string) (*SendCouponResponse, error) {
	if c.config.MchID == "" {
		return nil, fmt.Errorf("商户号未配置")
	}
	return &SendCouponResponse{CouponID: "mock_coupon_" + openid}, nil
}

// ListUserCoupons 查询用户卡券（Mock）
func (c *Client) ListUserCoupons(openid string) ([]Coupon, error) {
	if c.config.MchID == "" {
		return nil, fmt.Errorf("商户号未配置")
	}
	return []Coupon{}, nil
}

// UseCoupon 核销卡券（Mock）
func (c *Client) UseCoupon(openid, couponID, orderNo string) error {
	if c.config.MchID == "" {
		return fmt.Errorf("商户号未配置")
	}
	return nil
}
