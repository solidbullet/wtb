# 汪托帮点餐系统 — 代码规范

> 适用：backend 单体及 services/ 下各业务代码包 | 基于 Go 1.24 + Gin + GORM (PostgreSQL)
> 目标执行者：Kimi / AI 编程助手
> 所有业务代码包必须严格遵循此规范，确保代码风格一致

---

## 1. 项目结构规范

每个业务代码包目录结构：

```
services/{name}/
├── go.mod                     # module github.com/wtb-ordering/services/{name}
├── main.go                    # 入口：初始化配置→DB→Redis→启动HTTP
├── config/
│   └── config.go              # 配置结构体 + Load() 函数
├── model/
│   └── {entity}.go            # 每个实体一个文件，GORM tag
├── repository/
│   └── {entity}_repo.go       # 数据访问层，只做数据库操作
├── service/
│   └── {entity}_service.go    # 业务逻辑层，调用 repository + 其他服务
├── handler/
│   └── {domain}_handler.go    # HTTP handler，解析请求→调用 service→返回响应
├── router/
│   └── router.go              # 路由注册：func SetupRouter(deps) *gin.Engine
├── migration/
│   └── migrate.go             # func AutoMigrate(db *gorm.DB) error
├── {file}_test.go             # 单元测试（与源文件同目录）
└── TEST_REPORT.md             # 测试报告（每个服务完成后填写）
```

### 禁止事项

- ❌ handler 里直接写 SQL
- ❌ model 里写业务逻辑
- ❌ 在 main.go 里写超过 50 行的初始化代码
- ❌ 硬编码配置（端口、DSN、密钥等）

---

## 2. 命名规范

### 2.1 文件

```
model/user.go               → 小写+下划线
repository/user_repo.go     → 小写+下划线
service/user_service.go     → 小写+下划线
handler/user_handler.go     → 小写+下划线
handler/user_handler_test.go → 测试文件加 _test
```

### 2.2 Go 标识符

```go
// 类型名：大写驼峰
type UserService struct {}
type CreateOrderRequest struct {}

// 函数/方法：大写驼峰（导出）或小写驼峰（未导出）
func NewUserService(repo *UserRepo) *UserService {}
func (s *UserService) findByOpenID(openid string) (*model.User, error) {}

// 常量
const DefaultPageSize = 20
const OrderStatusPending = "pending"

// 变量
var jwtSecret []byte
var serviceURLs map[string]string
```

### 2.3 数据库

```sql
-- 表名：小写+下划线+复数
users, order_items, recharge_records

-- 字段名：小写+下划线
created_at, member_level, out_trade_no

-- 索引名：idx_{table}_{field}
CREATE INDEX idx_users_openid ON users(openid);
CREATE INDEX idx_orders_user ON orders(user_id);
```

### 2.4 API 路径

```
/api/{service}/{resource}        → 公开接口
/api/{service}/internal/{action} → 内部接口（走服务间直连）
/api/{service}/admin/{resource}  → 后台管理接口（需管理员角色）
```

---

## 3. 分层规范

### 3.1 Model 层

每个实体对应一个 struct，使用 GORM tag：

```go
package model

import "time"

type User struct {
    ID               uint      `gorm:"primaryKey" json:"id"`
    OpenID           string    `gorm:"uniqueIndex;size:64" json:"openid"`
    Nickname         string    `gorm:"size:100" json:"nickname"`
    MemberLevel      int16     `gorm:"default:0" json:"member_level"`
    Balance          int       `gorm:"default:0" json:"balance"`
    TotalConsumption int       `gorm:"default:0" json:"total_consumption"`
    TotalOrders      int       `gorm:"default:0" json:"total_orders"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }
```

**规范**：
- 所有 model 必须定义 `TableName()` 方法，返回复数小写表名
- `json` tag 使用下划线命名
- `gorm` tag 定义索引、默认值、约束
- 金额字段统一用 `int`（单位：分）

### 3.2 Repository 层

每个实体一个 repository，只做数据库操作，不包含业务逻辑：

```go
package repository

import (
    "gorm.io/gorm"
    "github.com/wtb-ordering/services/user/model"
)

type UserRepo struct {
    db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
    return &UserRepo{db: db}
}

func (r *UserRepo) FindByOpenID(openid string) (*model.User, error) {
    var user model.User
    err := r.db.Where("openid = ?", openid).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepo) Create(user *model.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepo) UpdateBalance(userID uint, amount int) error {
    return r.db.Model(&model.User{}).Where("id = ?", userID).
        UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *UserRepo) FindByID(id uint) (*model.User, error) {
    var user model.User
    err := r.db.First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

**规范**：
- 方法名使用 Find/Create/Update/Delete/List 前缀
- 参数化查询，不使用字符串拼接
- 返回 `(*Model, error)` 或 `([]Model, error)`
- 不在此层处理 `gorm.ErrRecordNotFound`（交给 service 层）
- 分页方法签名：`List(page, pageSize int) ([]Model, int64, error)`

### 3.3 Service 层

业务逻辑层，调用 repository + 其他业务代码包：

```go
package service

import (
    "errors"
    "github.com/wtb-ordering/services/user/model"
    "github.com/wtb-ordering/services/user/repository"
)

type UserService struct {
    repo        *repository.UserRepo
    wechatClient *wechat.Client  // 微信 API
    jwtSecret   []byte
}

func NewUserService(repo *repository.UserRepo, wc *wechat.Client, secret []byte) *UserService {
    return &UserService{repo: repo, wechatClient: wc, jwtSecret: secret}
}

func (s *UserService) WxLogin(code string) (string, *model.User, error) {
    // 1. 调用微信 code2session
    session, err := s.wechatClient.Code2Session(code)
    if err != nil {
        return "", nil, errors.New("微信登录失败")
    }
    // 2. 查找或创建用户
    user, err := s.repo.FindByOpenID(session.OpenID)
    if errors.Is(err, gorm.ErrRecordNotFound) {
        user = &model.User{OpenID: session.OpenID}
        if err := s.repo.Create(user); err != nil {
            return "", nil, err
        }
    } else if err != nil {
        return "", nil, err
    }
    // 3. 生成 JWT
    token, err := jwt.GenerateToken(...)
    return token, user, nil
}
```

**规范**：
- 业务错误使用 `errors.New()` 或 `fmt.Errorf()`，描述中文
- 调用其他服务使用 `pkg/httpclient` 包
- 事务操作放 service 层，repository 只提供单条 SQL 操作

### 3.4 Handler 层

```go
package handler

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/wtb-ordering/pkg/response"
    "github.com/wtb-ordering/services/user/service"
)

type UserHandler struct {
    svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
    return &UserHandler{svc: svc}
}

// WxLogin POST /api/user/wx-login
func (h *UserHandler) WxLogin(c *gin.Context) {
    var req struct {
        Code string `json:"code" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 40001, "参数错误: code 不能为空")
        return
    }

    token, user, err := h.svc.WxLogin(req.Code)
    if err != nil {
        response.Error(c, 50001, err.Error())
        return
    }

    response.Success(c, gin.H{
        "token": token,
        "user":  user,
    })
}

// GetProfile GET /api/user/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
    userID := c.GetString("user_id") // 从 JWT 中间件注入
    id, _ := strconv.ParseUint(userID, 10, 64)

    user, err := h.svc.GetProfile(uint(id))
    if err != nil {
        response.Error(c, 50001, "获取用户信息失败")
        return
    }

    response.Success(c, user)
}
```

**规范**：
- 每个 handler 方法上方注释：`// MethodName METHOD /path`
- 请求体绑定使用 `c.ShouldBindJSON(&req)`
- 路径参数使用 `c.Param("id")`
- 查询参数使用 `c.Query("page")`、`c.DefaultQuery("pageSize", "20")`
- JWT 中的 user_id 从 `c.GetString("user_id")` 获取
- 所有响应统一走 `response.Success()` / `response.Error()` / `response.SuccessPage()`

### 3.5 Router 层

```go
package router

import (
    "github.com/gin-gonic/gin"
    "github.com/wtb-ordering/pkg/jwt"
    "github.com/wtb-ordering/services/user/handler"
)

func SetupRouter(h *handler.UserHandler, jwtSecret []byte) *gin.Engine {
    r := gin.Default()

    // 健康检查
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    api := r.Group("/api/user")
    {
        // 公开接口
        api.POST("/wx-login", h.WxLogin)

        // 需要认证
        auth := api.Group("")
        auth.Use(jwt.AuthMiddleware(jwtSecret))
        {
            auth.GET("/profile", h.GetProfile)
            auth.GET("/consumption", h.GetConsumption)
            auth.GET("/consumption/summary", h.GetConsumptionSummary)
            auth.POST("/recharge", h.Recharge)
            auth.GET("/pets", h.ListPets)
            auth.POST("/pets", h.AddPet)
        }

        // 内部接口（不走网关，服务间直连）
        internal := api.Group("/internal")
        {
            internal.POST("/balance/deduct", h.DeductBalance)
            internal.POST("/balance/refund", h.RefundBalance)
            internal.GET("/:id", h.GetUserInternal)
        }
    }

    return r
}
```

### 3.6 main.go

```go
package main

import (
    "fmt"
    "log"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/wtb-ordering/services/user/config"
    "github.com/wtb-ordering/services/user/model"
    "github.com/wtb-ordering/services/user/repository"
    "github.com/wtb-ordering/services/user/service"
    "github.com/wtb-ordering/services/user/handler"
    "github.com/wtb-ordering/services/user/router"
    "github.com/wtb-ordering/pkg/jwt"
)

func main() {
    cfg := config.Load()

    // 数据库
    db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect database: %v", err)
    }
    db.AutoMigrate(&model.User{}, &model.RechargeRecord{}, /* ... */)

    // 依赖注入
    jwt.Init(cfg.JWTSecret)
    userRepo := repository.NewUserRepo(db)
    userSvc := service.NewUserService(userRepo, cfg.Wechat)
    userHandler := handler.NewUserHandler(userSvc)

    // 启动
    r := router.SetupRouter(userHandler, cfg.JWTSecret)
    addr := fmt.Sprintf(":%d", cfg.Port)
    log.Printf("user-service starting on %s", addr)
    if err := r.Run(addr); err != nil {
        log.Fatalf("failed to start server: %v", err)
    }
}
```

---

## 4. 错误处理规范

### 4.1 错误码统一

```go
const (
    CodeSuccess       = 200
    CodeBadRequest    = 40001  // 参数错误
    CodeOutOfStock    = 40002  // 库存不足
    CodeBalanceLow    = 40003  // 余额不足
    CodeSeatOccupied  = 40004  // 座位已被占用
    CodeInvalidStatus = 40005  // 订单状态不允许此操作
    CodeActivityFull  = 40006  // 活动名额已满
    CodeLowPoints     = 40007  // 积分不足
    CodeDuplicated    = 40008  // 不可重复操作
    CodeRechargeFail  = 40009  // 充值失败
    CodePayFail       = 40010  // 支付失败
    CodeUnauthorized  = 40101  // 未登录
    CodeTokenExpired  = 40102  // Token过期
    CodeForbidden     = 40301  // 非管理员
    CodeInternalError = 50001  // 内部错误
)
```

### 4.2 错误传递

```go
// repository: 直接返回 gorm 错误
func (r *UserRepo) FindByOpenID(openid string) (*User, error) {
    var user User
    err := r.db.Where("openid = ?", openid).First(&user).Error
    return &user, err  // 调用方判断 gorm.ErrRecordNotFound
}

// service: 包装为业务错误
func (s *UserService) DeductBalance(userID uint, amount int) error {
    user, err := s.repo.FindByID(userID)
    if err != nil {
        return fmt.Errorf("用户不存在")
    }
    if user.Balance < amount {
        return fmt.Errorf("余额不足")  // handler 层转换为 40003
    }
    return s.repo.UpdateBalance(userID, -amount)
}

// handler: 转换为 HTTP 响应
func (h *UserHandler) DeductBalance(c *gin.Context) {
    // ...
    if err := h.svc.DeductBalance(req.UserID, req.Amount); err != nil {
        if err.Error() == "余额不足" {
            response.Error(c, 40003, err.Error())
            return
        }
        response.Error(c, 50001, err.Error())
        return
    }
    response.Success(c, result)
}
```

---

## 5. 测试规范

### 5.1 单元测试

```go
// handler/user_handler_test.go
package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

func TestWxLogin_Success(t *testing.T) {
    // 1. 创建 Mock service
    mockSvc := &mocks.UserService{}  // 或用接口
    // 2. 创建 handler
    h := NewUserHandler(mockSvc)
    // 3. 构建请求
    body := map[string]string{"code": "test_code"}
    jsonBody, _ := json.Marshal(body)
    req := httptest.NewRequest("POST", "/api/user/wx-login", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    // 4. 创建 gin context
    c, _ := gin.CreateTestContext(w)
    c.Request = req
    // 5. 调用 handler
    h.WxLogin(c)
    // 6. 断言
    if w.Code != 200 {
        t.Errorf("expected 200, got %d", w.Code)
    }
}
```

### 5.2 测试命名

```
Test{FunctionName}_{Scenario}          → TestWxLogin_Success
Test{FunctionName}_{Scenario}_Error    → TestWxLogin_InvalidCode
Test{FunctionName}_{EdgeCase}          → TestDeductBalance_Insufficient
```

### 5.3 覆盖率目标

| 层 | 覆盖率目标 |
|-----|----------|
| model | 不需要 |
| repository | > 70% |
| service | > 75% |
| handler | > 70% |
| pricing-service | > 80%（核心业务） |

---

## 6. 配置规范

### 6.1 配置文件结构

```go
package config

import (
    "os"
)

type Config struct {
    Port      int
    DSN       string       // PostgreSQL DSN
    RedisAddr string       // Redis 地址
    JWTSecret string
    Wechat    WechatConfig
    Services  ServicesConfig
}

type WechatConfig struct {
    AppID     string
    AppSecret string
    MchID     string
    APIv3Key  string
}

type ServicesConfig struct {
    UserServiceURL      string
    SeatServiceURL      string
    MenuServiceURL      string
    OrderServiceURL     string
    PaymentServiceURL   string
    PointsServiceURL    string
    ActivityServiceURL  string
    PricingServiceURL   string
    AnalyticsServiceURL string
}

func Load() *Config {
    return &Config{
        Port:      getEnvInt("PORT", 8081),
        DSN:       getEnv("DATABASE_DSN", "host=/tmp user=admin dbname=wtb_user sslmode=disable TimeZone=Asia/Shanghai"),
        RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
        JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-in-production"),
        Wechat: WechatConfig{
            AppID:     getEnv("WECHAT_APPID", ""),
            AppSecret: getEnv("WECHAT_SECRET", ""),
        },
        Services: ServicesConfig{
            UserServiceURL:     getEnv("USER_SERVICE_URL", "http://localhost:8081"),
            SeatServiceURL:     getEnv("SEAT_SERVICE_URL", "http://localhost:8082"),
            MenuServiceURL:     getEnv("MENU_SERVICE_URL", "http://localhost:8083"),
            OrderServiceURL:    getEnv("ORDER_SERVICE_URL", "http://localhost:8084"),
            PaymentServiceURL:  getEnv("PAYMENT_SERVICE_URL", "http://localhost:8085"),
            PointsServiceURL:   getEnv("POINTS_SERVICE_URL", "http://localhost:8086"),
            ActivityServiceURL: getEnv("ACTIVITY_SERVICE_URL", "http://localhost:8087"),
            PricingServiceURL:  getEnv("PRICING_SERVICE_URL", "http://localhost:8088"),
            AnalyticsServiceURL: getEnv("ANALYTICS_SERVICE_URL", "http://localhost:8089"),
        },
    }
}

func getEnv(key, defaultVal string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
    if val := os.Getenv(key); val != "" {
        if i, err := strconv.Atoi(val); err == nil {
            return i
        }
    }
    return defaultVal
}
```

### 6.2 环境变量优先级

1. 系统环境变量（生产环境）
2. `.env` 文件（开发环境）
3. 代码默认值（本地开发）

---

## 7. Git 提交规范

```
feat(service): 描述
fix(service): 描述
test(service): 描述
docs: 描述
chore: 描述

示例：
feat(user): 完成用户服务所有API + 测试
fix(order): 修复购物车并发安全问题
test(pricing): 补充优惠叠加逻辑测试
docs: 更新API文档
chore: 初始化项目脚手架
```

---

## 8. 依赖清单（go.mod）

每个业务代码包的标准依赖：

```
require (
    github.com/gin-gonic/gin          v1.10+     // HTTP 框架
    gorm.io/gorm                       v1.26+     // ORM
    gorm.io/driver/postgres            v1.5+      // PG 驱动
    github.com/golang-jwt/jwt/v5       v5.2+      // JWT
    github.com/redis/go-redis/v9       v9.5+      // Redis 客户端（购物车/锁）
    github.com/xuri/excelize/v2        v2.8+      // Excel 导出（仅 analytics）
)
```

---

> 所有代码必须符合此规范。发现不符合的代码块，立即修正。
