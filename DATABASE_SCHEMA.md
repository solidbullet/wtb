# 汪托帮点餐系统 — 数据库 Schema 汇总

> PostgreSQL 18.3 | 9 个独立数据库 | 每个服务 1 个库
> GORM AutoMigrate 自动建表，以下为参考 DDL
> 金额单位：分（INTEGER）

---

## 数据库分布总览

```
wtb_user       → 用户、宠物、充值记录、余额日志、消费记录
wtb_seat       → 区域、座位、座位状态日志
wtb_menu       → 分类、菜品、菜品价格、菜品库存
wtb_order      → 订单、订单明细、订单状态日志
wtb_payment    → 支付单、支付记录、退款记录、充值订单
wtb_points     → 积分规则、用户积分、积分日志、兑换商品、兑换订单
wtb_activity   → 公告、活动、活动报名
wtb_pricing    → 价格规则、营销活动、套餐
wtb_analytics  →（按需创建汇总表，多数查询跨库）
```

---

## 1. wtb_user — 用户服务

### 1.1 users

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| openid | VARCHAR(64) | UNIQUE, NOT NULL | 微信 OpenID |
| unionid | VARCHAR(64) | DEFAULT '' | 微信 UnionID |
| nickname | VARCHAR(100) | DEFAULT '' | 昵称 |
| avatar_url | VARCHAR(500) | DEFAULT '' | 头像 URL |
| phone | VARCHAR(20) | DEFAULT '' | 手机号 |
| member_level | SMALLINT | DEFAULT 0 | 0普通 1会员 2充值 |
| balance | INTEGER | DEFAULT 0 | 余额（分） |
| total_consumption | INTEGER | DEFAULT 0 | 累计消费（分） |
| total_orders | INTEGER | DEFAULT 0 | 累计订单数 |
| created_at | TIMESTAMP | DEFAULT NOW() | — |
| updated_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_users_openid (openid)`

### 1.2 recharge_records

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| user_id | BIGINT | NOT NULL | FK → users.id |
| amount | INTEGER | NOT NULL | 充值金额（分） |
| gifted_amount | INTEGER | DEFAULT 0 | 赠送金额（分） |
| channel | VARCHAR(20) | DEFAULT 'wxpay' | — |
| status | VARCHAR(20) | DEFAULT 'pending' | pending/paid/failed |
| created_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_recharge_user (user_id)`

### 1.3 balance_logs

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| user_id | BIGINT | NOT NULL | FK → users.id |
| type | VARCHAR(20) | NOT NULL | recharge/deduct/refund |
| amount | INTEGER | NOT NULL | 变动金额（分） |
| order_no | VARCHAR(64) | DEFAULT '' | 关联订单号 |
| remark | VARCHAR(255) | DEFAULT '' | — |
| created_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_balance_log_user (user_id)`

### 1.4 consumption_records

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| user_id | BIGINT | NOT NULL | FK → users.id |
| order_id | BIGINT | NOT NULL | 关联订单 ID |
| amount | INTEGER | NOT NULL | 消费金额（分） |
| dish_count | INTEGER | DEFAULT 0 | 菜品数量 |
| created_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_consumption_user (user_id)`

### 1.5 pet_profiles

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| user_id | BIGINT | NOT NULL | FK → users.id |
| name | VARCHAR(50) | NOT NULL | 宠物名 |
| breed | VARCHAR(50) | DEFAULT '' | 品种 |
| weight | NUMERIC(5,2) | DEFAULT 0 | 体重（kg） |
| birthday | DATE | — | 生日 |
| created_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_pet_user (user_id)`

---

## 2. wtb_seat — 座位服务

### 2.1 areas

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| name | VARCHAR(50) | NOT NULL | 区域名 |
| sort_order | INTEGER | DEFAULT 0 | 排序 |
| created_at | TIMESTAMP | DEFAULT NOW() | — |

### 2.2 seats

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| area_id | BIGINT | NOT NULL | FK → areas.id |
| name | VARCHAR(50) | NOT NULL | 座位名 |
| type | VARCHAR(20) | DEFAULT 'normal' | normal/booth/outdoor |
| capacity | INTEGER | DEFAULT 4 | 容纳人数 |
| qrcode_url | VARCHAR(500) | DEFAULT '' | 二维码图片 URL |
| status | VARCHAR(20) | DEFAULT 'available' | available/occupied/reserved/cleaning |
| created_at | TIMESTAMP | DEFAULT NOW() | — |
| updated_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_seats_area (area_id)`

### 2.3 seat_status_logs

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| seat_id | BIGINT | NOT NULL | FK → seats.id |
| old_status | VARCHAR(20) | — | — |
| new_status | VARCHAR(20) | — | — |
| order_id | BIGINT | — | 关联订单 ID |
| changed_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_seat_log_seat (seat_id)`

---

## 3. wtb_menu — 菜品服务

### 3.1 categories

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| name | VARCHAR(50) | NOT NULL | 分类名 |
| parent_id | BIGINT | DEFAULT 0 | 父分类 ID（0=顶级） |
| sort_order | INTEGER | DEFAULT 0 | — |
| status | SMALLINT | DEFAULT 1 | 1启用 0禁用 |
| created_at | TIMESTAMP | DEFAULT NOW() | — |

### 3.2 dishes

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| category_id | BIGINT | NOT NULL | FK → categories.id |
| name | VARCHAR(100) | NOT NULL | 菜品名 |
| subtitle | VARCHAR(200) | DEFAULT '' | 副标题 |
| description | TEXT | — | 描述 |
| images | TEXT | — | JSON 数组 |
| tags | VARCHAR(200) | DEFAULT '' | 逗号分隔 |
| status | SMALLINT | DEFAULT 1 | 1上架 0下架 |
| search_vector | TSVECTOR | GENERATED | 全文搜索向量 |
| created_at | TIMESTAMP | DEFAULT NOW() | — |
| updated_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_dishes_category (category_id)`, `idx_dishes_search USING GIN (search_vector)`

### 3.3 dish_prices

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| dish_id | BIGINT | NOT NULL | FK → dishes.id |
| price_type | VARCHAR(20) | NOT NULL | normal/member/time_slot |
| price | INTEGER | NOT NULL | 价格（分） |
| start_time | TIMESTAMP | — | 时段价开始 |
| end_time | TIMESTAMP | — | 时段价结束 |

索引：`idx_dish_prices_dish (dish_id)`

### 3.4 dish_stocks

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| dish_id | BIGINT | NOT NULL | FK → dishes.id |
| daily_limit | INTEGER | DEFAULT -1 | -1=无限 |
| sold_count | INTEGER | DEFAULT 0 | 已售数量 |
| date | DATE | NOT NULL | 日期 |

唯一约束：`UNIQUE (dish_id, date)`

---

## 4. wtb_order — 订单服务

### 4.1 orders

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| order_no | VARCHAR(32) | UNIQUE, NOT NULL | 订单号 WTB+时间戳 |
| seat_id | BIGINT | NOT NULL | FK → seats.id（跨库） |
| user_id | BIGINT | NOT NULL | FK → users.id（跨库） |
| status | VARCHAR(20) | DEFAULT 'pending' | pending/confirmed/cooking/served/paid/completed/cancelled/refunded |
| total_amount | INTEGER | NOT NULL | 原价总额（分） |
| discount_amount | INTEGER | DEFAULT 0 | 优惠金额（分） |
| pay_amount | INTEGER | NOT NULL | 实付金额（分） |
| remark | VARCHAR(500) | DEFAULT '' | 备注 |
| created_at | TIMESTAMP | DEFAULT NOW() | — |
| updated_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_orders_user (user_id)`, `idx_orders_seat (seat_id)`, `idx_orders_no (order_no)`

### 4.2 order_items

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| order_id | BIGINT | NOT NULL | FK → orders.id |
| dish_id | BIGINT | NOT NULL | FK → dishes.id（跨库） |
| dish_name | VARCHAR(100) | NOT NULL | 快照：下单时菜品名 |
| quantity | INTEGER | NOT NULL | 数量 |
| unit_price | INTEGER | NOT NULL | 快照：下单时单价（分） |

索引：`idx_order_items_order (order_id)`

### 4.3 order_status_logs

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| order_id | BIGINT | NOT NULL | FK → orders.id |
| from_status | VARCHAR(20) | — | — |
| to_status | VARCHAR(20) | — | — |
| operator | VARCHAR(50) | DEFAULT '' | 操作人 |
| created_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_order_status_log_order (order_id)`

---

## 5. wtb_payment — 支付服务

### 5.1 payment_orders

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| order_no | VARCHAR(32) | NOT NULL | 订单号 |
| out_trade_no | VARCHAR(32) | UNIQUE, NOT NULL | 商户订单号 |
| user_id | BIGINT | NOT NULL | FK → users.id（跨库） |
| amount | INTEGER | NOT NULL | 支付金额（分） |
| channel | VARCHAR(20) | NOT NULL | wxpay/balance |
| status | VARCHAR(20) | DEFAULT 'pending' | pending/paid/closed/refunded |
| wx_prepay_id | VARCHAR(64) | DEFAULT '' | 微信预支付 ID |
| created_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_payment_orders_no (order_no)`

### 5.2 payment_records

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| payment_order_id | BIGINT | NOT NULL | FK → payment_orders.id |
| channel | VARCHAR(20) | NOT NULL | — |
| amount | INTEGER | NOT NULL | — |
| transaction_id | VARCHAR(64) | DEFAULT '' | 微信交易单号 |
| paid_at | TIMESTAMP | — | 支付时间 |

索引：`idx_payment_records_order (payment_order_id)`

### 5.3 refund_records

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| payment_order_id | BIGINT | NOT NULL | FK → payment_orders.id |
| refund_no | VARCHAR(32) | UNIQUE, NOT NULL | 退款单号 |
| amount | INTEGER | NOT NULL | 退款金额（分） |
| reason | VARCHAR(200) | DEFAULT '' | — |
| status | VARCHAR(20) | DEFAULT 'pending' | pending/success/failed |
| created_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_refund_records_order (payment_order_id)`

### 5.4 recharge_orders

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| user_id | BIGINT | NOT NULL | FK → users.id（跨库） |
| amount | INTEGER | NOT NULL | 充值金额（分） |
| gifted_amount | INTEGER | DEFAULT 0 | 赠送金额（分） |
| discount_rate | NUMERIC(3,2) | DEFAULT 1.00 | 折扣率 |
| final_amount | INTEGER | NOT NULL | 到账金额（分） |
| status | VARCHAR(20) | DEFAULT 'pending' | — |
| created_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_recharge_orders_user (user_id)`

---

## 6. wtb_points — 积分服务

### 6.1 points_rules

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| name | VARCHAR(50) | NOT NULL | 规则名 |
| type | VARCHAR(30) | NOT NULL | consumption/recharge/sign_in |
| config_json | TEXT | NOT NULL | 规则配置 JSON |
| status | SMALLINT | DEFAULT 1 | — |

### 6.2 user_points

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| user_id | BIGINT | UNIQUE, NOT NULL | FK → users.id（跨库） |
| total_points | INTEGER | DEFAULT 0 | 累计获得 |
| used_points | INTEGER | DEFAULT 0 | 已使用 |
| frozen_points | INTEGER | DEFAULT 0 | 冻结中 |
| updated_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_user_points_user (user_id)`

### 6.3 points_logs

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| user_id | BIGINT | NOT NULL | FK → users.id（跨库） |
| type | VARCHAR(20) | NOT NULL | gain/exchange/expire/adjust |
| points | INTEGER | NOT NULL | 变动积分（正=获得，负=消耗） |
| source_id | VARCHAR(64) | DEFAULT '' | 关联订单号 |
| remark | VARCHAR(200) | DEFAULT '' | — |
| created_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_points_logs_user_time (user_id, created_at)`

### 6.4 exchange_goods

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| name | VARCHAR(100) | NOT NULL | 商品名 |
| image | VARCHAR(500) | DEFAULT '' | 图片 URL |
| points_price | INTEGER | NOT NULL | 兑换所需积分 |
| stock | INTEGER | DEFAULT 0 | 库存 |
| type | VARCHAR(20) | DEFAULT 'physical' | physical/coupon |
| status | SMALLINT | DEFAULT 1 | — |

### 6.5 exchange_orders

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| user_id | BIGINT | NOT NULL | FK → users.id（跨库） |
| goods_id | BIGINT | NOT NULL | FK → exchange_goods.id |
| points_cost | INTEGER | NOT NULL | 消耗积分 |
| status | VARCHAR(20) | DEFAULT 'pending' | pending/delivered/cancelled |
| created_at | TIMESTAMP | DEFAULT NOW() | — |

索引：`idx_exchange_orders_user (user_id)`

---

## 7. wtb_activity — 活动服务

### 7.1 announcements

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| title | VARCHAR(100) | NOT NULL | — |
| content | TEXT | — | — |
| type | VARCHAR(20) | DEFAULT 'text' | text/image |
| image | VARCHAR(500) | DEFAULT '' | — |
| link_type | VARCHAR(20) | DEFAULT '' | recharge/activity/external |
| link_target | VARCHAR(500) | DEFAULT '' | — |
| sort_order | INTEGER | DEFAULT 0 | — |
| start_time | TIMESTAMP | NOT NULL | — |
| end_time | TIMESTAMP | NOT NULL | — |
| status | SMALLINT | DEFAULT 1 | — |

索引：`idx_announcements_time (start_time, end_time)`

### 7.2 activities

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| title | VARCHAR(100) | NOT NULL | — |
| description | TEXT | — | — |
| image | VARCHAR(500) | DEFAULT '' | — |
| max_participants | INTEGER | DEFAULT -1 | -1=不限 |
| current_participants | INTEGER | DEFAULT 0 | — |
| event_time | TIMESTAMP | — | 活动时间 |
| location | VARCHAR(200) | DEFAULT '' | — |
| status | VARCHAR(20) | DEFAULT 'draft' | draft/published/cancelled/ended |
| created_at | TIMESTAMP | DEFAULT NOW() | — |

### 7.3 activity_registrations

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| activity_id | BIGINT | NOT NULL | FK → activities.id |
| user_id | BIGINT | NOT NULL | FK → users.id（跨库） |
| name | VARCHAR(50) | DEFAULT '' | — |
| phone | VARCHAR(20) | DEFAULT '' | — |
| remark | VARCHAR(200) | DEFAULT '' | — |
| status | VARCHAR(20) | DEFAULT 'registered' | registered/cancelled |
| created_at | TIMESTAMP | DEFAULT NOW() | — |

唯一约束：`UNIQUE (user_id, activity_id)` — 同一用户不可重复报名

---

## 8. wtb_pricing — 营销定价服务

### 8.1 price_rules

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| dish_id | BIGINT | NOT NULL | FK → dishes.id（跨库） |
| rule_type | VARCHAR(20) | NOT NULL | normal/member/time_slot |
| price | INTEGER | NOT NULL | 价格（分） |
| start_time | TIMESTAMP | — | 时段价生效时间 |
| end_time | TIMESTAMP | — | 时段价结束时间 |
| status | SMALLINT | DEFAULT 1 | — |

索引：`idx_price_rules_dish (dish_id)`

### 8.2 promotions

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| name | VARCHAR(100) | NOT NULL | 活动名 |
| type | VARCHAR(30) | NOT NULL | full_reduction/discount |
| config_json | TEXT | NOT NULL | 规则 JSON |
| start_time | TIMESTAMP | NOT NULL | — |
| end_time | TIMESTAMP | NOT NULL | — |
| status | SMALLINT | DEFAULT 1 | — |

索引：`idx_promotions_time (start_time, end_time)`

**config_json 示例**：

```json
// full_reduction: { "threshold": 5000, "reduce": 1000 }
// discount:      { "rate": 0.85, "max_reduce": 2000 }
```

### 8.3 combos

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGSERIAL | PK | — |
| name | VARCHAR(100) | NOT NULL | 套餐名 |
| price | INTEGER | NOT NULL | 套餐价格（分） |
| dish_list | TEXT | NOT NULL | JSON: `[{dish_id, quantity}]` |
| status | SMALLINT | DEFAULT 1 | — |

---

## 9. 跨库关系图

```
users (wtb_user)
  │
  ├──→ orders.user_id (wtb_order)
  │      └──→ order_items.order_id
  │
  ├──→ payment_orders.user_id (wtb_payment)
  │      ├──→ payment_records.payment_order_id
  │      └──→ refund_records.payment_order_id
  │
  ├──→ recharge_orders.user_id (wtb_payment)
  ├──→ recharge_records.user_id (wtb_user)
  ├──→ balance_logs.user_id (wtb_user)
  ├──→ consumption_records.user_id (wtb_user)
  ├──→ pet_profiles.user_id (wtb_user)
  │
  ├──→ user_points.user_id (wtb_points)
  │      └──→ points_logs.user_id
  │
  ├──→ exchange_orders.user_id (wtb_points)
  │
  └──→ activity_registrations.user_id (wtb_activity)

seats (wtb_seat)
  │
  ├──→ orders.seat_id (wtb_order)
  └──→ seat_status_logs.seat_id

dishes (wtb_menu)
  │
  ├──→ dish_prices.dish_id (wtb_menu)
  ├──→ dish_stocks.dish_id (wtb_menu)
  ├──→ order_items.dish_id (wtb_order)
  └──→ price_rules.dish_id (wtb_pricing)

activities (wtb_activity)
  │
  └──→ activity_registrations.activity_id
```

---

> 跨库查询不依赖外键约束，通过应用层代码关联。字段类型、命名、金额单位全库统一。
