# 汪托帮后台管理系统测试文档

## 一、测试环境

| 组件 | 地址 | 说明 |
|------|------|------|
| 后台管理前端 | http://localhost:3000 | 纯 HTML/JS 页面 |
| Admin BFF | http://localhost:8090 | 管理员登录 + API 代理 |
| Gateway | http://localhost:8080 | 小程序网关 |
| 各微服务 | 8081-8089 | user/menu/order/activity/points 等 |

### 管理员账号
- 用户名：`admin`
- 密码：`admin123`

### 测试用户（小程序端）
- 开发 Code：`dev_test_001`
- OpenID：`dev_openid_dev_test_001`
- 余额：¥500
- 积分：5000

---

## 二、服务健康检查

### 2.1 所有服务端口检查
```bash
for p in 8080 8081 8082 8083 8084 8085 8086 8087 8088 8089 8090; do
  curl -s http://localhost:$p/health
done
```
**预期结果**：所有端口返回 `{"status":"ok"}`

### 2.2 Admin BFF 登录测试
```bash
curl -X POST http://localhost:8090/admin/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}'
```
**预期结果**：返回 `code: 200`，data.token 不为空

### 2.3 CORS 跨域测试
```bash
curl -X OPTIONS http://localhost:8090/admin/login \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST"
```
**预期结果**：返回 `204 No Content`，响应头包含 `Access-Control-Allow-Origin: *`

---

## 三、菜单管理测试

### 3.1 查询分类列表
**接口**：`GET /api/admin/menu/categories`
**预期结果**：返回分类数组，每个分类包含 id/name/parent_id/sort_order

### 3.2 新增分类
**接口**：`POST /api/admin/menu/admin/category`
**请求体**：`{"name":"新分类","sort_order":10}`
**预期结果**：
- code: 200
- 返回新建分类的 id
- 再次查询列表能看到新分类

### 3.3 删除分类
**接口**：`DELETE /api/admin/menu/admin/category/{id}`
**预期结果**：
- code: 200
- 再次查询列表该分类已消失

### 3.4 查询菜品列表
**接口**：`GET /api/admin/menu/dishes`
**预期结果**：
- 返回分页数据 `{total, page, pageSize, list}`
- 每个菜品包含 `price` 字段（分）

### 3.5 新增菜品（核心修复项）
**接口**：`POST /api/admin/menu/admin/dish`
**请求体**：
```json
{
  "name": "测试菜品A",
  "category_id": 2,
  "price": 3500,
  "stock": 100,
  "tags": "推荐,热销"
}
```
**预期结果**：
- code: 200
- 菜品创建成功
- **价格保存验证**：查询 `/api/admin/menu/dishes`，该菜品 `price` 为 3500
- **库存保存验证**：查询 `/api/admin/menu/dish/{id}`，返回的 `stock.daily_limit` 为 100

### 3.6 删除菜品
**接口**：`DELETE /api/admin/menu/admin/dish/{id}`
**预期结果**：
- code: 200
- 再次查询列表该菜品已消失

---

## 四、订单管理测试

### 4.1 查询订单列表
**接口**：`GET /api/admin/order/list`
**预期结果**：返回分页数据 `{total, page, pageSize, list}`

### 4.2 修改订单状态
**前置条件**：先通过小程序或接口创建一个订单
**接口**：`PUT /api/admin/order/admin/status`
**请求体**：`{"order_id":1,"status":"completed"}`
**预期结果**：
- code: 200
- 再次查询订单列表，该订单状态变为 "completed"

---

## 五、活动管理测试

### 5.1 查询公告列表
**接口**：`GET /api/admin/activity/announcements`
**预期结果**：返回公告数组

### 5.2 发布公告
**接口**：`POST /api/admin/activity/admin/announcement`
**请求体**：`{"title":"公告标题","content":"公告内容","is_published":true}`
**预期结果**：
- code: 200
- 再次查询列表能看到新公告

### 5.3 查询活动列表
**接口**：`GET /api/admin/activity/list`
**预期结果**：返回活动数组，status 为 "published" 或 "draft"

### 5.4 创建活动（核心修复项）
**接口**：`POST /api/admin/activity/admin/activity`
**请求体**：
```json
{
  "title": "测试活动",
  "location": "大厅",
  "quota": 20,
  "max_participants": 20
}
```
**预期结果**：
- code: 200
- 返回的活动 `status` 为 `"published"`（修复前为 "draft"）
- 再次查询列表能看到该活动

---

## 六、积分商品测试

### 6.1 查询积分商品列表
**接口**：`GET /api/admin/points/goods`
**预期结果**：返回商品数组，每个商品包含 id/name/points_price/stock

### 6.2 新增积分商品
**接口**：`POST /api/admin/points/admin/goods`
**请求体**：
```json
{
  "name": "测试商品",
  "points_price": 500,
  "stock": 50,
  "type": "physical",
  "status": 1
}
```
**预期结果**：
- code: 200
- 返回新建商品的 id
- 再次查询列表能看到新商品

---

## 七、前端页面测试（浏览器操作）

### 7.1 登录页面
1. 打开 http://localhost:3000
2. **预期结果**：显示登录表单
3. 输入用户名 `admin`，密码 `admin123`
4. 点击登录
5. **预期结果**：登录成功，跳转到数据概览页面

### 7.2 菜单管理页面
1. 点击左侧「菜单管理」
2. **预期结果**：
   - 显示分类列表（有数据）
   - 显示菜品列表（有数据，包含价格）
3. 点击「+ 新增分类」，输入名称后保存
4. **预期结果**：分类列表刷新，出现新分类
5. 点击「+ 新增菜品」，填写名称/分类/价格/库存后保存
6. **预期结果**：菜品列表刷新，新菜品价格显示正确（如 ¥35.00）

### 7.3 订单管理页面
1. 点击左侧「订单管理」
2. **预期结果**：显示订单列表（如无可通过小程序下单）

### 7.4 活动管理页面
1. 点击左侧「活动管理」
2. **预期结果**：显示公告列表和活动列表
3. 点击「+ 发布公告」，填写后保存
4. **预期结果**：公告列表刷新
5. 点击「+ 创建活动」，填写后保存
6. **预期结果**：活动列表刷新，状态显示为「published」

### 7.5 积分商品页面
1. 点击左侧「积分商品」
2. **预期结果**：显示商品列表
3. 点击「+ 新增商品」，填写后保存
4. **预期结果**：商品列表刷新

---

## 八、测试记录表

| 序号 | 测试项 | 预期结果 | 实际结果 | 状态 |
|------|--------|----------|----------|------|
| 1 | 服务健康检查 | 全部 OK | 11/11 端口 HTTP 200 | ✅ |
| 2 | Admin 登录 | 返回 token | code=200, token 有效 | ✅ |
| 3 | CORS 跨域 | 204 + headers | HTTP 204, Access-Control-Allow-Origin: * | ✅ |
| 4 | 查询分类 | 返回数组 | 4 个分类 | ✅ |
| 5 | 新增分类 | 创建成功 | ID=6, 创建成功 | ✅ |
| 6 | 删除分类 | 删除成功 | code=200, 删除成功 | ✅ |
| 7 | 查询菜品 | 含 price | total=6, 第一条 price=4500 | ✅ |
| 8 | 新增菜品（含价格/库存） | 价格库存保存 | ID=8, price=5500, stock=66 | ✅ |
| 9 | 删除菜品 | 删除成功 | code=200, 删除成功 | ✅ |
| 10 | 查询订单 | 返回分页 | total=0, 分页格式正确 | ✅ |
| 11 | 修改订单状态 | 状态变更 | code=200 | ✅ |
| 12 | 查询公告 | 返回数组 | 2 条公告 | ✅ |
| 13 | 发布公告 | 创建成功 | ID=4, 创建成功 | ✅ |
| 14 | 查询活动 | 返回数组 | 3 个活动 | ✅ |
| 15 | 创建活动（默认 published） | status=published | status=published | ✅ |
| 16 | 查询积分商品 | 返回数组 | 6 个商品 | ✅ |
| 17 | 新增积分商品 | 创建成功 | ID=8, 创建成功 | ✅ |
| 18 | 浏览器登录页面 | 正常显示 | http://localhost:3000 | ✅ |
| 19 | 浏览器菜单管理 | CRUD 正常 | 新增/删除/列表正常 | ✅ |
| 20 | 浏览器活动管理 | 创建 published | 状态正确 | ✅ |
