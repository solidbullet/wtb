# 汪托帮点餐系统 — API 接口文档 V1.0

> 对应 DEVELOPMENT_PLAN.md V4.0 | PostgreSQL 18.3
> 所有金额单位：分（整数）
> 统一响应格式见 0.1 节

---

## 0. 通用约定

### 0.1 统一响应格式

```json
// 成功
{ "code": 200, "message": "ok", "data": {...} }

// 业务错误
{ "code": 40001, "message": "库存不足", "data": null }

// 分页
{
  "code": 200,
  "message": "ok",
  "data": {
    "total": 100,
    "page": 1,
    "pageSize": 20,
    "list": [...]
  }
}
```

### 0.2 认证

```
Authorization: Bearer {JWT}
```

- 公开接口：`POST /api/user/wx-login`、`POST /api/pay/callback/wx`
- 内部接口：`/api/{service}/internal/*` 不走网关，服务间直连

### 0.3 HTTP 状态码规范

| HTTP Status | 含义 |
|-------------|------|
| 200 | 全部通过（含业务错误，看 `code` 字段） |
| 401 | 未认证 / token 过期 |
| 403 | 无权限 |
| 500 | 服务内部错误 |

### 0.4 业务错误码

| code | 含义 |
|------|------|
| 200 | 成功 |
| 40001 | 参数错误 |
| 40002 | 库存不足 |
| 40003 | 余额不足 |
| 40004 | 座位已被占用 |
| 40005 | 订单状态不允许此操作 |
| 40006 | 活动名额已满 |
| 40007 | 积分不足 |
| 40008 | 已报名/不可重复操作 |
| 40009 | 充值失败 |
| 40010 | 支付失败 |
| 40101 | 未登录 |
| 40102 | Token 过期 |
| 40301 | 非管理员 |
| 50001 | 内部服务调用失败 |

---

## 1. 用户服务 (user-service) — 端口 8081

### 1.1 微信登录

```
POST /api/user/wx-login
公开接口，无需认证
```

**请求体**：
```json
{
  "code": "wx_auth_code_from_miniprogram"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "token": "eyJhbGciOi...",
    "user": {
      "id": 1,
      "openid": "oxxxxxxxxxxxxxx",
      "nickname": "微信用户",
      "avatar_url": "https://...",
      "phone": "",
      "member_level": 0,
      "balance": 0
    }
  }
}
```

**错误**：
| code | message | 场景 |
|------|---------|------|
| 40001 | code 无效 | 微信返回 errcode != 0 |
| 50001 | 微信服务异常 | 网络超时 |

---

### 1.2 用户信息

```
GET /api/user/profile
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "openid": "oxxxxxxxxxxxxxx",
    "nickname": "微信用户",
    "avatar_url": "https://...",
    "phone": "13800138000",
    "member_level": 1,
    "balance": 50000,
    "total_consumption": 120000,
    "total_orders": 15,
    "created_at": "2026-01-15T10:30:00Z"
  }
}
```

---

### 1.3 消费记录

```
GET /api/user/consumption?page=1&pageSize=20
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "total": 15,
    "page": 1,
    "pageSize": 20,
    "list": [
      {
        "id": 1,
        "order_id": 1001,
        "amount": 8800,
        "dish_count": 3,
        "created_at": "2026-06-01T12:30:00Z"
      }
    ]
  }
}
```

---

### 1.4 消费汇总

```
GET /api/user/consumption/summary
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "total_amount": 120000,
    "total_orders": 15,
    "avg_amount": 8000,
    "this_month_amount": 25000,
    "this_month_orders": 3
  }
}
```

---

### 1.5 充值

```
POST /api/user/recharge
需要认证
```

**请求体**：
```json
{
  "amount": 10000,
  "channel": "wxpay"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "recharge_order_id": 2001,
    "amount": 10000,
    "gifted_amount": 2000,
    "final_amount": 12000,
    "wx_pay_params": {
      "prepay_id": "wx...",
      "nonce_str": "...",
      "sign": "...",
      "timestamp": "1717200000"
    }
  }
}
```

**错误**：
| code | message |
|------|---------|
| 40001 | 充值金额无效 |
| 50001 | 创建支付单失败 |

---

### 1.6 余额扣款（内部接口）

```
POST /api/user/balance/deduct
内部接口，服务间调用，不走网关
```

**请求体**：
```json
{
  "user_id": 1,
  "amount": 8800,
  "order_no": "WTB20260601123456"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "balance_before": 50000,
    "balance_after": 41200
  }
}
```

**错误**：
| code | message |
|------|---------|
| 40003 | 余额不足 |

---

### 1.7 余额退款（内部接口）

```
POST /api/user/balance/refund
内部接口，服务间调用，不走网关
```

**请求体**：
```json
{
  "user_id": 1,
  "amount": 8800,
  "order_no": "WTB20260601123456",
  "remark": "订单退款"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "balance_before": 41200,
    "balance_after": 50000
  }
}
```

---

### 1.8 宠物列表

```
GET /api/user/pets
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "name": "旺财",
      "breed": "金毛",
      "weight": 28.5,
      "birthday": "2024-03-15",
      "created_at": "2026-01-20T08:00:00Z"
    }
  ]
}
```

---

### 1.9 添加宠物

```
POST /api/user/pets
需要认证
```

**请求体**：
```json
{
  "name": "旺财",
  "breed": "金毛",
  "weight": 28.5,
  "birthday": "2024-03-15"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "name": "旺财",
    "breed": "金毛",
    "weight": 28.5,
    "birthday": "2024-03-15",
    "created_at": "2026-01-20T08:00:00Z"
  }
}
```

---

### 1.10 内部查询用户

```
GET /api/user/internal/:id
内部接口，服务间调用
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "member_level": 1,
    "balance": 50000
  }
}
```

---

## 2. 座位服务 (seat-service) — 端口 8082

### 2.1 区域列表

```
GET /api/seat/areas
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "name": "室内A区",
      "sort_order": 1,
      "seat_count": 8
    }
  ]
}
```

---

### 2.2 新增区域

```
POST /api/seat/areas
需要管理员认证
```

**请求体**：
```json
{
  "name": "室内A区",
  "sort_order": 1
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "name": "室内A区",
    "sort_order": 1
  }
}
```

---

### 2.3 座位列表

```
GET /api/seat/list?area_id=1
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "area_id": 1,
      "name": "A1",
      "type": "normal",
      "capacity": 4,
      "qrcode_url": "https://.../seat/1.png",
      "status": "available",
      "created_at": "2026-01-15T10:00:00Z"
    }
  ]
}
```

---

### 2.4 座位详情

```
GET /api/seat/:id
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "area_id": 1,
    "area_name": "室内A区",
    "name": "A1",
    "type": "normal",
    "capacity": 4,
    "qrcode_url": "https://.../seat/1.png",
    "status": "available",
    "status_logs": [
      {
        "old_status": "occupied",
        "new_status": "available",
        "order_id": 1001,
        "changed_at": "2026-06-01T13:00:00Z"
      }
    ]
  }
}
```

---

### 2.5 批量生成二维码

```
POST /api/seat/qrcode/batch
需要管理员认证
```

**请求体**：
```json
{
  "seat_ids": [1, 2, 3]
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": [
    { "seat_id": 1, "qrcode_url": "https://.../seat/1.png" },
    { "seat_id": 2, "qrcode_url": "https://.../seat/2.png" },
    { "seat_id": 3, "qrcode_url": "https://.../seat/3.png" }
  ]
}
```

---

### 2.6 扫码解析

```
GET /api/seat/scan?code=SEAT_CODE_1
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "seat_id": 1,
    "area_id": 1,
    "area_name": "室内A区",
    "seat_name": "A1",
    "status": "available"
  }
}
```

**错误**：
| code | message |
|------|---------|
| 40001 | 无效的二维码 |
| 40004 | 座位已被占用 |

---

### 2.7 内部查询座位

```
GET /api/seat/internal/:id
内部接口
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "area_id": 1,
    "name": "A1",
    "status": "available"
  }
}
```

---

## 3. 菜品服务 (menu-service) — 端口 8083

### 3.1 分类树

```
GET /api/menu/categories
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "name": "主食",
      "sort_order": 1,
      "children": [
        { "id": 2, "name": "饭类", "parent_id": 1, "sort_order": 1 },
        { "id": 3, "name": "面类", "parent_id": 1, "sort_order": 2 }
      ]
    },
    {
      "id": 4,
      "name": "饮品",
      "sort_order": 2,
      "children": []
    }
  ]
}
```

---

### 3.2 菜品列表

```
GET /api/menu/dishes?category_id=2&page=1&pageSize=20
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "total": 15,
    "page": 1,
    "pageSize": 20,
    "list": [
      {
        "id": 1,
        "category_id": 2,
        "name": "红烧肉饭",
        "subtitle": "秘制酱汁",
        "images": ["https://..."],
        "tags": ["热门", "推荐"],
        "prices": [
          { "price_type": "normal", "price": 3800 },
          { "price_type": "member", "price": 3200 }
        ],
        "stock": {
          "daily_limit": 50,
          "sold_count": 12,
          "available": 38
        }
      }
    ]
  }
}
```

---

### 3.3 菜品详情

```
GET /api/menu/dish/:id
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "category_id": 2,
    "category_name": "饭类",
    "name": "红烧肉饭",
    "subtitle": "秘制酱汁",
    "description": "精选五花肉，慢炖两小时...",
    "images": ["https://..."],
    "tags": ["热门", "推荐"],
    "prices": [
      { "price_type": "normal", "price": 3800 },
      { "price_type": "member", "price": 3200 },
      { "price_type": "time_slot", "price": 2800, "start_time": "14:00", "end_time": "17:00" }
    ],
    "stock": {
      "daily_limit": 50,
      "sold_count": 12,
      "available": 38
    },
    "status": 1,
    "created_at": "2026-01-15T10:00:00Z"
  }
}
```

---

### 3.4 搜索

```
GET /api/menu/search?q=红烧肉
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "name": "红烧肉饭",
      "subtitle": "秘制酱汁",
      "images": ["https://..."],
      "prices": [{ "price_type": "normal", "price": 3800 }]
    },
    {
      "id": 5,
      "name": "红烧牛肉面",
      "subtitle": "",
      "images": ["https://..."],
      "prices": [{ "price_type": "normal", "price": 2800 }]
    }
  ]
}
```

---

### 3.5 批量查询菜品（内部接口）

```
POST /api/menu/dishes/batch
内部接口
```

**请求体**：
```json
{
  "dish_ids": [1, 2, 3]
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "name": "红烧肉饭",
      "prices": [
        { "price_type": "normal", "price": 3800 },
        { "price_type": "member", "price": 3200 }
      ],
      "stock": { "daily_limit": 50, "sold_count": 12 }
    }
  ]
}
```

---

### 3.6 后台管理 — 分类 CRUD

```
POST   /api/menu/admin/category      新增分类
PUT    /api/menu/admin/category/:id   更新分类
DELETE /api/menu/admin/category/:id   删除分类
需要管理员认证
```

**请求体 (POST/PUT)**：
```json
{
  "name": "饭类",
  "parent_id": 1,
  "sort_order": 1
}
```

---

### 3.7 后台管理 — 菜品 CRUD

```
POST   /api/menu/admin/dish          新增菜品
PUT    /api/menu/admin/dish/:id      更新菜品
DELETE /api/menu/admin/dish/:id      删除菜品
需要管理员认证
```

**请求体 (POST/PUT)**：
```json
{
  "category_id": 2,
  "name": "红烧肉饭",
  "subtitle": "秘制酱汁",
  "description": "精选五花肉...",
  "images": ["https://..."],
  "tags": ["热门", "推荐"],
  "prices": [
    { "price_type": "normal", "price": 3800 },
    { "price_type": "member", "price": 3200 }
  ]
}
```

---

### 3.8 后台管理 — 设置库存

```
POST /api/menu/admin/stock
需要管理员认证
```

**请求体**：
```json
{
  "dish_id": 1,
  "date": "2026-06-01",
  "daily_limit": 50
}
```

---

## 4. 营销定价服务 (pricing-service) — 端口 8088

### 4.1 计算订单价格

```
POST /api/pricing/calculate
需要认证
```

**请求体**：
```json
{
  "user_level": 1,
  "items": [
    { "dish_id": 1, "quantity": 2 },
    { "dish_id": 2, "quantity": 1 }
  ],
  "coupon_id": ""
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "items": [
      {
        "dish_id": 1,
        "dish_name": "红烧肉饭",
        "quantity": 2,
        "unit_price": 3200,
        "subtotal": 6400,
        "price_type": "member"
      },
      {
        "dish_id": 2,
        "dish_name": "可乐",
        "quantity": 1,
        "unit_price": 800,
        "subtotal": 800,
        "price_type": "normal"
      }
    ],
    "total_amount": 7200,
    "discount_amount": 1000,
    "final_amount": 6200,
    "applied_promotions": [
      {
        "promotion_id": 1,
        "name": "满50减10",
        "discount": 1000
      }
    ]
  }
}
```

**错误**：
| code | message |
|------|---------|
| 40001 | 菜品ID不存在 |
| 50001 | 调用菜品服务失败 |

---

### 4.2 菜品当前价

```
GET /api/pricing/dish/:id?user_level=1
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "dish_id": 1,
    "price": 3200,
    "price_type": "member",
    "original_price": 3800
  }
}
```

---

### 4.3 营销活动列表

```
GET /api/pricing/promotions
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "name": "满50减10",
      "type": "full_reduction",
      "config": {
        "threshold": 5000,
        "reduce": 1000
      },
      "start_time": "2026-06-01T00:00:00Z",
      "end_time": "2026-06-30T23:59:59Z"
    }
  ]
}
```

---

### 4.4 后台管理

```
POST /api/pricing/admin/rule       新增价格规则
POST /api/pricing/admin/promotion  新增营销活动
POST /api/pricing/admin/combo      新增套餐
需要管理员认证
```

**新增价格规则 — 请求体**：
```json
{
  "dish_id": 1,
  "rule_type": "member",
  "price": 3200
}
```

**新增满减活动 — 请求体**：
```json
{
  "name": "满50减10",
  "type": "full_reduction",
  "config_json": "{\"threshold\": 5000, \"reduce\": 1000}",
  "start_time": "2026-06-01T00:00:00Z",
  "end_time": "2026-06-30T23:59:59Z"
}
```

**新增套餐 — 请求体**：
```json
{
  "name": "双人套餐",
  "price": 6800,
  "dish_list": "[{\"dish_id\": 1, \"quantity\": 2}, {\"dish_id\": 3, \"quantity\": 2}]"
}
```

---

## 5. 订单服务 (order-service) — 端口 8084

### 5.1 加入购物车

```
POST /api/order/cart/add
需要认证
```

**请求体**：
```json
{
  "seat_id": 1,
  "dish_id": 1,
  "quantity": 2,
  "remark": "少辣"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "cart_items": [
      {
        "dish_id": 1,
        "dish_name": "红烧肉饭",
        "quantity": 2,
        "unit_price": 3800,
        "remark": "少辣"
      }
    ]
  }
}
```

---

### 5.2 购物车列表

```
GET /api/order/cart/list?seat_id=1
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "seat_id": 1,
    "items": [
      {
        "user_id": 1,
        "user_nickname": "小明",
        "dish_id": 1,
        "dish_name": "红烧肉饭",
        "quantity": 2,
        "unit_price": 3800,
        "remark": "少辣"
      },
      {
        "user_id": 2,
        "user_nickname": "小红",
        "dish_id": 2,
        "dish_name": "可乐",
        "quantity": 1,
        "unit_price": 800,
        "remark": ""
      }
    ],
    "total_items": 3,
    "total_amount": 8400
  }
}
```

---

### 5.3 修改购物车

```
PUT /api/order/cart/update
需要认证
```

**请求体**：
```json
{
  "seat_id": 1,
  "dish_id": 1,
  "quantity": 1,
  "remark": "不要葱"
}
```

---

### 5.4 移除购物车商品

```
DELETE /api/order/cart/remove
需要认证
```

**请求体**：
```json
{
  "seat_id": 1,
  "dish_id": 1
}
```

---

### 5.5 提交订单

```
POST /api/order/create
需要认证
```

**请求体**：
```json
{
  "seat_id": 1
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "order_id": 1001,
    "order_no": "WTB20260601123456",
    "status": "pending",
    "items": [
      { "dish_id": 1, "dish_name": "红烧肉饭", "quantity": 2, "unit_price": 3200 }
    ],
    "total_amount": 7200,
    "discount_amount": 1000,
    "pay_amount": 6200,
    "created_at": "2026-06-01T12:34:56Z"
  }
}
```

**错误**：
| code | message |
|------|---------|
| 40002 | 库存不足 |
| 40004 | 座位已被占用 |
| 50001 | 价格计算失败 |

---

### 5.6 发起支付

```
POST /api/order/:id/pay
需要认证
```

**请求体**：
```json
{
  "channel": "wxpay"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "wx_pay_params": {
      "prepay_id": "wx...",
      "nonce_str": "...",
      "sign": "...",
      "timestamp": "1717200000"
    }
  }
}
```

---

### 5.7 订单状态

```
GET /api/order/:id/status
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "order_id": 1001,
    "order_no": "WTB20260601123456",
    "status": "cooking",
    "status_logs": [
      { "from_status": "pending", "to_status": "confirmed", "created_at": "..." },
      { "from_status": "confirmed", "to_status": "cooking", "created_at": "..." }
    ]
  }
}
```

---

### 5.8 我的订单列表

```
GET /api/order/list?page=1&pageSize=20
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "total": 15,
    "page": 1,
    "pageSize": 20,
    "list": [
      {
        "id": 1001,
        "order_no": "WTB20260601123456",
        "seat_name": "A1",
        "status": "completed",
        "total_amount": 7200,
        "pay_amount": 6200,
        "created_at": "2026-06-01T12:34:56Z"
      }
    ]
  }
}
```

---

### 5.9 支付回调（内部接口）

```
POST /api/order/internal/notify
内部接口，由 payment-service 调用
```

**请求体**：
```json
{
  "order_no": "WTB20260601123456",
  "status": "paid",
  "transaction_id": "420000..."
}
```

---

### 5.10 后台 — 订单列表

```
GET /api/order/admin/list?status=cooking&page=1&pageSize=20
需要管理员认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "total": 5,
    "page": 1,
    "pageSize": 20,
    "list": [
      {
        "id": 1001,
        "order_no": "WTB20260601123456",
        "seat_name": "A1",
        "user_nickname": "小明",
        "status": "cooking",
        "items": [
          { "dish_name": "红烧肉饭", "quantity": 2 }
        ],
        "pay_amount": 6200,
        "created_at": "2026-06-01T12:34:56Z"
      }
    ]
  }
}
```

---

### 5.11 后台 — 修改订单状态

```
PUT /api/order/admin/status
需要管理员认证
```

**请求体**：
```json
{
  "order_id": 1001,
  "status": "cooking"
}
```

**状态流转规则**：
```
pending    → confirmed / cancelled
confirmed  → cooking
cooking    → served
served     → paid (如未支付)
paid       → completed
pending    → cancelled（超时自动）
paid       → refunded（管理员退单）
```

**错误**：
| code | message |
|------|---------|
| 40005 | 订单状态不允许此操作 |

---

### 5.12 后台 — 退单

```
POST /api/order/admin/refund
需要管理员认证
```

**请求体**：
```json
{
  "order_id": 1001,
  "reason": "菜品质量问题"
}
```

---

## 6. 支付服务 (payment-service) — 端口 8085

### 6.1 创建支付单

```
POST /api/pay/create
需要认证
```

**请求体**：
```json
{
  "order_no": "WTB20260601123456",
  "amount": 6200,
  "channel": "wxpay"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "out_trade_no": "PAY20260601123456",
    "amount": 6200,
    "channel": "wxpay",
    "status": "pending"
  }
}
```

---

### 6.2 微信统一下单

```
POST /api/pay/wx/prepay
需要认证
```

**请求体**：
```json
{
  "out_trade_no": "PAY20260601123456"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "prepay_id": "wx...",
    "nonce_str": "abc...",
    "sign": "...",
    "timestamp": "1717200000",
    "package": "prepay_id=wx..."
  }
}
```

---

### 6.3 余额支付

```
POST /api/pay/balance
需要认证
```

**请求体**：
```json
{
  "out_trade_no": "PAY20260601123456"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "out_trade_no": "PAY20260601123456",
    "status": "paid",
    "paid_at": "2026-06-01T12:35:00Z"
  }
}
```

**错误**：
| code | message |
|------|---------|
| 40003 | 余额不足 |

---

### 6.4 充值

```
POST /api/pay/recharge
需要认证
```

**请求体**：
```json
{
  "amount": 10000,
  "channel": "wxpay"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "recharge_order_id": 2001,
    "amount": 10000,
    "gifted_amount": 2000,
    "discount_rate": 1.00,
    "final_amount": 12000,
    "wx_pay_params": {
      "prepay_id": "wx...",
      "nonce_str": "...",
      "sign": "...",
      "timestamp": "1717200000"
    }
  }
}
```

---

### 6.5 微信支付回调

```
POST /api/pay/callback/wx
公开接口，微信服务器调用
```

**请求体**（微信回调格式，简化）：
```json
{
  "out_trade_no": "PAY20260601123456",
  "transaction_id": "420000...",
  "status": "SUCCESS"
}
```

**处理逻辑**：
1. 更新 payment_order 状态 → paid
2. 回调 order-service `/api/order/internal/notify`
3. 回调 user-service `/api/user/internal/consume`（增加消费记录）
4. 回调 points-service `/api/points/internal/grant`（发放积分）
5. 如使用了卡券 → 核销

---

### 6.6 退款

```
POST /api/pay/refund
需要管理员认证
```

**请求体**：
```json
{
  "out_trade_no": "PAY20260601123456",
  "amount": 6200,
  "reason": "订单退单"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "refund_no": "REF20260601130000",
    "amount": 6200,
    "status": "pending"
  }
}
```

---

### 6.7 查询支付单

```
GET /api/pay/query/:outTradeNo
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "out_trade_no": "PAY20260601123456",
    "order_no": "WTB20260601123456",
    "amount": 6200,
    "channel": "wxpay",
    "status": "paid",
    "paid_at": "2026-06-01T12:35:00Z"
  }
}
```

---

## 7. 积分服务 (points-service) — 端口 8086

### 7.1 积分余额

```
GET /api/points/account
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "total_points": 1500,
    "used_points": 200,
    "frozen_points": 0,
    "available_points": 1300
  }
}
```

---

### 7.2 积分流水

```
GET /api/points/logs?page=1&pageSize=20
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "total": 10,
    "page": 1,
    "pageSize": 20,
    "list": [
      {
        "id": 1,
        "type": "gain",
        "points": 120,
        "source_id": "WTB20260601123456",
        "remark": "消费获积分 ×1.5(会员)",
        "created_at": "2026-06-01T12:35:00Z"
      },
      {
        "id": 2,
        "type": "exchange",
        "points": -200,
        "source_id": "EXC20260601001",
        "remark": "兑换狗粮一袋",
        "created_at": "2026-06-01T14:00:00Z"
      }
    ]
  }
}
```

---

### 7.3 兑换商品列表

```
GET /api/points/goods
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "name": "狗粮500g",
      "image": "https://...",
      "points_price": 200,
      "stock": 50,
      "type": "physical"
    }
  ]
}
```

---

### 7.4 积分兑换

```
POST /api/points/exchange
需要认证
```

**请求体**：
```json
{
  "goods_id": 1
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "exchange_order_id": 1,
    "goods_name": "狗粮500g",
    "points_cost": 200,
    "points_after": 1100
  }
}
```

**错误**：
| code | message |
|------|---------|
| 40002 | 已售罄 |
| 40007 | 积分不足 |

---

### 7.5 发放积分（内部接口）

```
POST /api/points/internal/grant
内部接口，订单支付完成后调用
```

**请求体**：
```json
{
  "user_id": 1,
  "user_level": 1,
  "amount": 6200,
  "source_id": "WTB20260601123456",
  "source_type": "consumption"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "points_granted": 93,
    "multiplier": 1.5,
    "total_points_after": 1593
  }
}
```

---

### 7.6 后台管理

```
POST /api/points/admin/rule   配置积分规则
POST /api/points/admin/goods  新增兑换商品
需要管理员认证
```

**配置积分规则 — 请求体**：
```json
{
  "name": "消费积分",
  "type": "consumption",
  "config_json": "{\"rate\": 1}"
}
```

**新增兑换商品 — 请求体**：
```json
{
  "name": "狗粮500g",
  "image": "https://...",
  "points_price": 200,
  "stock": 50,
  "type": "physical"
}
```

---

## 8. 活动服务 (activity-service) — 端口 8087

### 8.1 生效公告

```
GET /api/activity/announcements
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "title": "六月充值特惠",
      "content": "充值100送20",
      "type": "text",
      "image": "",
      "link_type": "recharge",
      "link_target": "",
      "sort_order": 1
    }
  ]
}
```

---

### 8.2 充值折扣配置

```
GET /api/activity/recharge-discount
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "discount_rate": 0.90,
    "description": "全场充值9折",
    "start_time": "2026-06-01T00:00:00Z",
    "end_time": "2026-06-30T23:59:59Z"
  }
}
```

---

### 8.3 活动列表

```
GET /api/activity/list?page=1&pageSize=20
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "total": 3,
    "page": 1,
    "pageSize": 20,
    "list": [
      {
        "id": 1,
        "title": "狗狗聚会日",
        "description": "带上你的狗狗来玩...",
        "image": "https://...",
        "max_participants": 30,
        "current_participants": 18,
        "event_time": "2026-06-15T14:00:00Z",
        "location": "户外草坪",
        "status": "published"
      }
    ]
  }
}
```

---

### 8.4 活动报名

```
POST /api/activity/:id/register
需要认证
```

**请求体**：
```json
{
  "name": "小明",
  "phone": "13800138000",
  "remark": "带一只金毛"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "registration_id": 1,
    "activity_title": "狗狗聚会日",
    "status": "registered"
  }
}
```

**错误**：
| code | message |
|------|---------|
| 40006 | 活动名额已满 |
| 40008 | 已报名 |

---

### 8.5 我的报名

```
GET /api/activity/my-registrations
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "activity_id": 1,
      "activity_title": "狗狗聚会日",
      "status": "registered",
      "created_at": "2026-06-01T10:00:00Z"
    }
  ]
}
```

---

### 8.6 取消报名

```
PUT /api/activity/:id/cancel
需要认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "activity_id": 1,
    "status": "cancelled"
  }
}
```

---

### 8.7 内部查折扣

```
GET /api/activity/internal/discount
内部接口
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "discount_rate": 0.90
  }
}
```

---

### 8.8 后台管理

```
POST /api/activity/admin/announcement  新增公告
POST /api/activity/admin/activity      新增活动
需要管理员认证
```

**新增公告 — 请求体**：
```json
{
  "title": "六月充值特惠",
  "content": "充值100送20",
  "type": "text",
  "link_type": "recharge",
  "start_time": "2026-06-01T00:00:00Z",
  "end_time": "2026-06-30T23:59:59Z"
}
```

**新增活动 — 请求体**：
```json
{
  "title": "狗狗聚会日",
  "description": "带上你的狗狗来玩...",
  "image": "https://...",
  "max_participants": 30,
  "event_time": "2026-06-15T14:00:00Z",
  "location": "户外草坪"
}
```

---

## 9. 数据统计服务 (analytics-service) — 端口 8089

### 9.1 仪表盘

```
GET /api/analytics/dashboard
需要管理员认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "today": {
      "revenue": 56000,
      "orders": 48,
      "avg_order": 1167,
      "new_users": 5
    },
    "this_month": {
      "revenue": 1200000,
      "orders": 980,
      "avg_order": 1224
    }
  }
}
```

---

### 9.2 营收统计

```
GET /api/analytics/revenue?start_date=2026-06-01&end_date=2026-06-30
需要管理员认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "total_revenue": 1200000,
    "total_orders": 980,
    "avg_order": 1224,
    "daily": [
      { "date": "2026-06-01", "revenue": 42000, "orders": 35 },
      { "date": "2026-06-02", "revenue": 38000, "orders": 30 }
    ]
  }
}
```

---

### 9.3 菜品销量排行

```
GET /api/analytics/dishes?start_date=2026-06-01&end_date=2026-06-30&limit=10
需要管理员认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": [
    {
      "dish_id": 1,
      "dish_name": "红烧肉饭",
      "total_quantity": 320,
      "total_amount": 1216000
    }
  ]
}
```

---

### 9.4 会员分布

```
GET /api/analytics/members
需要管理员认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "total_users": 500,
    "normal": 300,
    "member": 150,
    "recharge": 50
  }
}
```

---

### 9.5 积分统计

```
GET /api/analytics/points
需要管理员认证
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "total_granted": 50000,
    "total_exchanged": 12000,
    "avg_per_user": 100
  }
}
```

---

### 9.6 Excel 导出

```
POST /api/analytics/export
需要管理员认证
```

**请求体**：
```json
{
  "type": "revenue",
  "start_date": "2026-06-01",
  "end_date": "2026-06-30"
}
```

**响应体**：
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "file_url": "https://.../exports/revenue_202606.xlsx"
  }
}
```

---

## 附录A：内部服务调用汇总

当 payment-service 收到微信支付成功回调后，需依次调用：

```
1. PATCH payment_order.status = "paid"

2. POST http://localhost:8084/api/order/internal/notify
   请求: { "order_no": "...", "status": "paid", "transaction_id": "..." }

3. POST http://localhost:8081/api/user/consumption  (内部接口，更新消费记录)
   请求: { "user_id": 1, "order_id": 1001, "amount": 6200, "dish_count": 3 }

4. POST http://localhost:8086/api/points/internal/grant
   请求: { "user_id": 1, "user_level": 1, "amount": 6200, "source_id": "WTB...", "source_type": "consumption" }

5. 如有卡券核销: wechat.Coupon.Use(openid, coupon_id, order_no)
```

---

> 文档版本 V1.0 | 与 DEVELOPMENT_PLAN.md V4.0 配套
> 每个接口开发完成后，可用本文档对照验证请求/响应格式。
