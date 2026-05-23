# AI 数据查询接入方案（豆包/Coze）

> 目标：让客户（餐厅老板）无需开发管理后台，直接通过 豆包 App 用自然语言查询营业数据。
> 例如："今天的充值总额是多少？"、"今天卖了多少单？"、"最受欢迎的菜品是什么？"

---

## 一、技术现状与方案选择

### 关键限制（必读）

**Coze（扣子）当前不支持 MCP 协议**（截至 2025 Q4 仍在 Roadmap 中，尚未实现）。  
豆包通过 Coze 发布 Bot，因此：

- ❌ **MCP Server 无法直接被豆包使用**（今天不可用）
- ✅ **Coze OpenAPI 插件 可以立即工作**（推荐方案）

### 双轨策略

| 方案 | 用途 | 优先级 |
|------|------|--------|
| **Path A: Coze OpenAPI 插件** | 豆包 App 对话查询（现在就能用） | **P0 - 必须实现** |
| **Path B: MCP Server (SSE)** | Kimi/Claude 等客户端使用，未来 Coze 支持后可迁移 | P1 - 有时间再做 |

**本文档主要描述 Path A 的实现**，Path B 作为附录提供。

---

## 二、整体架构

```
┌─────────────┐      自然语言      ┌─────────────┐      HTTP API      ┌─────────────────────┐
│  客户手机    │ ───────────────→ │   豆包 App   │ ───────────────→ │  analytics-service  │
│  (豆包 Bot) │  "今天营收多少"   │  (Coze Bot) │  X-API-Key 认证   │   (port 8089)       │
└─────────────┘                   └─────────────┘                   └─────────────────────┘
                                                                             │
                                                                             ▼
                                                                   ┌──────────────────┐
                                                                   │  5 个 PostgreSQL  │
                                                                   │  聚合查询 (只读)  │
                                                                   └──────────────────┘
```

### 为什么选择 analytics-service？

- 项目已有 `analytics-service`（端口 8089），但当前所有接口返回 mock 零值
- API_SPEC.md 中已定义了 dashboard/revenue/orders/dishes/members/points 等路由
- 它是跨服务聚合查询的自然落脚点，直接重写即可

---

## 三、数据库连接策略

analytics-service 需要直连多个数据库做聚合查询。

**采用方案：直连多库**（每个数据库维护一个独立连接池）

需要连接的数据库：

```go
var (
    orderDB   *gorm.DB  // wtb_order   - 订单数据
    userDB    *gorm.DB  // wtb_user    - 用户、充值、余额
    paymentDB *gorm.DB  // wtb_payment - 支付、退款
    menuDB    *gorm.DB  // wtb_menu    - 菜品、分类
    pointsDB  *gorm.DB  // wtb_points  - 积分
)
```

> 注意：每种数据库的 DSN 从环境变量读取，格式与 backend 服务一致。

---

## 四、API 接口设计（Coze 插件用）

### 认证方式

所有接口必须在请求头中携带：
```
X-API-Key: <商家密钥>
```

### 接口清单

| 接口 | 方法 | 说明 | 示例问题 |
|------|------|------|----------|
| `/api/analytics/dashboard` | GET | 今日营业总览 | "今天营业情况怎么样？" |
| `/api/analytics/revenue` | GET | 营收统计（支持日期范围） | "今天的充值总额是多少？" |
| `/api/analytics/orders` | GET | 订单统计 | "今天卖了多少单？" |
| `/api/analytics/dishes` | GET | 菜品销售排行 | "最受欢迎的菜品是什么？" |
| `/api/analytics/members` | GET | 会员/用户统计 | "有多少会员？今天新增几个？" |
| `/api/analytics/points` | GET | 积分统计 | "今天积分发放了多少？" |

---

### 4.1 GET /api/analytics/dashboard

**功能**：今日关键指标一键总览，Coze 回答日常询问时优先调用此接口。

**响应：**
```json
{
  "today": {
    "revenue": 158000,
    "orders": 42,
    "avg_order": 3761,
    "new_users": 5,
    "recharge_amount": 50000
  },
  "this_month": {
    "revenue": 1250000,
    "orders": 340,
    "avg_order": 3676
  }
}
```

**对应 SQL（Go 代码中需要执行的查询）：**

```sql
-- 今日营收（分）
SELECT COALESCE(SUM(pay_amount), 0) FROM orders 
WHERE DATE(created_at) = CURRENT_DATE AND status IN ('paid', 'completed');

-- 今日订单数
SELECT COUNT(*) FROM orders 
WHERE DATE(created_at) = CURRENT_DATE AND status IN ('paid', 'completed');

-- 今日客单价（后端做除法，orders > 0 时 revenue/orders）

-- 今日新用户
SELECT COUNT(*) FROM users WHERE DATE(created_at) = CURRENT_DATE;

-- 今日充值总额（分）
SELECT COALESCE(SUM(amount), 0) FROM recharge_records 
WHERE DATE(created_at) = CURRENT_DATE AND status = 'success';

-- 本月营收（同今日，时间范围改为当月）
SELECT COALESCE(SUM(pay_amount), 0) FROM orders 
WHERE DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_DATE) 
  AND status IN ('paid', 'completed');
```

---

### 4.2 GET /api/analytics/revenue

**功能**：营收明细，支持按日期范围查询。

**Query 参数：**
- `start_date` (YYYY-MM-DD, 可选，默认：6 天前)
- `end_date` (YYYY-MM-DD, 可选，默认：今天)

**响应：**
```json
{
  "total_revenue": 158000,
  "total_orders": 42,
  "avg_order_amount": 3761,
  "daily": [
    {"date": "2026-05-17", "revenue": 120000, "orders": 35},
    {"date": "2026-05-18", "revenue": 95000, "orders": 28},
    {"date": "2026-05-23", "revenue": 158000, "orders": 42}
  ]
}
```

**对应 SQL：**
```sql
SELECT 
  DATE(created_at) as date,
  COALESCE(SUM(pay_amount), 0) as revenue,
  COUNT(*) as orders
FROM orders
WHERE created_at BETWEEN ? AND ?
  AND status IN ('paid', 'completed')
GROUP BY DATE(created_at)
ORDER BY date;
```

---

### 4.3 GET /api/analytics/orders

**功能**：订单统计与列表。

**Query 参数：**
- `date` (YYYY-MM-DD, 可选，默认：今天)
- `status` (可选，如：completed, cancelled)

**响应：**
```json
{
  "total": 42,
  "completed": 38,
  "cancelled": 4,
  "total_amount": 158000,
  "list": [
    {
      "order_no": "W20260523001",
      "amount": 8900,
      "status": "completed",
      "created_at": "2026-05-23T12:30:00+08:00"
    }
  ]
}
```

**对应 SQL：**
```sql
-- 总数与状态分布
SELECT 
  COUNT(*) as total,
  COUNT(*) FILTER (WHERE status = 'completed') as completed,
  COUNT(*) FILTER (WHERE status = 'cancelled') as cancelled,
  COALESCE(SUM(pay_amount), 0) as total_amount
FROM orders
WHERE DATE(created_at) = ?;

-- 列表（分页 limit 20）
SELECT order_no, pay_amount as amount, status, created_at
FROM orders
WHERE DATE(created_at) = ?
ORDER BY created_at DESC
LIMIT 20;
```

---

### 4.4 GET /api/analytics/dishes

**功能**：菜品销售排行。

**Query 参数：**
- `date` (YYYY-MM-DD, 可选，默认：今天)
- `limit` (int, 可选，默认：10)

**响应：**
```json
{
  "total_sold": 156,
  "list": [
    {
      "dish_id": 16,
      "dish_name": "红烧肉饭",
      "total_quantity": 25,
      "total_amount": 98750
    }
  ]
}
```

**对应 SQL：**
```sql
SELECT 
  oi.dish_id,
  oi.dish_name,
  SUM(oi.quantity) as total_quantity,
  SUM(oi.quantity * oi.unit_price) as total_amount
FROM order_items oi
JOIN orders o ON oi.order_id = o.id
WHERE DATE(o.created_at) = ? AND o.status IN ('paid', 'completed')
GROUP BY oi.dish_id, oi.dish_name
ORDER BY total_quantity DESC
LIMIT ?;
```

---

### 4.5 GET /api/analytics/members

**功能**：用户与会员统计。

**响应：**
```json
{
  "total_users": 1200,
  "normal": 800,
  "member": 300,
  "recharge": 100,
  "new_today": 5
}
```

**对应 SQL：**
```sql
SELECT 
  COUNT(*) as total_users,
  COUNT(*) FILTER (WHERE member_level = 0) as normal,
  COUNT(*) FILTER (WHERE member_level = 1) as member,
  COUNT(*) FILTER (WHERE member_level = 2) as recharge,
  COUNT(*) FILTER (WHERE DATE(created_at) = CURRENT_DATE) as new_today
FROM users;
```

---

### 4.6 GET /api/analytics/points

**功能**：积分发放与兑换统计。

**响应：**
```json
{
  "total_granted": 5000,
  "total_exchanged": 1200,
  "avg_per_user": 45
}
```

**对应 SQL：**
```sql
-- 今日积分发放
SELECT COALESCE(SUM(points), 0) FROM points_logs 
WHERE DATE(created_at) = CURRENT_DATE AND type = 'gain';

-- 今日积分兑换
SELECT COALESCE(SUM(ABS(points)), 0) FROM points_logs 
WHERE DATE(created_at) = CURRENT_DATE AND type = 'exchange';

-- 人均积分（总持有 / 总人数）
SELECT COALESCE(AVG(total_points), 0) FROM user_points;
```

---

## 五、代码目录结构

基于现有 `services/analytics/` 目录进行改造：

```
services/analytics/
├── main.go                      # 入口：初始化 5 个 DB 连接 + Gin 引擎
├── config/
│   └── config.go                # 多库 DSN 配置、API Key 列表
├── model/
│   └── response.go              # 所有 API 的 JSON 响应结构体
├── middleware/
│   └── auth.go                  # X-API-Key 校验中间件
├── repository/
│   ├── order_repo.go            # wtb_order 查询
│   ├── user_repo.go             # wtb_user 查询
│   ├── payment_repo.go          # wtb_payment 查询（备用）
│   ├── menu_repo.go             # wtb_menu 查询（备用）
│   └── points_repo.go           # wtb_points 查询
├── service/
│   └── analytics_service.go     # 业务逻辑：组合多个 repo 的查询结果
├── handler/
│   ├── analytics_handler.go     # 6 个 HTTP handler
│   └── mcp_handler.go           # 【可选】MCP SSE 端点
├── router/
│   └── router.go                # 路由注册
├── tools/
│   └── openapi_generator.go     # 生成项目根目录 coze_openapi.yaml
└── mcp/
    └── server.go                # 【可选】MCP Server 实现
```

---

## 六、关键代码规范（供 Kimi 自动生成）

### 6.1 环境变量

analytics-service 启动时读取以下环境变量：

```bash
# 各数据库 DSN（host 用本机地址或 EasyTier 虚拟 IP）
ORDER_DB_DSN="host=10.144.144.3 user=admin password=JyRUj7wlNjU0uVHh dbname=wtb_order port=5432 sslmode=disable"
USER_DB_DSN="host=10.144.144.3 user=admin password=JyRUj7wlNjU0uVHh dbname=wtb_user port=5432 sslmode=disable"
PAYMENT_DB_DSN="host=10.144.144.3 user=admin password=JyRUj7wlNjU0uVHh dbname=wtb_payment port=5432 sslmode=disable"
MENU_DB_DSN="host=10.144.144.3 user=admin password=JyRUj7wlNjU0uVHh dbname=wtb_menu port=5432 sslmode=disable"
POINTS_DB_DSN="host=10.144.144.3 user=admin password=JyRUj7wlNjU0uVHh dbname=wtb_points port=5432 sslmode=disable"

# API 密钥（多租户用逗号分隔 key:name）
API_KEYS="ak_wtb_demo:汪托帮 demo 店"

# 服务端口
PORT=8089
```

### 6.2 响应结构体（model/response.go）

```go
package model

type DashboardResponse struct {
    Today     DashboardStats `json:"today"`
    ThisMonth DashboardStats `json:"this_month"`
}

type DashboardStats struct {
    Revenue        int64 `json:"revenue"`         // 分
    Orders         int64 `json:"orders"`
    AvgOrder       int64 `json:"avg_order"`       // 分
    NewUsers       int64 `json:"new_users"`
    RechargeAmount int64 `json:"recharge_amount"` // 分
}

type RevenueResponse struct {
    TotalRevenue   int64         `json:"total_revenue"`
    TotalOrders    int64         `json:"total_orders"`
    AvgOrderAmount int64         `json:"avg_order_amount"`
    Daily          []DailyRevenue `json:"daily"`
}

type DailyRevenue struct {
    Date    string `json:"date"`
    Revenue int64  `json:"revenue"`
    Orders  int64  `json:"orders"`
}

// ... OrdersResponse, DishesResponse, MembersResponse, PointsResponse 类似定义
```

### 6.3 认证中间件（middleware/auth.go）

```go
package middleware

// APIKeyAuth 从请求头读取 X-API-Key，校验是否在允许列表中
func APIKeyAuth(validKeys map[string]string) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := c.GetHeader("X-API-Key")
        if key == "" || validKeys[key] == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid api key"})
            return
        }
        c.Set("merchant_name", validKeys[key])
        c.Next()
    }
}
```

### 6.4 Repository 模式

每个 repo 文件只负责一个数据库的查询，使用原生 SQL 或 GORM Raw：

```go
// repository/order_repo.go
package repository

type OrderRepo struct {
    db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
    return &OrderRepo{db: db}
}

func (r *OrderRepo) GetTodayRevenue() (int64, error) {
    var result int64
    err := r.db.Raw(`
        SELECT COALESCE(SUM(pay_amount), 0) FROM orders 
        WHERE DATE(created_at) = CURRENT_DATE AND status IN ('paid', 'completed')
    `).Scan(&result).Error
    return result, err
}

// ... 其他查询方法
```

### 6.5 Service 层

```go
// service/analytics_service.go
package service

type AnalyticsService struct {
    orderRepo   *repository.OrderRepo
    userRepo    *repository.UserRepo
    pointsRepo  *repository.PointsRepo
}

func (s *AnalyticsService) GetDashboard() (*model.DashboardResponse, error) {
    // 并发查询多个指标
    var revenue, orders, newUsers, recharge int64
    var err error
    
    // 使用 errgroup 或串行查询
    revenue, err = s.orderRepo.GetTodayRevenue()
    // ... 其他查询
    
    avgOrder := int64(0)
    if orders > 0 {
        avgOrder = revenue / orders
    }
    
    return &model.DashboardResponse{
        Today: model.DashboardStats{
            Revenue: revenue,
            Orders: orders,
            AvgOrder: avgOrder,
            NewUsers: newUsers,
            RechargeAmount: recharge,
        },
        // ... this_month 同理
    }, nil
}
```

---

## 七、OpenAPI 规范（coze_openapi.yaml）

项目根目录需要生成 `coze_openapi.yaml`，供 Coze 导入插件使用。  
**要求**：每个接口的 `description` 必须写清楚中文用途和示例问题，帮助 Coze LLM 理解何时调用该接口。

关键片段示例：

```yaml
paths:
  /api/analytics/dashboard:
    get:
      summary: 今日营业总览
      description: |
        查询今日关键营业指标：营收总额、订单数、客单价、新用户、充值总额。
        当用户问"今天营业情况怎么样"、"今天赚了多少钱"、"今天卖了多少单"时调用。
      responses:
        "200":
          description: 成功
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DashboardResponse"
```

> 完整 YAML 由 `tools/openapi_generator.go` 生成，或手动维护在项目根目录。

---

## 八、Coze 接入步骤（部署后操作）

1. **部署 analytics-service** 到公网或 EasyTier 可达地址
2. **登录 [coze.com](https://www.coze.com)**（使用豆包账号）
3. **创建插件**：
   - 点击 "创建插件"
   - 选择 "Import from OpenAPI"
   - 上传项目根目录的 `coze_openapi.yaml`
   - 设置认证类型：API Key
   - Header 名称：`X-API-Key`
   - 填入测试用的 API Key 值
4. **创建 Bot**：
   - 新建 Bot，名称如"汪托帮数据助手"
   - 在"插件"中启用刚才创建的插件
   - 填写提示词（Prompt）：
     ```
     你是餐厅营业数据助手，可以通过调用 API 帮老板查询营业数据。
     你能查询的内容包括：今日营收、订单数、客单价、新用户、充值总额、
     菜品销售排行、会员统计、积分统计等。
     请用中文口语化回答，金额自动除以 100 转换为元并保留两位小数。
     ```
5. **选择发布方式**（推荐用私密模式，见下方"权限控制"章节）：
   - **私密模式（推荐）**：不发布到豆包公共市场，只生成私密分享链接/二维码，发给管理人员
   - **公开模式**：发布到豆包，所有人可搜索（不推荐用于营业数据）
   - **飞书模式**：发布到飞书，在飞书后台配置可见范围（适合已有飞书企业的团队）
6. **测试**：管理人员打开豆包 App，通过分享链接或搜索找到 Bot，问"今天的充值总额是多少？"

---

## 九、权限控制与可见性

营业数据属于敏感信息，必须限制只有管理人员能访问。以下是三种可选方案：

### 方案 A：私密分享链接（最简单，推荐）

**原理**：Bot 不发布到豆包公共市场，只生成私密链接或二维码，通过微信发给管理人员。

**操作**：
1. 在 Coze 创建 Bot 后，**不要点击"发布到豆包"**
2. 选择"生成分享链接"或"生成二维码"
3. 把链接/二维码发给管理人员（如：张经理、李老板、王会计）
4. 他们点击链接 → 自动在豆包中打开 Bot → 即可长期使用

**优点**：
- 外人搜不到这个 Bot
- 零配置，点几下就行

**缺点**：
- 链接如果被转发，拿到的人也能打开 Bot（需要配合 API Key 做第二道防线）

### 方案 B：飞书可见范围控制（最安全）

**原理**：如果管理人员都在飞书上办公，把 Bot 发布到飞书而不是豆包，在飞书后台精确指定谁能使用。

**操作**：
1. Coze 发布渠道选择"飞书"
2. 登录飞书开放平台 → 进入应用详情
3. "版本管理与发布" → "可用范围配置" → 添加指定人员
4. 只有被添加的人员能在飞书 App 里搜到这个 Bot

**优点**：
- 精确到个人，员工离职自动失效
- 外人绝对看不到

**缺点**：
- 需要在飞书 App 中使用，不是在豆包里
- 需要有飞书企业版管理员权限

### 方案 C：双重保险（最佳实践）

**第一层（Bot 可见性）**：使用私密分享链接，避免被公开搜索
**第二层（数据访问）**：API Key 认证，即使链接泄露也查不到数据
**第三层（查询范围）**：只读接口，不能修改数据

```
外人拿到链接 → 打开 Bot → 问"今天营收多少"
                    ↓
            Coze 调用 API（带了 X-API-Key）
                    ↓
            服务器校验 Key → Key 错误 → 返回 401
                    ↓
            豆包告诉用户："查询失败，无权访问"
```

### 多管理人员配置

如果多个管理人员需要各自独立的权限，可以在环境变量中配置多个 API Key：

```bash
API_KEYS="ak_zhangsan:张三,ak_lisi:李四,ak_wangwu:王五"
```

每个 Key 对应一个姓名，服务器日志可以记录是谁在查询，便于审计。

---

## 十、安全设计

1. **只读接口**：所有接口均为 GET，禁止写操作
2. **API Key 认证**：简单有效的访问控制，后续可升级为数据库存储+过期机制
3. **HTTPS 强制**：生产环境必须走 HTTPS
4. **速率限制**：建议每个 API Key 限制 100 次/分钟
5. **时间范围限制**：禁止查询过大的日期范围（如超过 90 天），防止慢查询拖垮 DB

---

## 十一、开发任务清单（Kimi 执行清单）

按以下顺序实现代码：

- [ ] **1. config/config.go**：读取多库 DSN 和 API_KEYS 环境变量
- [ ] **2. model/response.go**：定义 6 个接口的响应结构体
- [ ] **3. middleware/auth.go**：X-API-Key 校验中间件
- [ ] **4. repository/*_repo.go**：为 order/user/points 写具体 SQL 查询方法
- [ ] **5. service/analytics_service.go**：组合 repo 查询，计算 avg 等派生字段
- [ ] **6. handler/analytics_handler.go**：6 个 Gin handler，处理 query 参数，调用 service
- [ ] **7. router/router.go**：注册路由，挂载 auth 中间件
- [ ] **8. main.go**：初始化 5 个 GORM 连接，启动 Gin 服务
- [ ] **9. tools/openapi_generator.go**：生成根目录 coze_openapi.yaml
- [ ] **10. 本地测试**：curl 每个接口验证返回数据正确
- [ ] **11. 【可选】mcp/server.go**：实现 MCP Server (SSE 模式)

---

## 附录：Path B - MCP Server（可选）

如果未来 Coze 支持 MCP，或需要在 Kimi/Claude 中使用，可同时部署 MCP Server。

**暴露的 Tools：**

| Tool 名称 | 说明 |
|-----------|------|
| `get_today_dashboard` | 今日营业总览 |
| `get_revenue_stats` | 营收统计 |
| `get_order_stats` | 订单统计 |
| `get_dish_ranking` | 菜品排行 |
| `get_member_stats` | 会员统计 |
| `get_points_stats` | 积分统计 |

**Transport**：SSE (Server-Sent Events) over HTTP  
**Endpoint**：`/mcp/sse`

实现参考：`mcp/server.go` 使用 `github.com/mark3labs/mcp-go` SDK。
