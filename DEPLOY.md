# 汪托帮点餐系统 - 部署指南

## 一、服务器要求

- Linux 服务器（推荐 Ubuntu 22.04 / CentOS 8）
- 已安装 Docker 和 docker-compose
- 域名 `wtb.lqqnw.cn` 已解析到服务器公网 IP
- 开放端口：80（HTTP）、443（HTTPS）

## 二、上传文件到服务器

将 `wtb-deploy.tar.gz` 上传到服务器任意目录（例如 `/opt/wtb`）：

```bash
# 在你的电脑上执行
scp wtb-deploy.tar.gz root@你的服务器IP:/opt/

# 登录服务器
ssh root@你的服务器IP
cd /opt
tar xzf wtb-deploy.tar.gz
cd wtb-deploy
```

## 三、一键部署

```bash
chmod +x deploy.sh
./deploy.sh
```

部署完成后，服务自动运行在：
- **HTTP**: http://wtb.lqqnw.cn
- **HTTPS**: https://wtb.lqqnw.cn （推荐 ✅）

## 四、HTTPS 自动配置说明

本项目使用 **Caddy** 作为反向代理，它会：
1. 自动为 `wtb.lqqnw.cn` 申请 Let's Encrypt 免费证书
2. 自动将 HTTP 重定向到 HTTPS
3. 自动续期证书（无需手动操作）

> ⚠️ 首次申请证书可能需要 10-30 秒，请稍等后再访问 HTTPS。

## 五、验证部署

```bash
# 查看服务运行状态
docker-compose ps

# 查看后端日志
docker-compose logs -f backend

# 查看 Caddy 日志（证书申请情况）
docker-compose logs -f caddy

# 测试 API
curl https://wtb.lqqnw.cn/api/menu/list
```

## 六、配置微信小程序

### 6.1 修改小程序 API 地址

打开 `miniprogram/utils/config.js`，将 `LAN_IP` 改为域名：

```javascript
const LAN_IP = 'wtb.lqqnw.cn'
```

同时修改协议为 HTTPS：

```javascript
const API_BASE = isDevTools
  ? 'http://localhost:8080'
  : `https://${LAN_IP}`
```

### 6.2 修改云函数中的后端地址

打开 `cloudfunctions/proxy/index.js`：

```javascript
const API_BASE = 'https://wtb.lqqnw.cn'
```

### 6.3 配置微信小程序后台

1. 登录 [微信公众平台](https://mp.weixin.qq.com/)
2. 进入「开发」→「开发管理」→「开发设置」
3. 在「服务器域名」→「request 合法域名」中添加：
   ```
   https://wtb.lqqnw.cn
   ```
4. 如果使用了云函数，还需要在「上传域名」中添加：
   ```
   https://wtb.lqqnw.cn
   ```

### 6.4 上传小程序代码

1. 微信开发者工具 → 点击「上传」
2. 填写版本号和项目备注
3. 上传成功后，登录微信公众平台 →「版本管理」→ 将上传的版本设为「体验版」或「提交审核」

## 七、配置真实微信登录（可选）

当前后端支持 fallback 模式（未配置 AppSecret 时自动生成 mock_openid），如果要使用真实的微信登录：

1. 在 [微信公众平台](https://mp.weixin.qq.com/) →「开发」→「开发管理」→「开发设置」中获取 **AppID** 和 **AppSecret**

2. 修改 `docker-compose.yml`，取消注释并填写：
   ```yaml
   backend:
     environment:
       WX_APPID: "wx1e4315d1974c72f6"
       WX_APPSECRET: "your_real_app_secret_here"
   ```

3. 重新部署：
   ```bash
   docker-compose down
   docker-compose up -d
   ```

## 八、数据持久化

数据库数据存储在 Docker Volume 中：

```bash
# 查看数据卷
docker volume ls
```

即使删除容器，数据也不会丢失。如需备份：

```bash
docker exec wtb-postgres pg_dumpall -U admin > backup.sql
```

## 九、常见问题

### Q1: 访问 https://wtb.lqqnw.cn 显示证书错误？
- 首次部署后等待 30 秒，Caddy 正在申请证书
- 检查域名是否正确解析到服务器 IP：`ping wtb.lqqnw.cn`
- 查看 Caddy 日志：`docker-compose logs -f caddy`

### Q2: 部署后 API 无法访问？
- 检查服务器安全组/防火墙是否开放 80 和 443 端口
- 检查服务状态：`docker-compose ps`
- 检查后端日志：`docker-compose logs -f backend`

### Q3: 如何更新代码后重新部署？
```bash
# 上传新代码到服务器后
docker-compose down
docker-compose up --build -d
```

### Q4: 如何查看后端日志？
```bash
docker-compose logs -f backend
```

## 十、技术架构

```
                          ┌─────────────────┐
                          │   微信小程序     │
                          │   (云函数)       │
                          └────────┬────────┘
                                   │ HTTPS
                          ┌────────▼────────┐
                     443  │   wtb-caddy      │  自动证书
                     80   │   (Caddy)        │  HTTP→HTTPS重定向
                          └────────┬────────┘
                                   │
                          ┌────────▼────────┐
                          │   wtb-backend    │ :8080
                          │   (Go Docker)    │
                          └────────┬────────┘
                                   │
                          ┌────────▼────────┐
                          │   wtb-postgres   │ :5432
                          │   (PostgreSQL)   │
                          └─────────────────┘
```
