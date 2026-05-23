# 汪托帮点餐系统 PRD V4.0

> 三级会员 + 积分 + 卡券 + 活动，单门店 Go 微服务架构
> 产品部 | 2025年6月

---

## 第一章 技术架构总览

### 1.1 技术选型

| 层级 | 技术栈 | 说明 |
|------|--------|------|
| **后端** | Go 1.24 + Gin/Fiber | 统一语言，全服务 Go |
| **数据库** | PostgreSQL 18.3 | 各服务独立数据库 |
| **缓存** | Redis | 购物车实时同步、库存锁、Session |
| **API 网关** | Go 自研或 Traefik | 路由分发、JWT 鉴权、限流 |
| **服务注册/配置** | 简化方案：环境变量 + YAML 配置文件 | 单门店无需 Nacos；集群部署时可升级为 Consul/etcd |
| **异步消息** | Redis Pub/Sub 或直接 HTTP 回调 | 单门店规模下不需要 Kafka/RabbitMQ |
| **微信小程序** | 原生 / Taro + Vant Weapp | — |
| **后台管理** | React 18 + Ant Design | — |
| **监控** | 基础：Go pprof + 日志；可扩展 Prometheus | 单门店够用 |
| **部署** | Docker Compose（单机）/ K8s（扩展） | — |

### 1.2 服务拆分（优化后）

与 V3 相比：移除独立的**卡券服务**和**消息通知服务**，改为直接调用微信原生 API。

| 服务 | 负责人 | 核心职责 | 独立数据库 |
|------|--------|----------|------------|
| API Gateway | 基础设施组 | 路由、JWT 鉴权、限流、TraceID | 无 |
| 用户服务 | A | 微信登录、三级会员、余额、宠物 | user_db |
| 座位服务 | B | 区域/座位/二维码、状态流转 | seat_db |
| 菜品服务 | C | 分类、菜品、库存、图片 | menu_db |
| 订单服务 | D | 购物车、订单状态机、多人同步 | order_db |
| 支付服务 | E | 微信支付、余额支付、充值、退款 | payment_db |
| 积分服务 | F | 积分发放、兑换、积分商城 | points_db |
| 活动服务 | G | 公告轮播、充值折扣、活动报名 | activity_db |
| 营销定价 | H | 双价、时段价、套餐、满减 | pricing_db |
| 数据统计 | I | 报表、经营分析、导出 | analytics_db |
| 后台聚合 | J | BFF 聚合 + 权限 | 无（聚合层） |
| 小程序前端 | K | 扫码点餐、积分、活动、会员 | 无 |
| 后台管理前端 | L | 运营后台、数据可视化 | 无 |

**服务数量：从 16 个减少到 13 个**（砍掉独立卡券服务、消息通知服务、注册/配置中心）。

---

## 第二章 微信原生能力替代自建服务

### 2.1 卡券 → 微信商家券 API

**结论：不再自建卡券服务，直接调用微信支付"商家券"API。**

微信商家券（Merchant Coupon）原生支持：

| 功能 | 微信商家券 API | 原自建卡券服务 |
|------|---------------|---------------|
| 创建卡券模板 | `POST /v3/marketing/favor/coupon-stocks` | 后台 CRUD + coupon_template 表 |
| 发放卡券 | `POST /v3/marketing/favor/users/{openid}/coupons` | user_coupon 表 + 发放记录 |
| 用户卡券包 | 微信卡包原生展示 | 需自建卡券包 UI |
| 核销 | `POST /v3/marketing/favor/users/{openid}/coupons/{coupon_id}` | coupon_use_record 表 |
| 转赠 | 微信原生转赠（需在模板中开启） | 自建转赠 + 新卡券生成 |
| 过期提醒 | 微信自动推送 | 自建定时任务 + 推送 |
| 查询 | `GET /v3/marketing/favor/users/{openid}/coupons` | 自建查询接口 |

**简化方案**：

1. 在微信支付后台/API 创建卡券模板（满减券、折扣券、兑换券等）
2. 后端仅需一个 `wechat_coupon.go` 工具模块，封装微信商家券 API 调用
3. 充值后调用发放接口 → 卡券自动进入用户微信卡包
4. 下单结算时调用查询接口获取用户可用卡券列表
5. 支付时调用核销接口抵扣金额
6. **不需要** coupon_template / user_coupon / coupon_send_record / coupon_use_record 四张表

**代码量对比**：
- 原方案：独立微服务 + 4 张表 + 10+ API 接口 + 过期定时任务
- 新方案：1 个 Go 包，约 300 行代码

### 2.2 消息推送 → 微信订阅消息

**不再自建通知服务，直接调用微信小程序订阅消息。**

| 通知场景 | 微信订阅消息模板 |
|----------|------------------|
| 积分到账 | `points_received` |
| 卡券到账 | 微信自动推送（商家券自带） |
| 卡券过期 | 微信自动推送（商家券自带） |
| 活动报名确认 | `activity_registered` |
| 活动开始提醒 | `activity_reminder` |
| 订单状态变更 | `order_status_changed` |
| 会员升级 | `member_upgraded` |

简化方案：封装一个 `wechat_notify.go` 工具包，其他服务直接调用。

### 2.3 后厨打印

保留轻量级方案：后端生成打印数据 → 通过 WebSocket/SSE 推送到后厨打印终端。不建独立服务，订单服务内嵌打印模块。

### 2.4 语音播报

依赖后厨管理终端（Windows/Android）的本地能力，后端仅推送文本内容。

---

## 第三章 核心业务服务详细设计

### 3.1 用户服务 (User Service)

**技术栈**：Go + Gin + PostgreSQL

#### 数据模型

| 表名 | 核心字段 |
|------|----------|
| user | id, openid, unionid, nickname, avatar_url, phone, member_level, balance, total_consumption, total_orders, created_at |
| recharge_record | id, user_id, amount, gifted_amount, channel, status, created_at |
| balance_log | id, user_id, type, amount, order_no, remark, created_at |
| consumption_record | id, user_id, order_id, amount, dish_count, created_at |
| pet_profile | id, user_id, name, breed, weight, birthday, created_at |

#### 三级会员

| 等级 | 获取方式 | 权益 |
|------|----------|------|
| 普通客户 | 注册即享 | 原价、可消费 |
| 会员客户 | 累计消费满指定金额/次数 | 会员价、积分×1.5 |
| 充值客户 | 主动充值任意金额 | 会员价、充值折扣、余额支付、积分×2、充值赠券 |

#### API

```
POST   /api/user/wx-login              # 微信 code 换 JWT
GET    /api/user/profile               # 用户信息+等级+余额
GET    /api/user/consumption           # 消费记录
GET    /api/user/consumption/summary   # 消费汇总
POST   /api/user/recharge              # 充值
POST   /api/user/balance/deduct       # 余额扣款（内部）
POST   /api/user/balance/refund       # 余额退款（内部）
GET    /api/user/pets                  # 宠物列表
POST   /api/user/pets                  # 添加宠物
GET    /api/user/internal/:id          # 内部RPC
```

---

### 3.2 座位服务 (Seat Service)

**技术栈**：Go + Gin + PostgreSQL

#### 数据模型

| 表名 | 核心字段 |
|------|----------|
| area | id, name, sort_order |
| seat | id, area_id, name, type, capacity, qrcode_url, status |
| seat_status_log | id, seat_id, old_status, new_status, order_id, changed_at |

#### API

```
GET    /api/seat/areas                 # 区域列表
POST   /api/seat/areas                 # 新增区域
GET    /api/seat/list                  # 座位列表
GET    /api/seat/:id                   # 座位详情
POST   /api/seat/qrcode/batch          # 批量生成二维码
GET    /api/seat/scan?code=xxx         # 扫码解析
GET    /api/seat/internal/:id          # 内部RPC
```

---

### 3.3 菜品服务 (Menu Service)

**技术栈**：Go + Gin + PostgreSQL

#### 数据模型

| 表名 | 核心字段 |
|------|----------|
| category | id, name, parent_id, sort_order, status |
| dish | id, category_id, name, subtitle, description, images, tags, status |
| dish_price | id, dish_id, price_type, price, start_time, end_time |
| dish_stock | id, dish_id, daily_limit, sold_count, date |

#### API

```
GET    /api/menu/categories            # 分类树
GET    /api/menu/dishes                # 菜品列表
GET    /api/menu/dish/:id              # 菜品详情
GET    /api/menu/search                # 搜索
POST   /api/menu/dishes/batch          # 批量查询（内部）
CRUD   /api/menu/admin/*               # 后台管理
```

---

### 3.4 订单服务 (Order Service)

**技术栈**：Go + Gin + PostgreSQL + Redis

#### 数据模型

| 表名 | 核心字段 |
|------|----------|
| cart | id, seat_id, user_id, dish_id, quantity, remark |
| order | id, order_no, seat_id, user_id, status, total_amount, discount_amount, pay_amount |
| order_item | id, order_id, dish_id, dish_name, quantity, unit_price |
| order_status_log | id, order_id, from_status, to_status, operator |

#### API

```
POST   /api/order/cart/add             # 加入购物车
GET    /api/order/cart/list            # 购物车列表（按座位）
PUT    /api/order/cart/update          # 修改数量/备注
DELETE /api/order/cart/remove          # 移除商品
POST   /api/order/create               # 提交订单
POST   /api/order/:id/pay              # 发起支付
GET    /api/order/:id/status           # 订单状态
GET    /api/order/list                 # 我的订单
POST   /api/order/internal/notify      # 支付回调（内部）
GET    /api/order/admin/list           # 后台订单
PUT    /api/order/admin/status         # 改状态
POST   /api/order/admin/refund         # 退单
```

**下单主流程**：
1. 校验座位状态 → 2. 校验库存 → 3. 生成订单 → 4. 调用定价服务算价格 → 5. 锁定库存 15 分钟 → 6. 调支付服务 → 7. 返回支付参数

---

### 3.5 支付服务 (Payment Service)

**技术栈**：Go + Gin + PostgreSQL + Redis

#### 数据模型

| 表名 | 核心字段 |
|------|----------|
| payment_order | id, order_no, out_trade_no, user_id, amount, channel, status, wx_prepay_id |
| payment_record | id, payment_order_id, channel, amount, transaction_id, paid_at |
| refund_record | id, payment_order_id, refund_no, amount, reason, status |
| recharge_order | id, user_id, amount, gifted_amount, discount_rate, final_amount, status |

#### API

```
POST   /api/pay/create                 # 创建支付单
POST   /api/pay/wx/prepay              # 微信统一下单
POST   /api/pay/balance                # 余额支付
POST   /api/pay/recharge               # 充值
POST   /api/pay/callback/wx            # 微信回调
POST   /api/pay/refund                 # 退款
GET    /api/pay/query/:outTradeNo      # 查询
```

---

### 3.6 卡券（微信商家券集成，非独立服务）

**不建独立微服务**。在多处嵌入微信商家券调用：

- **充值赠券**：用户服务充值时调用 `wechat.Coupon.Send(openid, stock_id)`
- **卡券列表**：小程序直接调微信 API 获取用户卡包，或后端代理
- **核销**：支付服务下单时调用 `wechat.Coupon.Use(openid, coupon_id, order_no)`
- **转赠**：微信原生支持，无需后端介入

后端仅需一个 `internal/wechat/coupon.go` 文件：

```go
package wechat

// 创建卡券批次
func CreateStock(req CreateStockRequest) (*CreateStockResponse, error)

// 发放卡券
func SendCoupon(openid, stockID string) (*SendCouponResponse, error)

// 查询用户卡券
func ListUserCoupons(openid string) ([]Coupon, error)

// 核销卡券
func UseCoupon(openid, couponID, orderNo string) error
```

---

### 3.7 积分服务 (Points Service)

**技术栈**：Go + Gin + PostgreSQL + Redis

#### 数据模型

| 表名 | 核心字段 |
|------|----------|
| points_rule | id, name, type, config_json, status |
| user_points | id, user_id, total_points, used_points, frozen_points |
| points_log | id, user_id, type, points, source_id, remark |
| exchange_goods | id, name, image, points_price, stock, type, status |
| exchange_order | id, user_id, goods_id, points_cost, status |

#### API

```
GET    /api/points/account             # 积分余额
GET    /api/points/logs                # 积分流水
GET    /api/points/goods               # 兑换商品列表
POST   /api/points/exchange            # 积分兑换
POST   /api/points/internal/grant      # 发放积分（内部）
CRUD   /api/points/admin/*             # 后台管理
```

**积分倍率**：普通×1、会员×1.5、充值×2

---

### 3.8 活动服务 (Activity Service)

**技术栈**：Go + Gin + PostgreSQL

#### 数据模型

| 表名 | 核心字段 |
|------|----------|
| announcement | id, title, content, type, image, link_type, link_target, sort_order, start_time, end_time |
| activity | id, title, description, image, max_participants, current_participants, event_time, location, status |
| activity_registration | id, activity_id, user_id, name, phone, remark, status |

#### API

```
GET    /api/activity/announcements              # 生效公告
GET    /api/activity/recharge-discount           # 充值折扣配置
GET    /api/activity/list                        # 活动列表
POST   /api/activity/:id/register                # 报名
GET    /api/activity/my-registrations            # 我的报名
PUT    /api/activity/:id/cancel                  # 取消报名
GET    /api/activity/internal/discount            # 内部查折扣
CRUD   /api/activity/admin/*                     # 后台管理
```

---

### 3.9 营销定价服务 (Pricing Service)

**技术栈**：Go + Gin + PostgreSQL

#### 数据模型

| 表名 | 核心字段 |
|------|----------|
| price_rule | id, dish_id, rule_type, price, start_time, end_time, status |
| promotion | id, name, type, config_json, start_time, end_time, status |
| combo | id, name, price, dish_list, status |

#### API

```
POST   /api/pricing/calculate          # 计算订单价格
GET    /api/pricing/dish/:id           # 菜品当前价
GET    /api/pricing/promotions         # 营销活动
CRUD   /api/pricing/admin/*            # 后台管理
```

---

### 3.10 数据统计服务 (Analytics Service)

**技术栈**：Go + Gin + PostgreSQL（单门店无需 ClickHouse）

#### API

```
GET    /api/analytics/dashboard        # 仪表盘
GET    /api/analytics/revenue          # 营收
GET    /api/analytics/dishes           # 菜品分析
GET    /api/analytics/members          # 会员分布
GET    /api/analytics/points           # 积分统计
GET    /api/analytics/coupons          # 卡券统计
GET    /api/analytics/activities       # 活动统计
POST   /api/analytics/export           # Excel 导出
```

---

### 3.11 后台聚合服务 (Admin BFF)

**技术栈**：Go + Gin（纯聚合，无数据库）

聚合所有下游服务 API，为后台前端提供统一接口层。权限控制、操作审计。

---

### 3.12 消息推送（内嵌，非独立服务）

各服务在需要推送时直接调用：

```go
wechat.SendSubscribeMsg(openid, templateID, data, page)
```

| 场景 | 模板 |
|------|------|
| 积分到账 | `您已获得{{amount}}积分` |
| 订单状态变更 | `订单{{order_no}} 已{{status}}` |
| 会员升级 | `恭喜您已升级为{{level}}` |
| 活动报名确认 | `您已成功报名{{title}}` |
| 活动开始提醒 | `{{title}}即将开始` |

---

## 第四章 前端应用需求

### 4.1 微信小程序

**技术栈**：原生或 Taro + Vant Weapp

| 模块 | 核心页面 |
|------|----------|
| 首页 | 公告轮播、充值折扣入口、活动推荐 |
| 扫码点餐 | 扫码页 → 菜单页 → 购物车 → 订单确认（含微信原生卡券选择）→ 支付页 |
| 订单 | 订单列表、状态时间轴、呼叫服务员 |
| 会员中心 | 等级展示+升级进度条、积分余额、消费曲线、充值中心 |
| 积分商城 | 兑换商品列表、兑换记录 |
| 卡券包 | **直接使用微信卡包原生页面** |
| 活动中心 | 活动列表、详情、报名 |

### 4.2 后台管理前端

**技术栈**：React 18 + TypeScript + Ant Design 5 + ECharts

| 模块 | 核心功能 |
|------|----------|
| 座位管理 | 可视化布局编辑器、二维码下载 |
| 菜品管理 | 分类树、菜品 CRUD、库存、价格 |
| 订单中心 | 实时看板、状态流转、退单 |
| 会员管理 | 三级会员列表、消费统计、积分调整 |
| 积分中心 | 规则配置、流水查询、兑换商品 |
| 活动中心 | 公告管理、活动发布、报名名单 |
| 价格策略 | 双价、时段价、套餐、满减 |
| 数据统计 | 仪表盘、报表导出 |
| 系统设置 | 充值档位、消息模板、账号权限 |

---

## 第五章 服务间通信

### 5.1 同步通信：RESTful API

统一规范：

```json
// 成功
{ "code": 200, "message": "ok", "data": {...} }

// 错误
{ "code": 40001, "message": "库存不足", "data": null }

// 分页
{ "code": 200, "data": { "total": 100, "page": 1, "pageSize": 20, "list": [...] } }
```

- 认证头：`Authorization: Bearer {JWT}`
- 链路追踪：`X-Trace-ID: {uuid}`
- 内部接口以 `/api/{service}/internal/` 开头，网关禁止外部路由

### 5.2 异步通信（简化方案）

单门店规模不引入消息队列。异步场景用**直接 HTTP 回调**：

```go
// 订单支付成功后
orderService.OnPaid(order) → http.Post(userServiceURL + "/internal/consume", ...)
orderService.OnPaid(order) → http.Post(pointsServiceURL + "/internal/grant", ...)
```

如需解耦，后续可平滑升级到 Redis Pub/Sub。

---

## 第六章 非功能需求

### 6.1 性能

- 小程序首屏 < 2s（3G）
- 订单提交 P99 < 800ms
- 微服务内部调用 P99 < 200ms
- 单门店峰值：100 单/分钟

### 6.2 安全

- 全链路 HTTPS
- 支付密钥环境变量存储，不入代码仓
- 后台关键操作二次验证
- JWT 过期时间 2h，refresh 7 天

---

## 第七章 项目结构建议

```
wtb-ordering/
├── gateway/            # API 网关
├── services/
│   ├── user/           # 用户服务
│   ├── seat/           # 座位服务
│   ├── menu/           # 菜品服务
│   ├── order/          # 订单服务
│   ├── payment/        # 支付服务
│   ├── points/         # 积分服务
│   ├── activity/       # 活动服务
│   ├── pricing/        # 营销定价
│   ├── analytics/      # 数据统计
│   └── admin/          # 后台聚合
├── internal/
│   └── wechat/         # 微信工具包（登录/支付/商家券/订阅消息）
├── pkg/
│   ├── jwt/            # JWT 工具
│   ├── response/       # 统一响应
│   └── httpclient/     # HTTP 客户端
├── miniprogram/        # 微信小程序
├── admin-web/          # 后台管理前端
├── docker-compose.yml
└── Makefile
```

---

## 附录：V3 → V4 变更摘要

| 变更项 | V3 | V4 | 理由 |
|--------|----|----|------|
| 技术栈 | Spring Boot / Go | 纯 Go | 统一语言，降低团队成本 |
| 卡券服务 | 独立微服务 + 4 张表 | 微信商家券 API | 微信原生支持全部功能 |
| 消息通知 | 独立微服务 + Redis Stream | 内嵌微信订阅消息调用 | 单门店无需独立服务 |
| 服务注册 | Nacos (Java 生态) | 环境变量 + YAML | Go 无原生 Nacos 支持 |
| 消息队列 | RabbitMQ / Redis Stream | 直接 HTTP 回调 | 单门店规模够用 |
| 统计存储 | ClickHouse/ES | PostgreSQL | 单门店数据量小 |
| 服务数量 | 16 | 13 | 砍掉 3 个非必要服务 |
