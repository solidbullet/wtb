package wechat

import (
	"fmt"
)

// PrepayResult 微信统一下单结果
type PrepayResult struct {
	PrepayID string `json:"prepay_id"`
}

// UnifiedOrder 统一下单（开发阶段 Mock）
func (c *Client) UnifiedOrder(outTradeNo, description string, amount int, openID string) (*PrepayResult, error) {
	if c.config.MchID == "" {
		return nil, fmt.Errorf("商户号未配置")
	}
	// TODO: 接入微信支付 APIv3
	return &PrepayResult{PrepayID: "mock_prepay_id_" + outTradeNo}, nil
}

// VerifyCallback 验证微信支付回调签名（开发阶段 Mock）
func (c *Client) VerifyCallback(body []byte, signature, timestamp, nonce, serial string) (bool, error) {
	if c.config.APIv3Key == "" {
		return true, nil // Mock 环境跳过验证
	}
	// TODO: 实现 APIv3 签名验证
	return true, nil
}
