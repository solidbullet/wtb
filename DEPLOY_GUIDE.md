# WTB 项目部署指南（生产版）

> ⚠️ **文档已整理**：本项目支持多环境部署，配置文件已按目标服务器分类存放：
> - `deploy/hooper/` —— 生产环境（hooper@10.144.144.1，nginx + Docker）
> - `deploy/qs-admin/` —— 内网环境（qs_admin@192.168.0.156，内网穿透，无外网）
>
> 本文档继续保留作为总览参考，具体部署脚本和配置文件请进入对应目录查看。

> **目标**：读完本文档后，部署一次成功，不出现图片 404、容器冲突、配置错误等低级问题。  
> **哲学**：把事后踩坑变成事前检查，不通过的检查项阻塞部署流程。

---

## 一、架构总览（必须先理解）

### 1.1 域名与服务映射

| 域名 | 用途 | nginx 处理方式 |
|------|------|---------------|
| `wtb.lqqnw.cn` | 小程序 API + 图片 | `/images/` → 本地文件；`/` → 后端容器 |
| `wtbadm.lqqnw.cn` | admin-web 管理后台 | `/images/` → 后端容器；`/api/` `/admin/` → 后端容器；`/` → 静态文件 |

### 1.2 图片存储统一架构（核心）

**绝对不能分裂。** 后端上传图片必须写入同一个物理目录，两个域名都能访问。

```
┌─────────────────────────────────────────┐
│  宿主机 /var/www/wtb/images/            │
│  （唯一真实图片目录）                    │
└──────────────┬──────────────────────────┘
               │
    ┌──────────┴──────────┐
    ▼                     ▼
nginx (wtb)          nginx (wtbadm)
alias 直接读取        proxy_pass → 后端容器
                         │
                    volume 映射
                         ▼
              容器 /miniprogram/images/
```

**关键**：docker-compose.yml 中 volume 必须是 `/var/www/wtb/images:/miniprogram/images`，
这样后端上传 → 实际写到 `/var/www/wtb/images/` → 两个域名都可见。

---

## 二、首次部署（从零开始，一次性）

### 2.1 服务器环境准备

```bash
见/Users/admin/workspace/jyq/document/account.md

# 检查 Docker
docker --version        # 期望 20.10+
docker compose version  # 期望 v2+

# 检查 nginx
nginx -v

# 创建目录（一次性）
sudo mkdir -p /var/www/wtb/images
sudo mkdir -p /home/hooper/wtb
sudo chown -R hooper:hooper /var/www/wtb
```

### 2.2 传输文件到服务器

> ⚠️ **重要**：公网 IP `118.145.193.23` 带宽小、易超时。务必使用 EasyTier 虚拟内网 `10.144.144.1` 传输。

**必须上传的文件清单：**

| 本地路径 | 服务器路径 | 说明 |
|----------|-----------|------|
| `backend/Dockerfile` | `~/wtb/Dockerfile` | 后端镜像构建文件 |
| `backend/router.go` 等源码 | `~/wtb/` 对应目录 | 后端源码（用于 Docker build） |
| `docker-compose.yml` | `~/wtb/docker-compose.yml` | 编排文件（见下方标准模板） |
| `init.sql` | `~/wtb/init.sql` | 数据库初始化 |
| `miniprogram/images/*.png` | `/var/www/wtb/images/` | 菜品图片（必须传！） |
| `admin-web/dist/` | `/var/www/wtbadm/` | 管理后台前端（如有更新） |

**推荐传输方式（EasyTier）：**

```bash
# 1. 本地压缩图片
cd workspace/jyq/wtb
tar czf /tmp/images.tar.gz -C miniprogram/images .

# 2. 通过 EasyTier 内网传（稳定、速度快）
rsync -avz -e "ssh" /tmp/images.tar.gz hooper@10.144.144.1:/tmp/

# 3. 服务器上解压到统一目录
ssh hooper@10.144.144.1 "sudo tar xzf /tmp/images.tar.gz -C /var/www/wtb/images/"
```

### 2.3 部署标准配置文件

**以下三份配置是"黄金标准"，直接复制使用，不要自己改。**

#### A. docker-compose.yml（`~/wtb/docker-compose.yml`）

```yaml
version: "3.8"

services:
  postgres:
    image: postgres:16-alpine
    container_name: wtb-postgres
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: JyRUj7wlNjU0uVHh
    ports:
      - 127.0.0.1:5432:5432
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - wtb-network
    healthcheck:
      test: [CMD-SHELL, pg_isready -U admin]
      interval: 5s
      timeout: 5s
      retries: 5

  backend:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: wtb-backend
    environment:
      DB_DSN: "host=postgres user=admin password=JyRUj7wlNjU0uVHh port=5432 sslmode=disable TimeZone=Asia/Shanghai"
      IMAGE_PATH: "/miniprogram/images"
      WX_APPID: "wx1e4315d1974c72f6"
      WX_APPSECRET: ""
    ports:
      - 127.0.0.1:18080:8080
    volumes:
      # ⚠️ 核心：统一挂载到 /var/www/wtb/images，确保上传后双域名可见
      - /var/www/wtb/images:/miniprogram/images
    networks:
      - wtb-network
    restart: unless-stopped
    depends_on:
      postgres:
        condition: service_healthy

volumes:
  pgdata:

networks:
  wtb-network:
```

#### B. nginx wtb.lqqnw.cn（`/etc/nginx/sites-enabled/wtb.lqqnw.cn`）

```nginx
server {
    server_name wtb.lqqnw.cn;

    location /images/ {
        alias /var/www/wtb/images/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    location / {
        proxy_pass http://127.0.0.1:18080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    listen 443 ssl;
    ssl_certificate /etc/letsencrypt/live/wtb.lqqnw.cn/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/wtb.lqqnw.cn/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;
}

server {
    if ($host = wtb.lqqnw.cn) {
        return 301 https://$host$request_uri;
    }
    listen 80;
    server_name wtb.lqqnw.cn;
    return 404;
}
```

#### C. nginx wtbadm.lqqnw.cn（`/etc/nginx/sites-enabled/wtbadm`）

> ⚠️ **关键**：`location ^~ /images/` 必须在 `location ~* \.(png|...)$` **之前**，否则正则匹配会抢先！

```nginx
server {
    server_name wtbadm.lqqnw.cn;

    root /var/www/wtbadm;
    index index.html;

    # 前端路由支持
    location / {
        try_files $uri $uri/ /index.html;
    }

    # ⚠️ 核心：/images/ 必须用 ^~，且放在正则 location 之前！
    location ^~ /images/ {
        proxy_pass http://127.0.0.1:18080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /api/ {
        proxy_pass http://127.0.0.1:18080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /admin/ {
        proxy_pass http://127.0.0.1:18080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /health {
        proxy_pass http://127.0.0.1:18080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
    }

    # 静态资源缓存（必须放在 /images/ 之后！）
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|webp|woff|woff2|ttf)$ {
        expires 30d;
        add_header Cache-Control "public";
    }

    listen 443 ssl;
    ssl_certificate /etc/letsencrypt/live/wtbadm.lqqnw.cn/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/wtbadm.lqqnw.cn/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;
}

server {
    if ($host = wtbadm.lqqnw.cn) {
        return 301 https://$host$request_uri;
    }
    listen 80;
    server_name wtbadm.lqqnw.cn;
    return 404;
}
```

### 2.4 启动服务

```bash
cd /home/hooper/wtb

# 验证 YAML 格式（不通过不能继续！）
docker compose config > /dev/null || { echo "❌ docker-compose.yml 格式错误"; exit 1; }

# 构建并启动
docker compose down 2>/dev/null || true
docker compose build backend
docker compose up -d

# 验证 nginx
sudo nginx -t && sudo nginx -s reload
```

> ⚠️ **重要**：`docker compose build backend` 在服务器上执行 `go build` 会重新下载所有 Go 依赖，**首次部署可能需要 10-30 分钟**。日常代码小改动请使用下方「超快速部署」方案，10 秒内完成。

---

## 三、日常更新部署（代码小改动）

> **核心原则**：服务器上没有 Go 模块缓存，`docker build` 每次都会重新下载依赖，非常慢。  
> **正确做法**：本地编译 → 只传二进制 → 替换容器内文件 → 重启。**全程 10-20 秒。**

### 3.1 超快速部署脚本（推荐）

```bash
#!/bin/bash
set -e

echo "========================================"
echo "  WTB 超快速部署（本地编译 + 热更新）"
echo "========================================"

# 1. 本地交叉编译（利用本机 Go 缓存，秒级完成）
echo "🔨 本地编译..."
cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o backend-linux .
cd ..

# 2. 只传变更的二进制（rsync 增量，通常 1-3 秒）
echo "📤 上传二进制..."
rsync -avz backend/backend-linux hooper@10.144.144.1:~/wtb/backend/

# 3. 服务器上替换并重启（5 秒内完成）
echo "🔄 热更新容器..."
ssh hooper@10.144.144.1 "
  docker cp ~/wtb/backend/backend-linux wtb-backend:/app/backend
  docker restart wtb-backend
  sleep 2
  echo '容器状态:'
  docker ps | grep wtb-backend
"

# 4. 验证
echo "🧪 验证..."
curl -s -o /dev/null -w "  wtb API: %{http_code}\n" https://wtb.lqqnw.cn/api/menu/categories
curl -s -o /dev/null -w "  wtb 图片: %{http_code}\n" https://wtb.lqqnw.cn/images/hongshao.png
curl -s -o /dev/null -w "  wtbadm 图片: %{http_code}\n" https://wtbadm.lqqnw.cn/images/hongshao.png

echo ""
echo "========================================"
echo "  ✅ 部署完成"
echo "========================================"
```

**适用场景**：
- ✅ 后端 Go 代码有改动（API、业务逻辑等）
- ✅ 只改动了一两个文件
- ✅ 需要秒级上线

**不适用场景**：
- ❌ Dockerfile 本身有改动（需要 `docker compose build`）
- ❌ 新增/修改了静态资源且 Dockerfile 的 `COPY` 层需要更新
- ❌ 依赖了 `miniprogram/images` 目录的新增图片（需要重新 build 镜像）

### 3.2 什么时候必须重新 build 镜像

以下情况必须走 `docker compose build backend`：

1. 修改了 `Dockerfile`
2. 在 `miniprogram/images/` 新增了图片且 Dockerfile 中 `COPY miniprogram/images /app/images` 需要同步
3. 首次部署（服务器上没有镜像缓存）
4. 需要清理 Docker 缓存层

```bash
cd /home/hooper/wtb
docker compose build backend
docker compose up -d
```

### 3.3 更新前强制检查脚本

**在服务器上运行此脚本，全部 PASS 才能继续部署：**

```bash
#!/bin/bash
set -e

PASS=0
FAIL=0

check() {
    if eval "$2" > /dev/null 2>&1; then
        echo "  ✅ PASS: $1"
        PASS=$((PASS+1))
    else
        echo "  ❌ FAIL: $1"
        FAIL=$((FAIL+1))
    fi
}

echo "========================================"
echo "  WTB 部署前强制检查清单"
echo "========================================"

echo ""
echo "【文件检查】"
check "init.sql 是普通文件（不是目录）" "[ -f /home/hooper/wtb/init.sql ]"
check "docker-compose.yml 存在" "[ -f /home/hooper/wtb/docker-compose.yml ]"
check "后端 Dockerfile 存在" "[ -f /home/hooper/wtb/Dockerfile ]"

echo ""
echo "【图片检查】"
check "nginx 图片目录存在" "[ -d /var/www/wtb/images ]"
check "图片目录有菜品图片" "ls /var/www/wtb/images/*.png > /dev/null 2>&1"
check "hongshao.png 存在" "[ -f /var/www/wtb/images/hongshao.png ]"

echo ""
echo "【配置检查】"
check "docker-compose.yml 格式正确" "cd /home/hooper/wtb && docker compose config > /dev/null 2>&1"
check "nginx 语法正确" "sudo nginx -t > /dev/null 2>&1"
check "wtbadm nginx 配置有 ^~ /images/" "grep -q 'location ^~ /images/' /etc/nginx/sites-enabled/wtbadm"
check "wtbadm nginx 配置 /images/ 在正则之前" "bash -c 'IMG_LINE=\$(grep -n \"location ^~ /images/\" /etc/nginx/sites-enabled/wtbadm | head -1 | cut -d: -f1); REG_LINE=\$(grep -n \"location ~*\" /etc/nginx/sites-enabled/wtbadm | head -1 | cut -d: -f1); [ -n \"\$IMG_LINE\" ] && [ -n \"\$REG_LINE\" ] && [ \"\$IMG_LINE\" -lt \"\$REG_LINE\" ]'"

echo ""
echo "【容器检查】"
check "postgres 容器在运行" "docker ps | grep -q wtb-postgres"
check "backend 容器在运行" "docker ps | grep -q wtb-backend"

echo ""
echo "【网络检查】"
check "wtb.lqqnw.cn API 可访问" "curl -s -o /dev/null -w '%{http_code}' https://wtb.lqqnw.cn/api/menu/categories | grep -q '200'"
check "wtb.lqqnw.cn 图片可访问" "curl -s -o /dev/null -w '%{http_code}' https://wtb.lqqnw.cn/images/hongshao.png | grep -q '200'"
check "wtbadm.lqqnw.cn 图片可访问" "curl -s -o /dev/null -w '%{http_code}' https://wtbadm.lqqnw.cn/images/hongshao.png | grep -q '200'"

echo ""
echo "========================================"
echo "  结果: $PASS 通过, $FAIL 失败"
echo "========================================"

if [ $FAIL -gt 0 ]; then
    echo "❌ 有检查项未通过，禁止部署！先修复上述问题。"
    exit 1
else
    echo "✅ 全部通过，可以安全部署。"
fi
```

保存为 `/home/hooper/wtb/pre-deploy-check.sh`，每次部署前执行：

```bash
chmod +x /home/hooper/wtb/pre-deploy-check.sh
/home/hooper/wtb/pre-deploy-check.sh
```

### 3.4 完整重建部署脚本（首次/镜像更新）

以下脚本适合首次部署或必须重新 build 镜像的场景。注意：服务器上 `go build` 下载依赖很慢（10-30 分钟），日常小改动请用上方「超快速部署」。

```bash
#!/bin/bash
set -e

echo "========================================"
echo "  WTB 完整重建部署脚本"
echo "========================================"

cd /home/hooper/wtb

# 1. 强制检查（不通过直接退出）
echo ""
echo "🔍 运行强制检查..."
./pre-deploy-check.sh || { echo "❌ 检查未通过，中止部署"; exit 1; }

# 2. 停止旧容器
echo ""
echo "🛑 停止旧容器..."
docker compose down 2>/dev/null || true

# 3. 构建（不用 --no-cache）
echo ""
echo "🔨 构建后端镜像（服务器下载依赖，可能需要 10-30 分钟）..."
docker compose build backend

# 4. 启动
echo ""
echo "🚀 启动服务..."
docker compose up -d

# 5. 等待就绪
echo ""
echo "⏳ 等待数据库就绪..."
sleep 8

# 6. 部署后验证
echo ""
echo "📋 部署后验证:"
curl -s -o /dev/null -w "  wtb API: %{http_code}\n" https://wtb.lqqnw.cn/api/menu/categories
curl -s -o /dev/null -w "  wtb 图片: %{http_code}\n" https://wtb.lqqnw.cn/images/hongshao.png
curl -s -o /dev/null -w "  wtbadm 图片: %{http_code}\n" https://wtbadm.lqqnw.cn/images/hongshao.png

echo ""
echo "========================================"
echo "  ✅ 部署完成"
echo "========================================"
```

---

---

## 四、小程序包大小检查

上传小程序前务必检查包大小，防止动态图片混入本地包：

```bash
# 检查总大小
du -sh miniprogram/

# 检查是否有异常大图（超过 500KB 的文件）
find miniprogram/ -type f -size +500k

# 动态上传的图片（dish_*.png）绝不应出现在 miniprogram/images/
ls miniprogram/images/dish_*.png 2>/dev/null && echo "❌ 发现动态图片，需删除" || echo "✅ 无动态图片"
```

**根因**：后端 `UploadImage` 接口默认保存路径为 `uploads/`（已修复），旧代码为 `../miniprogram/images`，导致本地开发时上传的图片直接写入小程序目录。

---

## 五、更新 admin-web 前端

admin-web 更新后需要重新构建并上传到服务器：

```bash
# 本地构建
cd admin-web
npm run build

# 上传到服务器（通过 EasyTier 内网）
rsync -avz --delete -e "ssh" dist/ hooper@10.144.144.1:/var/www/wtbadm/

# 不需要重启任何服务，nginx 直接服务静态文件
```

---

## 六、回滚操作

```bash
cd /home/hooper/wtb

# Docker 回滚
docker stop wtb-backend
docker rm wtb-backend
docker compose up -d backend

# 或者用上一个镜像
docker images | grep wtb-backend
# docker run ... wtb-backend:<上一个tag>
```

---

## 七、systemd 配置（二进制部署方案）

如果不用 Docker，直接运行 backend-linux：

```bash
sudo tee /etc/systemd/system/wtb-backend.service << 'EOF'
[Unit]
Description=WTB Backend Service
After=network.target postgresql.service

[Service]
Type=simple
User=hooper
WorkingDirectory=/home/hooper/wtb
ExecStart=/home/hooper/wtb/backend-linux
Restart=always
RestartSec=5
StartLimitInterval=60s
StartLimitBurst=3
Environment="DB_DSN=host=127.0.0.1 user=admin password=JyRUj7wlNjU0uVHh port=5432 sslmode=disable TimeZone=Asia/Shanghai"
Environment="IMAGE_PATH=/var/www/wtb/images"
Environment="WX_APPID=wx1e4315d1974c72f6"
Environment="WX_APPSECRET="

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable wtb-backend
sudo systemctl start wtb-backend
```

---

## 八、问题速查表

| 现象 | 快速诊断 | 解决 |
|------|---------|------|
| 容器启动报错 `exec format error` | `file backend-linux` 显示 ARM64 | 在服务器上构建，或交叉编译 `GOOS=linux GOARCH=amd64` |
| postgres 启动报错 `init.sql is a directory` | `ls -la init.sql` 显示 `drwx...` | `rm -rf init.sql` 后重建文件 |
| docker compose 报错 `version must be a string` | YAML 格式错误 | 改用双引号 `version: "3.8"` |
| 容器名冲突 | `docker ps -a` 看到旧容器 | `docker rm -f wtb-backend` |
| wtb 图片 404 | `curl -I https://wtb.lqqnw.cn/images/xxx.png` | 检查 `/var/www/wtb/images/` 是否有文件 |
| wtbadm 图片 404 | `curl -I https://wtbadm.lqqnw.cn/images/xxx.png` | 检查 nginx 配置中 `^~ /images/` 是否在正则之前 |
| 上传后小程序看不到 | 检查两个目录是否一致 | docker-compose volume 改为 `/var/www/wtb/images:/miniprogram/images` |
| 小程序上传报 80051（包超 2MB） | `du -sh miniprogram/` | 删除 `miniprogram/images/dish_*.png` 等动态大图 |
| 部署极慢（go build 卡死） | 服务器下载 Go 依赖慢 | 改用「超快速部署」：本地编译 + docker cp + restart |
| scp 上传超时 | 公网不稳定 | 使用 EasyTier 内网 `10.144.144.1` |

---

**最后更新：2026-05-27**
> 本次更新：新增「超快速部署」方案、小程序包大小检查、图片上传路径修复说明。  
**维护者：部署前务必运行 `pre-deploy-check.sh`，不通过禁止部署。**
