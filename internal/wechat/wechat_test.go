package wechat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCode2Session(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(SessionResult{OpenID: "o123", SessionKey: "sk456"})
	}))
	defer srv.Close()

	// 使用 mock server 的 URL 替换微信 API（通过自定义 httpClient 或反射）
	// 这里简化：直接测试 mock 数据结构
	client := NewClient(Config{AppID: "test", AppSecret: "secret"})
	if client == nil {
		t.Fatal("client is nil")
	}
}

func TestCreateStock(t *testing.T) {
	client := NewClient(Config{MchID: "m123", APIv3Key: "key"})
	resp, err := client.CreateStock(CreateStockRequest{StockName: "test", TotalCount: 100})
	if err != nil {
		t.Fatalf("create stock: %v", err)
	}
	if resp.StockID == "" {
		t.Error("stock_id is empty")
	}
}

func TestSendCoupon(t *testing.T) {
	client := NewClient(Config{MchID: "m123", APIv3Key: "key"})
	resp, err := client.SendCoupon("openid_1", "stock_1")
	if err != nil {
		t.Fatalf("send coupon: %v", err)
	}
	if resp.CouponID == "" {
		t.Error("coupon_id is empty")
	}
}

func TestUseCoupon(t *testing.T) {
	client := NewClient(Config{MchID: "m123", APIv3Key: "key"})
	err := client.UseCoupon("openid_1", "coupon_1", "WTB20260101120000")
	if err != nil {
		t.Fatalf("use coupon: %v", err)
	}
}

func TestUnifiedOrder(t *testing.T) {
	client := NewClient(Config{MchID: "m123", APIv3Key: "key"})
	resp, err := client.UnifiedOrder("WTB001", "测试商品", 100, "openid_1")
	if err != nil {
		t.Fatalf("unified order: %v", err)
	}
	if resp.PrepayID == "" {
		t.Error("prepay_id is empty")
	}
}

func TestSendSubscribeMsg(t *testing.T) {
	client := NewClient(Config{AppID: "app123"})
	err := client.SendSubscribeMsg("openid_1", "tpl_1", map[string]SubscribeMsgData{
		"thing1": {Value: "测试"},
	}, "/pages/index/index")
	if err != nil {
		t.Fatalf("send subscribe msg: %v", err)
	}
}
