# 汪托帮点餐系统 — 优化建议文档

> 基于 2026-05-24 代码库分析，按优先级从高到低排列。
> **已修复项标记为 ✅**，其余为待处理建议。

---

## 已修复（2026-05-24）

| 项 | 说明 |
|----|------|
| ✅ 1.1 购物车持久化 | 新建 `carts` 表（`services/order/model/cart.go`），`CartService` 改为通过 `CartRepo` 读写 PostgreSQL，服务重启不丢数据 |
| ✅ 1.2 余额扣款非原子 | `DeductBalance` 和 `RefundBalance` 使用 `DB.Transaction` + `WHERE balance >= ?` 原子 UPDATE，余额更新与日志写入在同一事务中 |
| ✅ 1.3 下单无事务 | `CreateOrder` 使用 `DB.Transaction` 包裹订单创建、订单项创建、状态日志创建 |
| ✅ 1.4 GORM AutoMigrate 生产执行 | 通过 `AUTO_MIGRATE=false` 环境变量可关闭自动迁移 |
| ✅ 2.1 8 个独立 DB 连接 | 每个 DB 连接池限制为 `MaxOpenConns=5, MaxIdleConns=2`，80 → 40 总连接数 |
| ✅ 2.2 优雅关闭 | 添加 `SIGTERM/SIGINT` 信号处理，10 秒超时优雅关闭 HTTP 服务，关闭所有 DB 连接 |
| ✅ 2.3 健康检查 | `/health` 增加数据库 Ping 检查，DB 不可用时返回 503 |
| ✅ 2.5 CORS 宽松 | 生产环境（`ENV=production`）启用 Origin 白名单；非生产环境保持宽松便于调试 |
| ✅ 4.4 管理员密码硬编码 | 改为从 `ADMIN_PASSWORD` 环境变量读取，默认 fallback `1234` |
| ✅ 5.1 代码压缩 | `order_service.go`、`order_repo.go`、`order_handler.go`、`order_item_repo.go`、`order_status_log_repo.go` 全部展开为标准 Go 格式 |
| ✅ 5.2 错误静默忽略 | `DeductBalance`/`RefundBalance` 的错误处理完善，余额日志写入失败会回滚事务 |
| ✅ 6.1 AdminListPets N+1 | 改为先收集所有 `user_id`，再 `FindByIDs` 批量查询用户，从 N 次查询降为 2 次 |
| ✅ 3.1 加减按钮无防抖 | `syncCart` 改为按 dishId 防抖（300ms trailing debounce），快速点击只发最后一次请求 |
| ✅ 3.3 登录时序竞态 | `app.js` 新增 `onLoginReady(callback)` 事件机制，`menu.js` 不再用 `setTimeout` 等登录 |
| ✅ 4.1 订单轮询 | 轮询间隔从 5s 改为 15s，增加请求去重（`loadingRef`），避免重复请求堆积 |
| ✅ 4.2 Token localStorage | 改为 `sessionStorage`，关闭标签页自动清除 |

---

## 一、数据安全与一致性（高优先级）

### 1.1 购物车持久化

**现状**：`CartService` 使用 `sync.RWMutex` 保护的 `map[string]map[string]string` 内存结构（`services/order/service/cart_service.go:18-25`），服务重启后所有购物车数据丢失，且无法水平扩展。

**建议**：将购物车数据存入 Redis，以 `seat_id` 为 key。如果暂时不想引入 Redis，至少存入 PostgreSQL（如 `wtb_order` 库新增 `carts` 表）。

### 1.2 余额扣款非原子操作

**现状**：`DeductBalance`（`services/user/service/user_service.go:276-311`）执行"查询用户 → 判断余额 → 更新余额"三步操作，无事务包裹、无行锁。并发扣款时可能超扣。

**建议**：
```go
// 使用 UPDATE ... WHERE balance >= ? 保证原子性
result := r.db.Model(&model.User{}).
    Where("id = ? AND balance >= ?", userID, amount).
    UpdateColumn("balance", gorm.Expr("balance - ?", amount))
if result.RowsAffected == 0 {
    return errors.New("余额不足或用户不存在")
}
```

### 1.3 下单无事务

**现状**：`CreateOrder`（`services/order/service/order_service.go:7-14`）依次执行 `orderRepo.Create`、`itemRepo.Create`（循环）、`logRepo.Create`、`cartSvc.Clear`，中间任何一步失败都会造成数据不一致。

**建议**：用 `gorm.DB.Transaction` 包裹整个下单流程。

### 1.4 GORM AutoMigrate 在生产环境自动执行

**现状**：`backend/main.go:105-129` 在每次启动时对所有 8 个数据库执行 `AutoMigrate`。生产环境可能意外修改表结构。

**建议**：通过环境变量控制（如 `AUTO_MIGRATE=true`），生产环境默认关闭，数据库变更走独立 migration 脚本或工具。

---

## 二、后端架构（高优先级）

### 2.1 8 个独立数据库连接

**现状**：`backend/main.go:70-101` 用同一 DSN 拼接不同 `dbname` 分别 `gorm.Open`，产生 8 个独立连接池。默认 GORM 连接池为 10，总共 80 个连接对 PostgreSQL 压力较大。

**建议**：
- 合并为一个连接，通过 `dbname` 参数切换（GORM 支持 `Table("db.table")` 跨库查询）
- 或设置合理的连接池上限：`sqlDB.SetMaxOpenConns(5)`、`sqlDB.SetMaxIdleConns(2)`

### 2.2 优雅关闭

**现状**：`r.Run(":8080")` 是阻塞调用，无 `SIGTERM`/`SIGINT` 处理。

**建议**：
```go
srv := &http.Server{Addr: ":8080", Handler: r}
go func() { srv.ListenAndServe() }()
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
srv.Shutdown(context.Background())
```

### 2.3 健康检查不检查依赖

**现状**：`/health`（`backend/router.go:59-61`）永远返回 200，不检查数据库连通性。

**建议**：在健康检查中 `ping` 数据库，返回各依赖状态。

### 2.4 无访问限流

**现状**：所有公开接口（登录、菜单、购物车、扫码）无限流保护，可被恶意刷接口。

**建议**：用 Gin 中间件（如 `ulule/limiter`）对公开接口做 IP 级别限流，尤其保护 `/api/user/wx-login` 和 `/api/order/cart/*`。

### 2.5 CORS 配置过于宽松

**现状**：`backend/router.go:23-39` 的 `CORSMiddleware` 对任何 `Origin` 返回 `Access-Control-Allow-Origin: *` 且开启 `Credentials: true`。浏览器会拒绝 credentials + wildcard 的组合，但配置本身不规范。

**建议**：生产环境白名单允许的域名，开发环境才放开。

---

## 三、小程序端优化（中优先级）

### 3.1 加减按钮无防抖

**现状**：`menu.js:208-248` 的 `plus()`/`minus()` 直接调 `syncCart()` 发 HTTP 请求，快速点击会产生大量并发请求。

**建议**：加 300ms 防抖，或使用"本地先更新 UI → 批量同步"模式。

### 3.2 syncCart 重复拉取整个购物车

**现状**：每次 `syncCart()` 成功后都会调 `loadCart()` 全量刷新（`menu.js:242-243`），且 `loadCart()` 内部又做大量价格计算。快速操作时浪费严重。

**建议**：`syncCart` 成功后仅用返回数据局部更新 cartItems，不再全量拉取；或服务端 `cart/add` 接口直接返回最新的购物车数据。

### 3.3 自动登录时序竞态

**现状**：`menu.js:36` 用 `setTimeout(() => this.loadMemberInfo(), 500)` 等待登录完成，纯靠运气。

**建议**：`app.js` 的 `autoLogin()` 完成后通过 EventBus 或全局标志通知各页面，页面监听该事件后再加载需要鉴权的数据。

### 3.4 品种列表硬编码

**现状**：宠物品种 `BREED_OPTIONS` 硬编码在前端（`miniprogram/pages/mine/pets/add.js:4`）。

**建议**：品种数据从后端接口获取，方便运营后台动态维护。

### 3.5 会员价计算逻辑前后端重复

**现状**：会员价计算逻辑（member_level ≥ 1 用 member_price）在 `menu.js:93-95`、`menu.js:126-130` 等多个位置重复，后端 `user_service.go:485-494` 也有一份。

**建议**：统一由后端计算并返回菜品最终价格（前端只做展示），避免两端逻辑不一致。

### 3.6 缺少骨架屏/加载状态

**现状**：菜品列表、宠物列表、订单列表等页面无骨架屏，数据加载期间显示空白，体验较差。

**建议**：关键列表页增加骨架屏或 `wx:if` 加载态。

---

## 四、后台管理端优化（中优先级）

### 4.1 订单提醒使用轮询而非推送

**现状**：`OrderAlert.jsx:42` 每隔 5 秒轮询 `/api/order/admin/today-paid`，空转消耗资源。

**建议**：
- 短期：轮询间隔适当增大（如 15-30s），并加请求去重
- 长期：使用 WebSocket 或 SSE 推送新订单通知

### 4.2 Token 存储在 localStorage

**现状**：`admin-web/src/api/index.js:4` 将 JWT 存在 `localStorage`，XSS 攻击可窃取。

**建议**：改用 httpOnly cookie（需后端配合 `Set-Cookie`），或至少使用 `sessionStorage` + 较短过期时间。

### 4.3 无 React Error Boundary

**现状**：`App.jsx` 没有任何 Error Boundary，子组件异常会导致整个页面白屏。

**建议**：在 Layout 层包裹 Error Boundary 组件。

### 4.4 管理员密码硬编码

**现状**：管理员登录验证在 `services/admin/handler/` 中使用硬编码凭据。

**建议**：管理员账号密码存入数据库（`wtb_user` 库 `admin_users` 表），密码加盐哈希存储。

---

## 五、代码质量（中优先级）

### 5.1 order_service.go 可读性极差

**现状**：`services/order/service/order_service.go` 整个文件只有 22 行，多行语句用分号拼接在一行，极度压缩。

**建议**：展开为标准 Go 格式（`gofmt` 即可自动修复）。

### 5.2 错误被静默忽略

**现状**：多处 `json.Unmarshal`、`repo.FindByID` 的错误被忽略（如 `cart_service.go:48`、`user_service.go:300`）。

**建议**：对所有 `_` 丢弃的 error 做判断，至少打日志。关键路径（如余额扣款后查新余额失败）应返回错误而非继续。

### 5.3 无结构化日志

**现状**：仅使用 `log.Printf` + `fmt.Printf`，无日志级别、无 traceID。

**建议**：引入 `zap` 或 `zerolog`，在 Gin 中间件中注入 requestID，方便排查问题。

### 5.4 错误码不统一

**现状**：错误码混用 40001、50001、40101、40102 等不同格式。

**建议**：定义统一的错误码枚举，如：
- 4xxxx：客户端错误（参数校验、权限等）
- 5xxxx：服务端错误

### 5.5 缺少单元测试

**现状**：项目中引用了 `_test.go` 文件路径但实际测试覆盖率很低。

**建议**：
- 优先给 service 层（`user_service.go`、`order_service.go`、`cart_service.go`）补测试
- 余额扣款、创建订单、会员升级等涉及金钱的核心逻辑必须覆盖

---

## 六、性能优化（低优先级）

### 6.1 AdminListPets N+1 查询

**现状**：`user_service.go:406-437` 的 `AdminListPets` 先查宠物列表，再对每个宠物单独查用户信息（`FindByID`）。

**建议**：一次 JOIN 查询或批量 `WHERE id IN (…)` 取出所有用户。

### 6.2 菜品图片无 CDN

**现状**：图片通过 `r.Static("/images", imagePath)` 直接从服务器文件系统提供。

**建议**：上传到 OSS/S3 + CDN，减轻服务器带宽压力。

### 6.3 数据库无索引策略文档

**现状**：依赖 GORM AutoMigrate 自动创建的索引（主键 + 外键），未显式定义业务索引。

**建议**：至少对以下字段加索引：
- `users.openid`
- `orders.user_id`、`orders.status`、`orders.created_at`
- `order_items.order_id`
- `pets.user_id`
- `recharge_records.user_id`

---

## 七、安全加固（低优先级）

### 7.1 公开接口无鉴权的购物车操作

**现状**：`/api/order/cart/*` 不加 `auth` 中间件，仅靠 `seat_id` 隔离。恶意用户可遍历 `seat_id` 操作他人购物车。

**建议**：购物车接口加鉴权，将 `user_id` 与 `seat_id` 绑定。

### 7.2 扫码回调无签名校验

**现状**：`POST /api/pay/callback/wx` 对接微信支付回调，需确认已正确校验微信签名。

**建议**：确认微信支付回调签名校验逻辑是否完整（检查 `internal/wechat/` 中支付通知模块的实现）。

---

## 八、运维相关（低优先级）

### 8.1 无 Docker Compose 健康检查

**现状**：`docker-compose.yml` 中 backend 服务无 `healthcheck`。

**建议**：添加 `healthcheck: curl -f http://localhost:8080/health || exit 1`。

### 8.2 部署脚本无回滚机制

**现状**：`deploy.sh` 直接编译上传，无蓝绿/滚动发布。

**建议**：至少保留上一个版本的二进制备份（`backend-linux.prev`），出问题时快速回滚。

### 8.3 无 APM/监控

**现状**：无 Prometheus metrics、无 tracing。

**建议**：至少加 Gin 的 prometheus 中间件，暴露 `/metrics` 端点，监控 QPS、延迟、错误率。

---

## 优先级汇总

| 优先级 | 条目 | 风险 |
|--------|------|------|
| 高 | 1.2 余额扣款非原子 | 资金损失 |
| 高 | 1.3 下单无事务 | 数据不一致 |
| 高 | 1.1 购物车内存存储 | 服务重启丢数据 |
| 高 | 1.4 AutoMigrate 生产执行 | 误改表结构 |
| 高 | 2.2 无优雅关闭 | 请求中断 |
| 高 | 2.3 健康检查不查依赖 | 无法感知故障 |
| 中 | 3.1 按钮无防抖 | 重复请求 |
| 中 | 3.3 登录时序竞态 | 偶发功能异常 |
| 中 | 4.1 轮询替代推送 | 资源浪费 |
| 中 | 4.4 密码硬编码 | 安全风险 |
| 中 | 5.1-5.5 代码质量 | 维护成本 |
| 低 | 6.x 性能优化 | 当前规模可承受 |
| 低 | 7.x 安全加固 | 内网运营为主 |
| 低 | 8.x 运维 | 单机部署够用 |
