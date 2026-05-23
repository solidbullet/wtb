# 汪托帮点餐系统 — 项目开发计划 V4.0

> 基于 PRD V4 Go版 | 适配本地开发环境 (macOS 26.3 ARM64, Go 1.24.1, PostgreSQL 18.3)
> 目标执行者：Kimi / 其他 AI 编程助手
> 每个微服务完成后 → 单元测试 → 测试报告

---

## 目录

- [0. 环境准备与项目脚手架](#0-环境准备与项目脚手架)
- [1. 共享包层 (pkg/)](#1-共享包层-pkg)
- [2. 内部微信工具包 (internal/wechat/)](#2-内部微信工具包-internalwechat)
- [3. 用户服务 (user-service)](#3-用户服务-user-service)
- [4. 座位服务 (seat-service)](#4-座位服务-seat-service)
- [5. 菜品服务 (menu-service)](#5-菜品服务-menu-service)
- [6. 营销定价服务 (pricing-service)](#6-营销定价服务-pricing-service)
- [7. 订单服务 (order-service)](#7-订单服务-order-service)
- [8. 支付服务 (payment-service)](#8-支付服务-payment-service)
- [9. 积分服务 (points-service)](#9-积分服务-points-service)
- [10. 活动服务 (activity-service)](#10-活动服务-activity-service)
- [11. 数据统计服务 (analytics-service)](#11-数据统计服务-analytics-service)
- [12. 后台聚合服务 (admin-bff)](#12-后台聚合服务-admin-bff)
- [13. API 网关 (gateway)](#13-api-网关-gateway)
- [14. Docker Compose 集成](#14-docker-compose-集成)
- [15. 前端应用（概要）](#15-前端应用概要)
- [附录A：测试报告模板](#附录a测试报告模板)
- [附录B：每阶段验证命令](#附录b每阶段验证命令)

---

## 执行约定

1. **每个微服务是一个独立的 Go module**，放在 `services/{name}/` 下
2. **每个服务完成后**，必须运行 `go test ./... -v -cover` 并生成测试报告（模板见附录A）
3. **数据库**：每个服务使用独立数据库，命名规则 `wtb_{service}`
4. **端口分配**：每个服务监听不同端口，避免冲突
5. **Git 提交**：每完成一个服务 + 测试通过，做一次 commit
6. **GORM 驱动**：使用 `gorm.io/driver/postgres`，连接串见下方

### 本地环境速查

```
系统:       macOS 26.3 ARM64
Go:         1.24.1  (/opt/homebrew/Cellar/go/1.24.1/libexec)
PostgreSQL: 18.3    (Homebrew, socket /tmp, port 5432, user=admin, 无密码)
Redis:      需安装 (brew install redis)
工作区:     /Users/admin/workspace/jyq/dp
```

### PostgreSQL 连接信息

```
Socket:   /tmp/.s.PGSQL.5432
Host:     localhost (TCP) 或 /tmp (Unix socket)
Port:     5432
User:     admin
Password: (空 — peer 认证)
```

### GORM 连接串 (DSN)

```
host=/tmp user=admin dbname=wtb_user sslmode=disable TimeZone=Asia/Shanghai
```

> 每个服务替换 `dbname=wtb_user` 为对应的数据库名。

---

## 0. 环境准备与项目脚手架

### 0.1 安装缺失工具

```bash
# 1. 安装 Redis（本地开发）
brew install redis
brew services start redis

# 2. 安装 air（Go 热重载）
go install github.com/air-verse/air@latest

# 3. 验证
go version     # go1.24.1
psql -h /tmp -U admin -d postgres -c "SELECT VERSION();"  # PostgreSQL 18.3
redis-cli ping  # PONG
```

### 0.2 创建项目根目录结构

```bash
cd /Users/admin/workspace/jyq/dp
mkdir -p wtb-ordering
cd wtb-ordering

# 初始化根 module（用于共享包）
go mod init github.com/wtb-ordering

mkdir -p pkg/{jwt,response,httpclient}
mkdir -p internal/wechat
mkdir -p services/{user,seat,menu,order,payment,points,activity,pricing,analytics,admin}
mkdir -p gateway
mkdir -p miniprogram
mkdir -p admin-web
mkdir -p configs   # 共享配置文件

touch Makefile
touch docker-compose.yml
```

### 0.3 创建 PostgreSQL 数据库

```bash
for db in user seat menu order payment points activity pricing analytics; do
  psql -h /tmp -U admin -d postgres -c "CREATE DATABASE wtb_${db} OWNER admin ENCODING 'UTF8';" 2>/dev/null
  echo "created wtb_${db}"
done
```

### 0.4 端口分配表

| 服务 | 端口 |
|------|------|
| gateway | 8080 |
| user-service | 8081 |
| seat-service | 8082 |
| menu-service | 8083 |
| order-service | 8084 |
| payment-service | 8085 |
| points-service | 8086 |
| activity-service | 8087 |
| pricing-service | 8088 |
| analytics-service | 8089 |
| admin-bff | 8090 |

### 0.5 检查点 0

完成后验证：
- `go version` → 1.24.1
- `redis-cli ping` → PONG
- `psql -h /tmp -U admin -d postgres -c "\l" | grep wtb_` → 显示 9 个数据库

---

## 1. 共享包层 (pkg/)

> 这三个包被所有微服务依赖，必须先完成。

### 1.1 pkg/response — 统一响应格式

**文件**：`pkg/response/response.go`

```go
package response

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

type PageData struct {
    Total    int64       `json:"total"`
    Page     int         `json:"page"`
    PageSize int         `json:"pageSize"`
    List     interface{} `json:"list"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{Code: 200, Message: "ok", Data: data})
}

func SuccessPage(c *gin.Context, total int64, page, pageSize int, list interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    200,
        Message: "ok",
        Data: PageData{Total: total, Page: page, PageSize: pageSize, List: list},
    })
}

func Error(c *gin.Context, code int, message string) {
    c.JSON(http.StatusOK, Response{Code: code, Message: message, Data: nil})
}

func ErrorWithStatus(c *gin.Context, httpStatus, code int, message string) {
    c.JSON(httpStatus, Response{Code: code, Message: message, Data: nil})
}
```

**测试**：`pkg/response/response_test.go`

```go
package response

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

func TestSuccess(t *testing.T) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    Success(c, map[string]string{"key": "val"})

    if w.Code != http.StatusOK { t.Errorf("expected 200, got %d", w.Code) }
    var resp Response
    json.Unmarshal(w.Body.Bytes(), &resp)
    if resp.Code != 200 { t.Errorf("expected code 200, got %d", resp.Code) }
}

func TestError(t *testing.T) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    Error(c, 40001, "库存不足")

    var resp Response
    json.Unmarshal(w.Body.Bytes(), &resp)
    if resp.Code != 40001 { t.Errorf("expected 40001, got %d", resp.Code) }
}

func TestSuccessPage(t *testing.T) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    SuccessPage(c, 100, 1, 20, []string{"a", "b"})

    var resp Response
    json.Unmarshal(w.Body.Bytes(), &resp)
    if resp.Code != 200 { t.Errorf("expected 200, got %d", resp.Code) }
}
```

### 1.2 pkg/jwt — JWT 工具

**文件**：`pkg/jwt/jwt.go`

```go
package jwt

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   string `json:"user_id"`
    OpenID   string `json:"openid"`
    Level    int    `json:"level"` // 0=普通 1=会员 2=充值
    jwt.RegisteredClaims
}

var jwtSecret []byte

func Init(secret string) { jwtSecret = []byte(secret) }

func GenerateToken(userID, openID string, level int) (string, error) {
    claims := Claims{
        UserID: userID, OpenID: openID, Level: level,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{},
        func(t *jwt.Token) (interface{}, error) { return jwtSecret, nil })
    if err != nil { return nil, err }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, errors.New("invalid token")
}
```

**测试**：`pkg/jwt/jwt_test.go`

```go
package jwt

import (
    "testing"
    "time"
)

func init() { Init("test-secret-key-for-unit-tests") }

func TestGenerateAndParse(t *testing.T) {
    token, err := GenerateToken("u1", "openid_xxx", 1)
    if err != nil { t.Fatalf("generate: %v", err) }

    claims, err := ParseToken(token)
    if err != nil { t.Fatalf("parse: %v", err) }
    if claims.UserID != "u1" { t.Errorf("user_id: %s", claims.UserID) }
    if claims.OpenID != "openid_xxx" { t.Errorf("openid: %s", claims.OpenID) }
    if claims.Level != 1 { t.Errorf("level: %d", claims.Level) }
}

func TestParseInvalidToken(t *testing.T) {
    _, err := ParseToken("invalid.token.here")
    if err == nil { t.Error("expected error for invalid token") }
}

func TestExpiredToken(t *testing.T) {
    claims := Claims{
        UserID: "u1", OpenID: "ox", Level: 0,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
        },
    }
    if claims.RegisteredClaims.ExpiresAt.Time.After(time.Now()) {
        t.Error("expected expired time")
    }
}
```

### 1.3 pkg/httpclient — HTTP 客户端

**文件**：`pkg/httpclient/client.go`

```go
package httpclient

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type Client struct {
    baseURL    string
    httpClient *http.Client
    headers    map[string]string
}

func New(baseURL string) *Client {
    return &Client{
        baseURL: baseURL,
        httpClient: &http.Client{Timeout: 10 * time.Second},
        headers:   map[string]string{"Content-Type": "application/json"},
    }
}

func (c *Client) SetHeader(key, value string) { c.headers[key] = value }

func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
    return c.do(ctx, "GET", path, nil, result)
}

func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
    return c.do(ctx, "POST", path, body, result)
}

func (c *Client) do(ctx context.Context, method, path string, body, result interface{}) error {
    var bodyReader io.Reader
    if body != nil {
        data, _ := json.Marshal(body)
        bodyReader = bytes.NewReader(data)
    }
    req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
    if err != nil { return fmt.Errorf("new request: %w", err) }
    for k, v := range c.headers { req.Header.Set(k, v) }

    resp, err := c.httpClient.Do(req)
    if err != nil { return fmt.Errorf("do: %w", err) }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return fmt.Errorf("http status %d", resp.StatusCode)
    }
    if result != nil {
        return json.NewDecoder(resp.Body).Decode(result)
    }
    return nil
}
```

**测试**：`pkg/httpclient/client_test.go`

```go
package httpclient

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestGet(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "GET" { t.Errorf("expected GET, got %s", r.Method) }
        json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
    }))
    defer srv.Close()

    client := New(srv.URL)
    var result map[string]string
    err := client.Get(context.Background(), "/test", &result)
    if err != nil { t.Fatalf("get: %v", err) }
    if result["status"] != "ok" { t.Errorf("unexpected: %v", result) }
}

func TestPost(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" { t.Errorf("expected POST, got %s", r.Method) }
        var body map[string]string
        json.NewDecoder(r.Body).Decode(&body)
        json.NewEncoder(w).Encode(body)
    }))
    defer srv.Close()

    client := New(srv.URL)
    var result map[string]string
    err := client.Post(context.Background(), "/test", map[string]string{"a": "1"}, &result)
    if err != nil { t.Fatalf("post: %v", err) }
    if result["a"] != "1" { t.Errorf("unexpected: %v", result) }
}
```

### 1.4 共享包 — 运行测试

```bash
cd pkg/jwt       && go mod init github.com/wtb-ordering/pkg/jwt && go mod tidy && go test -v -cover ./...
cd pkg/response  && go mod init github.com/wtb-ordering/pkg/response && go mod tidy && go test -v -cover ./...
cd pkg/httpclient && go mod init github.com/wtb-ordering/pkg/httpclient && go mod tidy && go test -v -cover ./...
```

### 检查点 1

- [ ] `pkg/jwt` 测试通过，覆盖率 > 80%
- [ ] `pkg/response` 测试通过，覆盖率 > 80%
- [ ] `pkg/httpclient` 测试通过，覆盖率 > 80%
- [ ] 测试报告已生成

---

## 2. 内部微信工具包 (internal/wechat/)

> 封装微信登录、支付、商家券、订阅消息 API。其他服务通过此包调用微信。
> 开发阶段使用 Mock + 环境变量配置。

### 2.1 目录结构

```
internal/wechat/
├── go.mod
├── wechat.go          # 客户端初始化
├── login.go           # code2session
├── pay.go             # 统一下单
├── coupon.go          # 商家券
├── notify.go          # 订阅消息
├── config.go          # 配置结构
├── wechat_test.go
└── README.md
```

### 2.2 核心代码

**`internal/wechat/config.go`**

```go
package wechat

type Config struct {
    AppID     string
    AppSecret string
    MchID     string // 商户号
    APIv3Key  string // APIv3 密钥
}
```

**`internal/wechat/wechat.go`**

```go
package wechat

import "net/http"

type Client struct {
    config     Config
    httpClient *http.Client
}

func NewClient(cfg Config) *Client {
    return &Client{config: cfg, httpClient: &http.Client{}}
}
```

**`internal/wechat/login.go`**

```go
package wechat

import (
    "encoding/json"
    "fmt"
    "net/url"
)

type SessionResult struct {
    OpenID     string `json:"openid"`
    SessionKey string `json:"session_key"`
    UnionID    string `json:"unionid"`
    ErrCode    int    `json:"errcode"`
    ErrMsg     string `json:"errmsg"`
}

func (c *Client) Code2Session(code string) (*SessionResult, error) {
    u := fmt.Sprintf(
        "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
        url.QueryEscape(c.config.AppID),
        url.QueryEscape(c.config.AppSecret),
        url.QueryEscape(code))

    resp, err := c.httpClient.Get(u)
    if err != nil { return nil, fmt.Errorf("wx code2session: %w", err) }
    defer resp.Body.Close()

    var result SessionResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    if result.ErrCode != 0 {
        return nil, fmt.Errorf("wx error %d: %s", result.ErrCode, result.ErrMsg)
    }
    return &result, nil
}
```

**`internal/wechat/pay.go`** — 统一下单 + 回调验证（开发阶段 Mock）

**`internal/wechat/coupon.go`** — 商家券发放/核销/查询

**`internal/wechat/notify.go`** — 订阅消息发送

### 2.3 测试

开发阶段 Mock 微信 API。创建测试时使用 `httptest.NewServer` 模拟微信响应。

### 检查点 2

- [ ] wechat 包编译通过
- [ ] Mock 测试覆盖 login/coupon/notify 三个模块
- [ ] 测试报告已生成

---

## 从第 3 章开始 — 微服务开发模板

每个微服务遵循统一结构：

```
services/{name}/
├── go.mod
├── main.go                  # 入口
├── config/
│   └── config.go            # 配置结构 + 加载
├── model/
│   └── *.go                 # GORM 模型（PG 驱动）
├── repository/
│   └── *.go                 # 数据访问层
├── service/
│   └── *.go                 # 业务逻辑层
├── handler/
│   └── *.go                 # HTTP handler
├── router/
│   └── router.go            # 路由注册
├── migration/
│   └── migrate.go           # AutoMigrate
├── *.go_test                # 单元测试（与源文件同目录）
└── TEST_REPORT.md           # 测试报告（完成后生成）
```

### GORM + PostgreSQL 使用要点

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// 连接
dsn := "host=/tmp user=admin dbname=wtb_user sslmode=disable TimeZone=Asia/Shanghai"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// AutoMigrate 建表
db.AutoMigrate(&User{}, &RechargeRecord{}, ...)
```

PostgreSQL 注意事项：
- GORM 使用 `gorm.io/driver/postgres`（不是 mysql driver）
- `AUTO_INCREMENT` 用 GORM 默认的 `SERIAL`/`BIGSERIAL`
- `ON UPDATE CURRENT_TIMESTAMP` 在 PG 中需手动处理（或在 GORM hook 中设置 `updated_at`）
- 不要在 DDL 中写 `ENGINE=InnoDB`、`DEFAULT CHARSET=utf8mb4` 等 MySQL 专属语法
- 所有建表交给 GORM AutoMigrate，以下 DDL 仅供阅读参考

---

## 3. 用户服务 (user-service)

> 端口 8081 | 数据库 wtb_user | 依赖：pkg/jwt, pkg/response, internal/wechat

### 3.1 任务清单

- [ ] 创建 Go module 和目录结构
- [ ] 定义 GORM 模型（user, recharge_record, balance_log, consumption_record, pet_profile）
- [ ] 实现 repository 层（CRUD + 查询）
- [ ] 实现 service 层（微信登录、会员等级判定、充值、余额扣款）
- [ ] 实现 handler 层（全部 API）
- [ ] 实现路由注册 + main.go
- [ ] 编写单元测试（每个 handler 至少一个测试用例）
- [ ] 编写集成测试（测试数据库操作）
- [ ] 生成测试报告
- [ ] Git commit

### 3.2 数据模型（PostgreSQL — 参考 DDL）

```sql
-- 由 GORM AutoMigrate 自动创建，以下是等效的 PostgreSQL DDL 参考

CREATE TABLE users (
    id               BIGSERIAL PRIMARY KEY,
    openid           VARCHAR(64)  NOT NULL UNIQUE,
    unionid          VARCHAR(64)  DEFAULT '',
    nickname         VARCHAR(100) DEFAULT '',
    avatar_url       VARCHAR(500) DEFAULT '',
    phone            VARCHAR(20)  DEFAULT '',
    member_level     SMALLINT     DEFAULT 0,      -- 0普通 1会员 2充值
    balance          INTEGER      DEFAULT 0,      -- 余额（分）
    total_consumption INTEGER     DEFAULT 0,      -- 累计消费（分）
    total_orders     INTEGER      DEFAULT 0,
    created_at       TIMESTAMP    DEFAULT NOW(),
    updated_at       TIMESTAMP    DEFAULT NOW()
);
CREATE INDEX idx_users_openid ON users(openid);

CREATE TABLE recharge_records (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT       NOT NULL,
    amount        INTEGER      NOT NULL,          -- 充值金额（分）
    gifted_amount INTEGER      DEFAULT 0,         -- 赠送金额（分）
    channel       VARCHAR(20)  DEFAULT 'wxpay',
    status        VARCHAR(20)  DEFAULT 'pending',
    created_at    TIMESTAMP    DEFAULT NOW()
);
CREATE INDEX idx_recharge_user ON recharge_records(user_id);

CREATE TABLE balance_logs (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT       NOT NULL,
    type       VARCHAR(20)  NOT NULL,             -- recharge/deduct/refund
    amount     INTEGER      NOT NULL,
    order_no   VARCHAR(64)  DEFAULT '',
    remark     VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP    DEFAULT NOW()
);
CREATE INDEX idx_balance_log_user ON balance_logs(user_id);

CREATE TABLE consumption_records (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT   NOT NULL,
    order_id   BIGINT   NOT NULL,
    amount     INTEGER  NOT NULL,
    dish_count INTEGER  DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_consumption_user ON consumption_records(user_id);

CREATE TABLE pet_profiles (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT       NOT NULL,
    name       VARCHAR(50)  NOT NULL,
    breed      VARCHAR(50)  DEFAULT '',
    weight     NUMERIC(5,2) DEFAULT 0,
    birthday   DATE         DEFAULT NULL,
    created_at TIMESTAMP    DEFAULT NOW()
);
CREATE INDEX idx_pet_user ON pet_profiles(user_id);
```

### 3.3 GORM Model 示例（Go 代码）

```go
package model

import "time"

type User struct {
    ID               uint      `gorm:"primaryKey"`
    OpenID           string    `gorm:"uniqueIndex;size:64"`
    UnionID          string    `gorm:"size:64"`
    Nickname         string    `gorm:"size:100"`
    AvatarURL        string    `gorm:"size:500"`
    Phone            string    `gorm:"size:20"`
    MemberLevel      int16     `gorm:"default:0"`
    Balance          int       `gorm:"default:0"`
    TotalConsumption int       `gorm:"default:0"`
    TotalOrders      int       `gorm:"default:0"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
```

### 3.4 API 清单

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/user/wx-login | 微信 code 换 JWT |
| GET | /api/user/profile | 用户信息+等级+余额 |
| GET | /api/user/consumption | 消费记录 |
| GET | /api/user/consumption/summary | 消费汇总 |
| POST | /api/user/recharge | 充值 |
| POST | /api/user/balance/deduct | 余额扣款（内部） |
| POST | /api/user/balance/refund | 余额退款（内部） |
| GET | /api/user/pets | 宠物列表 |
| POST | /api/user/pets | 添加宠物 |
| GET | /api/user/internal/:id | 内部 RPC |

### 3.5 关键业务逻辑

```
会员等级判定：
- 普通客户：注册默认 (member_level=0)
- 会员客户：total_consumption >= 指定金额 或 total_orders >= 指定次数 (member_level=1)
- 充值客户：有过充值记录即升级 (member_level=2)

积分倍率（供积分服务调用）：
- 普通 ×1、会员 ×1.5、充值 ×2
```

### 3.6 测试要求

1. **单元测试**：Mock wechat client，测试 wx-login handler
2. **单元测试**：测试 profile handler（含不同会员等级）
3. **单元测试**：测试 recharge handler（充值成功/失败）
4. **单元测试**：测试 balance deduct/refund（内部接口）
5. **集成测试**：使用真实 PostgreSQL 测试数据库 CRUD
6. 覆盖率目标：> 70%

### 3.7 测试报告

完成后填写 `services/user/TEST_REPORT.md`（参考附录A模板）。

---

## 4. 座位服务 (seat-service)

> 端口 8082 | 数据库 wtb_seat | 依赖：pkg/response

### 4.1 任务清单

- [ ] 创建 Go module 和目录结构
- [ ] 定义模型（area, seat, seat_status_log）
- [ ] repository 层
- [ ] service 层（状态流转、二维码生成）
- [ ] handler 层
- [ ] 路由 + main.go
- [ ] 测试
- [ ] 测试报告
- [ ] Git commit

### 4.2 数据模型（PostgreSQL 参考 DDL）

```sql
CREATE TABLE areas (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(50) NOT NULL,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE seats (
    id         BIGSERIAL PRIMARY KEY,
    area_id    BIGINT       NOT NULL,
    name       VARCHAR(50)  NOT NULL,
    type       VARCHAR(20)  DEFAULT 'normal',    -- normal/booth/outdoor
    capacity   INTEGER      DEFAULT 4,
    qrcode_url VARCHAR(500) DEFAULT '',
    status     VARCHAR(20)  DEFAULT 'available', -- available/occupied/reserved/cleaning
    created_at TIMESTAMP    DEFAULT NOW(),
    updated_at TIMESTAMP    DEFAULT NOW()
);
CREATE INDEX idx_seats_area ON seats(area_id);

CREATE TABLE seat_status_logs (
    id         BIGSERIAL PRIMARY KEY,
    seat_id    BIGINT      NOT NULL,
    old_status VARCHAR(20),
    new_status VARCHAR(20),
    order_id   BIGINT      DEFAULT NULL,
    changed_at TIMESTAMP   DEFAULT NOW()
);
CREATE INDEX idx_seat_log_seat ON seat_status_logs(seat_id);
```

### 4.3 API 清单

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/seat/areas | 区域列表 |
| POST | /api/seat/areas | 新增区域 |
| GET | /api/seat/list?area_id= | 座位列表 |
| GET | /api/seat/:id | 座位详情 |
| POST | /api/seat/qrcode/batch | 批量生成二维码 |
| GET | /api/seat/scan?code=xxx | 扫码解析 |
| GET | /api/seat/internal/:id | 内部 RPC |

### 4.4 测试要求

- [ ] 区域 CRUD 测试
- [ ] 座位 CRUD + 状态流转测试
- [ ] 二维码生成与解析测试
- 覆盖率目标：> 70%

---

## 5. 菜品服务 (menu-service)

> 端口 8083 | 数据库 wtb_menu | 依赖：pkg/response

### 5.1 任务清单

- [ ] 创建 Go module 和目录结构
- [ ] 定义模型（category, dish, dish_price, dish_stock）
- [ ] repository 层
- [ ] service 层（搜索、分类树）
- [ ] handler 层
- [ ] 路由 + main.go
- [ ] 测试
- [ ] 测试报告
- [ ] Git commit

### 5.2 数据模型（PostgreSQL 参考 DDL）

```sql
CREATE TABLE categories (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(50) NOT NULL,
    parent_id  BIGINT   DEFAULT 0,
    sort_order INTEGER  DEFAULT 0,
    status     SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE dishes (
    id          BIGSERIAL PRIMARY KEY,
    category_id BIGINT       NOT NULL,
    name        VARCHAR(100) NOT NULL,
    subtitle    VARCHAR(200) DEFAULT '',
    description TEXT,
    images      TEXT,                              -- JSON array
    tags        VARCHAR(200) DEFAULT '',           -- 逗号分隔
    status      SMALLINT     DEFAULT 1,
    created_at  TIMESTAMP    DEFAULT NOW(),
    updated_at  TIMESTAMP    DEFAULT NOW()
);
CREATE INDEX idx_dishes_category ON dishes(category_id);
-- PostgreSQL 全文搜索使用 GIN + tsvector
ALTER TABLE dishes ADD COLUMN search_vector tsvector
    GENERATED ALWAYS AS (to_tsvector('simple', name || ' ' || COALESCE(description, ''))) STORED;
CREATE INDEX idx_dishes_search ON dishes USING GIN(search_vector);

CREATE TABLE dish_prices (
    id         BIGSERIAL PRIMARY KEY,
    dish_id    BIGINT      NOT NULL,
    price_type VARCHAR(20) NOT NULL,              -- normal/member/time_slot
    price      INTEGER     NOT NULL,              -- 价格（分）
    start_time TIMESTAMP   DEFAULT NULL,
    end_time   TIMESTAMP   DEFAULT NULL
);
CREATE INDEX idx_dish_prices_dish ON dish_prices(dish_id);

CREATE TABLE dish_stocks (
    id          BIGSERIAL PRIMARY KEY,
    dish_id     BIGINT  NOT NULL,
    daily_limit INTEGER DEFAULT -1,               -- -1 无限
    sold_count  INTEGER DEFAULT 0,
    date        DATE    NOT NULL,
    UNIQUE (dish_id, date)
);
```

### 5.3 API 清单

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/menu/categories | 分类树 |
| GET | /api/menu/dishes?category_id=&page=&pageSize= | 菜品列表 |
| GET | /api/menu/dish/:id | 菜品详情 |
| GET | /api/menu/search?q= | 搜索（使用 tsvector） |
| POST | /api/menu/dishes/batch | 批量查询（内部） |
| POST | /api/menu/admin/category | 新增分类 |
| PUT | /api/menu/admin/category/:id | 更新分类 |
| POST | /api/menu/admin/dish | 新增菜品 |
| PUT | /api/menu/admin/dish/:id | 更新菜品 |
| DELETE | /api/menu/admin/dish/:id | 删除菜品 |
| POST | /api/menu/admin/stock | 设置库存 |

### 5.4 测试要求

- [ ] 分类树构建测试（含多级分类）
- [ ] 菜品 CRUD 测试
- [ ] 全文搜索测试（tsvector 查询）
- [ ] 库存扣减测试
- 覆盖率目标：> 70%

---

## 6. 营销定价服务 (pricing-service)

> 端口 8088 | 数据库 wtb_pricing | 依赖：pkg/response, pkg/httpclient（调用 menu-service）

### 6.1 任务清单

- [ ] 创建 Go module 和目录结构
- [ ] 定义模型（price_rule, promotion, combo）
- [ ] repository 层
- [ ] service 层（价格计算引擎、活动叠加逻辑）
- [ ] handler 层
- [ ] 路由 + main.go
- [ ] 测试
- [ ] 测试报告
- [ ] Git commit

### 6.2 数据模型（PostgreSQL 参考 DDL）

```sql
CREATE TABLE price_rules (
    id         BIGSERIAL PRIMARY KEY,
    dish_id    BIGINT      NOT NULL,
    rule_type  VARCHAR(20) NOT NULL,              -- normal/member/time_slot
    price      INTEGER     NOT NULL,
    start_time TIMESTAMP   DEFAULT NULL,
    end_time   TIMESTAMP   DEFAULT NULL,
    status     SMALLINT    DEFAULT 1
);
CREATE INDEX idx_price_rules_dish ON price_rules(dish_id);

CREATE TABLE promotions (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(30)  NOT NULL,            -- full_reduction/discount/combo
    config_json TEXT         NOT NULL,            -- 规则配置JSON
    start_time  TIMESTAMP    NOT NULL,
    end_time    TIMESTAMP    NOT NULL,
    status      SMALLINT     DEFAULT 1
);
CREATE INDEX idx_promotions_time ON promotions(start_time, end_time);

CREATE TABLE combos (
    id        BIGSERIAL PRIMARY KEY,
    name      VARCHAR(100) NOT NULL,
    price     INTEGER      NOT NULL,
    dish_list TEXT         NOT NULL,              -- JSON: [{dish_id, quantity}]
    status    SMALLINT     DEFAULT 1
);
```

### 6.3 API 清单

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/pricing/calculate | 计算订单价格 |
| GET | /api/pricing/dish/:id | 菜品当前价 |
| GET | /api/pricing/promotions | 营销活动列表 |
| POST | /api/pricing/admin/rule | 新增价格规则 |
| POST | /api/pricing/admin/promotion | 新增活动 |
| POST | /api/pricing/admin/combo | 新增套餐 |

### 6.4 核心算法

```go
// CalculateOrderPrice 计算订单价格
// 输入：{ user_level, dishes: [{dish_id, quantity}] }
// 输出：{ items: [{dish_id, unit_price, subtotal}], 
//          total, discount, final_amount, 
//          applied_promotions: [...] }
//
// 逻辑：
// 1. 批量查询菜品基础价（调用 menu-service 内部接口）
// 2. 根据 user_level 确定单价（会员价/时段价）
// 3. 匹配满减活动
// 4. 按优先级叠加优惠
// 5. 返回详细计算明细
```

### 6.5 测试要求

- [ ] 单菜品定价测试（普通/会员/时段）
- [ ] 满减活动测试
- [ ] 套餐价格测试
- [ ] 优惠叠加逻辑测试
- 覆盖率目标：> 80%（定价是核心逻辑，需要高覆盖）

---

## 7. 订单服务 (order-service)

> 端口 8084 | 数据库 wtb_order | Redis 购物车
> 依赖：pkg/response, pkg/httpclient（调用 user/seat/menu/pricing/payment）

### 7.1 任务清单

- [ ] 创建 Go module 和目录结构
- [ ] 定义模型（order, order_item, order_status_log）
- [ ] repository 层
- [ ] Redis 购物车操作封装
- [ ] service 层（下单主流程状态机）
- [ ] handler 层
- [ ] 路由 + main.go
- [ ] 测试
- [ ] 测试报告
- [ ] Git commit

### 7.2 数据模型（PostgreSQL 参考 DDL）

```sql
CREATE TABLE orders (
    id              BIGSERIAL PRIMARY KEY,
    order_no        VARCHAR(32) NOT NULL UNIQUE,
    seat_id         BIGINT      NOT NULL,
    user_id         BIGINT      NOT NULL,
    status          VARCHAR(20) DEFAULT 'pending', -- pending/confirmed/cooking/served/paid/completed/cancelled/refunded
    total_amount    INTEGER     NOT NULL,
    discount_amount INTEGER     DEFAULT 0,
    pay_amount      INTEGER     NOT NULL,
    remark          VARCHAR(500) DEFAULT '',
    created_at      TIMESTAMP   DEFAULT NOW(),
    updated_at      TIMESTAMP   DEFAULT NOW()
);
CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_seat ON orders(seat_id);
CREATE INDEX idx_orders_no ON orders(order_no);

CREATE TABLE order_items (
    id         BIGSERIAL PRIMARY KEY,
    order_id   BIGINT       NOT NULL,
    dish_id    BIGINT       NOT NULL,
    dish_name  VARCHAR(100) NOT NULL,
    quantity   INTEGER      NOT NULL,
    unit_price INTEGER      NOT NULL
);
CREATE INDEX idx_order_items_order ON order_items(order_id);

CREATE TABLE order_status_logs (
    id          BIGSERIAL PRIMARY KEY,
    order_id    BIGINT      NOT NULL,
    from_status VARCHAR(20),
    to_status   VARCHAR(20),
    operator    VARCHAR(50) DEFAULT '',
    created_at  TIMESTAMP   DEFAULT NOW()
);
CREATE INDEX idx_order_status_log_order ON order_status_logs(order_id);
```

**Redis 购物车结构**：
```
Key:   cart:{seat_id}
Value: Hash: { dish_id: JSON{user_id, dish_name, quantity, unit_price, remark} }
TTL:   30 分钟（座位空闲后自动清理）
```

### 7.3 API 清单

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/order/cart/add | 加入购物车 |
| GET | /api/order/cart/list?seat_id= | 购物车列表 |
| PUT | /api/order/cart/update | 修改数量/备注 |
| DELETE | /api/order/cart/remove | 移除商品 |
| POST | /api/order/create | 提交订单 |
| POST | /api/order/:id/pay | 发起支付 |
| GET | /api/order/:id/status | 订单状态 |
| GET | /api/order/list | 我的订单 |
| POST | /api/order/internal/notify | 支付回调（内部） |
| GET | /api/order/admin/list | 后台订单 |
| PUT | /api/order/admin/status | 改状态 |
| POST | /api/order/admin/refund | 退单 |

### 7.4 下单主流程（service 层实现）

```
CreateOrder(seatID, userID) → Order:
1. 校验座位状态（调用 seat-service）
2. 从 Redis 读取购物车
3. 校验库存（调用 menu-service）
4. 生成订单号 + 订单记录
5. 调用 pricing-service 计算价格
6. 锁定库存 15 分钟（Redis 分布式锁）
7. 返回订单详情（含支付参数由 payment-service 生成）
```

### 7.5 测试要求

- [ ] 购物车 Redis 操作测试
- [ ] 下单流程单元测试（Mock 所有下游服务）
- [ ] 订单状态机测试（所有合法/非法流转）
- [ ] 多人同时操作购物车测试（并发安全）
- 覆盖率目标：> 75%

---

## 8. 支付服务 (payment-service)

> 端口 8085 | 数据库 wtb_payment | Redis 防重
> 依赖：pkg/response, pkg/httpclient, internal/wechat

### 8.1 任务清单

- [ ] 创建 Go module 和目录结构
- [ ] 定义模型（payment_order, payment_record, refund_record, recharge_order）
- [ ] repository 层
- [ ] service 层（微信支付/余额支付/退款/充值）
- [ ] handler 层
- [ ] 路由 + main.go
- [ ] 测试
- [ ] 测试报告
- [ ] Git commit

### 8.2 数据模型（PostgreSQL 参考 DDL）

```sql
CREATE TABLE payment_orders (
    id            BIGSERIAL PRIMARY KEY,
    order_no      VARCHAR(32) NOT NULL,
    out_trade_no  VARCHAR(32) NOT NULL UNIQUE,
    user_id       BIGINT      NOT NULL,
    amount        INTEGER     NOT NULL,           -- 支付金额（分）
    channel       VARCHAR(20) NOT NULL,           -- wxpay/balance
    status        VARCHAR(20) DEFAULT 'pending',  -- pending/paid/closed/refunded
    wx_prepay_id  VARCHAR(64) DEFAULT '',
    created_at    TIMESTAMP   DEFAULT NOW()
);
CREATE INDEX idx_payment_orders_no ON payment_orders(order_no);

CREATE TABLE payment_records (
    id               BIGSERIAL PRIMARY KEY,
    payment_order_id BIGINT      NOT NULL,
    channel          VARCHAR(20) NOT NULL,
    amount           INTEGER     NOT NULL,
    transaction_id   VARCHAR(64) DEFAULT '',
    paid_at          TIMESTAMP   DEFAULT NULL
);
CREATE INDEX idx_payment_records_order ON payment_records(payment_order_id);

CREATE TABLE refund_records (
    id               BIGSERIAL PRIMARY KEY,
    payment_order_id BIGINT       NOT NULL,
    refund_no        VARCHAR(32)  NOT NULL UNIQUE,
    amount           INTEGER      NOT NULL,
    reason           VARCHAR(200) DEFAULT '',
    status           VARCHAR(20)  DEFAULT 'pending',
    created_at       TIMESTAMP    DEFAULT NOW()
);
CREATE INDEX idx_refund_records_order ON refund_records(payment_order_id);

CREATE TABLE recharge_orders (
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT       NOT NULL,
    amount         INTEGER      NOT NULL,
    gifted_amount  INTEGER      DEFAULT 0,
    discount_rate  NUMERIC(3,2) DEFAULT 1.00,
    final_amount   INTEGER      NOT NULL,
    status         VARCHAR(20)  DEFAULT 'pending',
    created_at     TIMESTAMP    DEFAULT NOW()
);
CREATE INDEX idx_recharge_orders_user ON recharge_orders(user_id);
```

### 8.3 API 清单

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/pay/create | 创建支付单 |
| POST | /api/pay/wx/prepay | 微信统一下单 |
| POST | /api/pay/balance | 余额支付 |
| POST | /api/pay/recharge | 充值 |
| POST | /api/pay/callback/wx | 微信回调 |
| POST | /api/pay/refund | 退款 |
| GET | /api/pay/query/:outTradeNo | 查询 |

### 8.4 测试要求

- [ ] 支付单创建测试
- [ ] 余额支付测试（含余额不足场景）
- [ ] 微信支付 Mock 测试
- [ ] 退款测试
- [ ] 支付回调幂等性测试
- 覆盖率目标：> 75%

---

## 9. 积分服务 (points-service)

> 端口 8086 | 数据库 wtb_points | Redis 积分排名缓存
> 依赖：pkg/response, pkg/httpclient

### 9.1 任务清单

- [ ] 创建 Go module 和目录结构
- [ ] 定义模型（points_rule, user_points, points_log, exchange_goods, exchange_order）
- [ ] repository 层
- [ ] service 层（积分发放倍率计算、兑换扣减）
- [ ] handler 层
- [ ] 路由 + main.go
- [ ] 测试
- [ ] 测试报告
- [ ] Git commit

### 9.2 数据模型（PostgreSQL 参考 DDL）

```sql
CREATE TABLE points_rules (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(50) NOT NULL,
    type        VARCHAR(30) NOT NULL,             -- consumption/recharge/sign_in
    config_json TEXT        NOT NULL,
    status      SMALLINT    DEFAULT 1
);

CREATE TABLE user_points (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT   NOT NULL UNIQUE,
    total_points  INTEGER  DEFAULT 0,
    used_points   INTEGER  DEFAULT 0,
    frozen_points INTEGER  DEFAULT 0,
    updated_at    TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_user_points_user ON user_points(user_id);

CREATE TABLE points_logs (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT       NOT NULL,
    type       VARCHAR(20)  NOT NULL,             -- gain/exchange/expire/adjust
    points     INTEGER      NOT NULL,
    source_id  VARCHAR(64)  DEFAULT '',           -- 关联订单号/兑换单号
    remark     VARCHAR(200) DEFAULT '',
    created_at TIMESTAMP    DEFAULT NOW()
);
CREATE INDEX idx_points_logs_user_time ON points_logs(user_id, created_at);

CREATE TABLE exchange_goods (
    id           BIGSERIAL PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    image        VARCHAR(500) DEFAULT '',
    points_price INTEGER      NOT NULL,
    stock        INTEGER      DEFAULT 0,
    type         VARCHAR(20)  DEFAULT 'physical',
    status       SMALLINT     DEFAULT 1
);

CREATE TABLE exchange_orders (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    goods_id    BIGINT      NOT NULL,
    points_cost INTEGER     NOT NULL,
    status      VARCHAR(20) DEFAULT 'pending',
    created_at  TIMESTAMP   DEFAULT NOW()
);
CREATE INDEX idx_exchange_orders_user ON exchange_orders(user_id);
```

### 9.3 API 清单

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/points/account | 积分余额 |
| GET | /api/points/logs | 积分流水 |
| GET | /api/points/goods | 兑换商品列表 |
| POST | /api/points/exchange | 积分兑换 |
| POST | /api/points/internal/grant | 发放积分（内部） |
| POST | /api/points/admin/rule | 配置积分规则 |
| POST | /api/points/admin/goods | 新增兑换商品 |

### 9.4 积分发放核心逻辑

```go
// GrantPoints(userID, level, amount, sourceType) error
// points = amount * multiplier
// multiplier:
//   level 0 (普通) → ×1
//   level 1 (会员) → ×1.5
//   level 2 (充值) → ×2
```

### 9.5 测试要求

- [ ] 积分发放测试（所有倍率）
- [ ] 积分兑换测试（含库存不足）
- [ ] 积分流水查询测试
- 覆盖率目标：> 70%

---

## 10. 活动服务 (activity-service)

> 端口 8087 | 数据库 wtb_activity
> 依赖：pkg/response, internal/wechat（订阅消息）

### 10.1 任务清单

- [ ] 创建 Go module 和目录结构
- [ ] 定义模型（announcement, activity, activity_registration）
- [ ] repository 层
- [ ] service 层
- [ ] handler 层
- [ ] 路由 + main.go
- [ ] 测试
- [ ] 测试报告
- [ ] Git commit

### 10.2 数据模型（PostgreSQL 参考 DDL）

```sql
CREATE TABLE announcements (
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(100) NOT NULL,
    content     TEXT,
    type        VARCHAR(20)  DEFAULT 'text',
    image       VARCHAR(500) DEFAULT '',
    link_type   VARCHAR(20)  DEFAULT '',
    link_target VARCHAR(500) DEFAULT '',
    sort_order  INTEGER      DEFAULT 0,
    start_time  TIMESTAMP    NOT NULL,
    end_time    TIMESTAMP    NOT NULL,
    status      SMALLINT     DEFAULT 1
);
CREATE INDEX idx_announcements_time ON announcements(start_time, end_time);

CREATE TABLE activities (
    id                   BIGSERIAL PRIMARY KEY,
    title                VARCHAR(100) NOT NULL,
    description          TEXT,
    image                VARCHAR(500) DEFAULT '',
    max_participants     INTEGER      DEFAULT -1,  -- -1 不限
    current_participants INTEGER      DEFAULT 0,
    event_time           TIMESTAMP    DEFAULT NULL,
    location             VARCHAR(200) DEFAULT '',
    status               VARCHAR(20)  DEFAULT 'draft', -- draft/published/cancelled/ended
    created_at           TIMESTAMP    DEFAULT NOW()
);

CREATE TABLE activity_registrations (
    id          BIGSERIAL PRIMARY KEY,
    activity_id BIGINT       NOT NULL,
    user_id     BIGINT       NOT NULL,
    name        VARCHAR(50)  DEFAULT '',
    phone       VARCHAR(20)  DEFAULT '',
    remark      VARCHAR(200) DEFAULT '',
    status      VARCHAR(20)  DEFAULT 'registered',  -- registered/cancelled
    created_at  TIMESTAMP    DEFAULT NOW(),
    UNIQUE (user_id, activity_id)
);
```

### 10.3 API 清单

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/activity/announcements | 生效公告 |
| GET | /api/activity/recharge-discount | 充值折扣配置 |
| GET | /api/activity/list | 活动列表 |
| POST | /api/activity/:id/register | 报名 |
| GET | /api/activity/my-registrations | 我的报名 |
| PUT | /api/activity/:id/cancel | 取消报名 |
| GET | /api/activity/internal/discount | 内部查折扣 |
| POST | /api/activity/admin/announcement | 新增公告 |
| POST | /api/activity/admin/activity | 新增活动 |

### 10.4 测试要求

- [ ] 公告 CRUD 测试
- [ ] 活动报名测试（含名额限制）
- [ ] 取消报名测试
- 覆盖率目标：> 70%

---

## 11. 数据统计服务 (analytics-service)

> 端口 8089 | 数据库 wtb_analytics（存汇总数据）+ 跨库查询
> 依赖：pkg/response

### 11.1 任务清单

- [ ] 创建 Go module 和目录结构
- [ ] 定义聚合查询逻辑（使用 PostgreSQL 窗口函数）
- [ ] handler 层
- [ ] 路由 + main.go
- [ ] 测试
- [ ] 测试报告
- [ ] Git commit

### 11.2 API 清单

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/analytics/dashboard | 今日仪表盘 |
| GET | /api/analytics/revenue | 营收统计 |
| GET | /api/analytics/dishes | 菜品销量排行 |
| GET | /api/analytics/members | 会员分布 |
| GET | /api/analytics/points | 积分统计 |
| GET | /api/analytics/coupons | 卡券统计 |
| GET | /api/analytics/activities | 活动统计 |
| POST | /api/analytics/export | Excel 导出 |

### 11.3 实现说明

- PostgreSQL 18.3 支持窗口函数（`ROW_NUMBER()`, `RANK()`, `SUM() OVER` 等）
- 导出使用 `github.com/xuri/excelize/v2` 生成 Excel
- 跨库查询使用 `dblink` 或在代码层聚合多个数据库的查询结果

### 11.4 测试要求

- [ ] 各统计接口响应格式测试
- 覆盖率目标：> 60%

---

## 12. 后台聚合服务 (admin-bff)

> 端口 8090 | 无数据库（聚合层）
> 依赖：pkg/response, pkg/jwt, pkg/httpclient（调用所有下游服务）

### 12.1 任务清单

- [ ] 创建 Go module 和目录结构
- [ ] 鉴权中间件（JWT + 角色校验）
- [ ] 代理转发所有 admin 接口
- [ ] 聚合查询（如：座位占用率 + 当前订单）
- [ ] 操作审计日志
- [ ] 路由 + main.go
- [ ] 测试
- [ ] 测试报告
- [ ] Git commit

### 12.2 核心功能

- 接收后台管理前端的请求
- 校验管理员权限（JWT 中包含 role）
- 转发到对应的下游微服务
- 提供聚合接口（一次请求查多个服务）

### 12.3 测试要求

- [ ] 鉴权中间件测试
- [ ] 代理转发测试
- 覆盖率目标：> 60%

---

## 13. API 网关 (gateway)

> 端口 8080 | 无数据库
> 依赖：pkg/jwt（验证）

### 13.1 功能

- 路由分发（根据路径前缀转发到对应微服务）
- JWT 鉴权（公开接口白名单：/api/user/wx-login, /api/pay/callback/wx）
- 限流（令牌桶，每秒 200 请求）
- TraceID 注入
- CORS 处理
- 请求日志

### 13.2 任务清单

- [ ] 创建 Go module
- [ ] 路由表配置
- [ ] JWT 鉴权中间件
- [ ] 反向代理
- [ ] 限流中间件
- [ ] 日志中间件
- [ ] main.go
- [ ] 测试
- [ ] 测试报告
- [ ] Git commit

### 13.3 路由配置示例

```go
var routes = map[string]string{
    "/api/user/":      "http://localhost:8081",
    "/api/seat/":      "http://localhost:8082",
    "/api/menu/":      "http://localhost:8083",
    "/api/order/":     "http://localhost:8084",
    "/api/pay/":       "http://localhost:8085",
    "/api/points/":    "http://localhost:8086",
    "/api/activity/":  "http://localhost:8087",
    "/api/pricing/":   "http://localhost:8088",
    "/api/analytics/": "http://localhost:8089",
    "/api/admin/":     "http://localhost:8090",
}

var publicPaths = []string{
    "/api/user/wx-login",
    "/api/pay/callback/wx",
}
```

### 13.4 测试要求

- [ ] 路由分发测试
- [ ] JWT 鉴权测试（通过/拒绝）
- [ ] 限流测试
- 覆盖率目标：> 70%

---

## 14. Docker Compose 集成

> 注意：本地 Docker daemon 需手动启动（Docker Desktop）
> 日常开发直接使用本地 Homebrew PostgreSQL + Redis，Docker 仅用于 CI/CD 或新成员快速搭建

### 14.1 docker-compose.yml

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:18-alpine
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: dev123456
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  pgdata:
```

### 14.2 init.sql（自动建库）

```sql
CREATE DATABASE wtb_user      OWNER admin ENCODING 'UTF8';
CREATE DATABASE wtb_seat      OWNER admin ENCODING 'UTF8';
CREATE DATABASE wtb_menu      OWNER admin ENCODING 'UTF8';
CREATE DATABASE wtb_order     OWNER admin ENCODING 'UTF8';
CREATE DATABASE wtb_payment   OWNER admin ENCODING 'UTF8';
CREATE DATABASE wtb_points    OWNER admin ENCODING 'UTF8';
CREATE DATABASE wtb_activity  OWNER admin ENCODING 'UTF8';
CREATE DATABASE wtb_pricing   OWNER admin ENCODING 'UTF8';
CREATE DATABASE wtb_analytics OWNER admin ENCODING 'UTF8';
```

### 14.3 Makefile

```makefile
.PHONY: help run-all test-all db-create

help:
	@echo "make db-create     # 创建所有数据库"
	@echo "make run-user      # 启动用户服务"
	@echo "make run-gateway   # 启动网关"
	@echo "make test-user     # 测试用户服务"
	@echo "make test-all      # 测试所有服务"

db-create:
	@for db in user seat menu order payment points activity pricing analytics; do \
		psql -h /tmp -U admin -d postgres -c "CREATE DATABASE wtb_$${db} OWNER admin ENCODING 'UTF8';" 2>/dev/null; \
		echo "created wtb_$${db}"; \
	done

run-user:
	cd services/user && go run main.go

run-gateway:
	cd gateway && go run main.go

test-user:
	cd services/user && go test ./... -v -cover

test-all:
	cd pkg/jwt && go test ./... -v -cover
	cd pkg/response && go test ./... -v -cover
	cd pkg/httpclient && go test ./... -v -cover
	@for svc in user seat menu order payment points activity pricing analytics admin; do \
		cd services/$$svc && go test ./... -v -cover && cd ../..; \
	done
	cd gateway && go test ./... -v -cover
```

---

## 15. 前端应用（概要）

> 完整开发由前端开发者负责，以下为后端同学需要了解的前端对接口要求。

### 15.1 微信小程序

| 页面 | 调用的后端 API |
|------|---------------|
| 首页 | /api/activity/announcements, /api/activity/recharge-discount |
| 扫码页 | /api/seat/scan |
| 菜单页 | /api/menu/categories, /api/menu/dishes |
| 购物车 | /api/order/cart/* |
| 订单确认 | /api/pricing/calculate, /api/pay/create |
| 支付页 | /api/pay/wx/prepay 或 /api/pay/balance |
| 订单列表 | /api/order/list, /api/order/:id/status |
| 会员中心 | /api/user/profile, /api/user/consumption |
| 积分商城 | /api/points/account, /api/points/goods, /api/points/exchange |
| 活动中心 | /api/activity/list, /api/activity/:id/register |

### 15.2 后台管理前端

- 所有管理接口通过 `/api/admin/` 前缀访问 → admin-bff 转发
- 使用 Ant Design Pro 脚手架搭建

---

## 执行顺序总结

```
Phase 0:  环境准备 (0.1 → 0.5)
Phase 1:  共享包 (pkg/jwt → pkg/response → pkg/httpclient)
Phase 2:  internal/wechat
Phase 3:  user-service      (端口 8081)
Phase 4:  seat-service      (端口 8082)
Phase 5:  menu-service      (端口 8083)
Phase 6:  pricing-service   (端口 8088) — 依赖 menu
Phase 7:  order-service     (端口 8084) — 依赖 user/seat/menu/pricing
Phase 8:  payment-service   (端口 8085) — 依赖 user/order
Phase 9:  points-service    (端口 8086) — 依赖 user
Phase 10: activity-service  (端口 8087)
Phase 11: analytics-service (端口 8089)
Phase 12: admin-bff         (端口 8090)
Phase 13: gateway           (端口 8080)
Phase 14: 前端应用
```

**依赖图**：
```
gateway ─────────────────────────────────────────┐
  │ (JWT 鉴权后路由)                              │
  ├── user    ← wechat (login)                   │
  ├── seat                                       │
  ├── menu                                       │
  ├── pricing ← menu (查菜品信息)                 │
  ├── order   ← user + seat + menu + pricing     │
  ├── payment ← user + order + wechat (pay)      │
  ├── points  ← user (查等级算倍率)               │
  ├── activity ← wechat (notify)                  │
  ├── analytics                                  │
  └── admin   → 聚合所有 (BFF)                    │
```

---

## 附录A：测试报告模板

每个微服务完成后，在服务目录下创建 `TEST_REPORT.md`：

```markdown
# {服务名称} 测试报告

- **日期**：YYYY-MM-DD
- **服务**：{service-name}
- **执行人**：Kimi

## 测试环境

| 项目 | 值 |
|------|-----|
| Go 版本 | 1.24.1 |
| PostgreSQL | 18.3 |
| Redis | 7.x (如适用) |

## 测试结果

\`\`\`
go test ./... -v -cover
\`\`\`

## 覆盖率

| 包 | 覆盖率 |
|----|--------|
| model | xx% |
| repository | xx% |
| service | xx% |
| handler | xx% |
| **总计** | **xx%** |

## 通过的测试用例

- [x] TestXxx — 描述
- [x] TestYyy — 描述

## 失败的测试用例

- [ ] TestZzz — 失败原因

## 已知问题

1. 问题描述 + 临时解决方案

## 结论

- 服务是否可进入下一阶段：是/否
- 阻塞项：无 / 有（描述）
```

---

## 附录B：每阶段验证命令

```bash
# === Phase 0 ===
go version                               # go1.24.1
psql -h /tmp -U admin -d postgres -c "\l" | grep wtb_   # 显示 9 个数据库
redis-cli ping                           # PONG

# === Phase 1: 共享包 ===
cd pkg/jwt        && go test -v -cover ./...
cd pkg/response   && go test -v -cover ./...
cd pkg/httpclient && go test -v -cover ./...

# === Phase 3-12: 微服务（以 user 为例） ===
cd services/user
go mod tidy
go build ./...                           # 编译
go test ./... -v -cover                  # 测试
go run main.go &                         # 启动
curl http://localhost:8081/api/user/profile -H "Authorization: Bearer <token>"

# === Phase 13: 网关 ===
cd gateway
go run main.go &
curl http://localhost:8080/api/user/profile -H "Authorization: Bearer <token>"
```

---

> 本计划按依赖顺序排列，每完成一个微服务并进行测试后，再进行下一个。
> 每个阶段完成后执行 `git add -A && git commit -m "feat({service}): 完成{服务名}开发 + 测试"`。
