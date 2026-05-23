# internal/wechat

封装微信登录、支付、商家券、订阅消息 API。

## 模块

- `config.go` — 配置结构
- `wechat.go` — 客户端初始化
- `login.go` — 小程序 code2session
- `pay.go` — 微信支付统一下单 + 回调验证
- `coupon.go` — 微信商家券（创建/发放/查询/核销）
- `notify.go` — 微信小程序订阅消息

## 开发阶段

所有支付/卡券/消息接口当前为 Mock 实现，配置空的商户号时返回模拟数据。
