# WTB 部署指南（qs_admin 后端服务器）

> **目标服务器**：`qs_admin@192.168.0.156`（内网 SSH 免密）  
> **公网入口**：`hooper@10.144.144.1` → nginx 反向代理 → qs_admin  
> **域名**：`wtb.lqqnw.cn` / `wtbadm.lqqnw.cn`  
> **核心原则**：**日常更新只拷贝二进制，不要重新构建镜像**

---

## 一、当前架构

```
┌─────────────────────────────────────────────────────────────────┐
│                          公网入口                                │
│                  hooper (10.144.144.1)                          │
│                     wtb.lqqnw.cn                                │
│                         │                                       │
│                    nginx (443)                                  │
│                         │ proxy_pass                            │
│                         ▼                                       │
│              EasyTier 虚拟内网                                   │
│         ┌───────────────────────────┐                           │
│         │  http://10.144.144.2:18080 │ ←── 指向 qs_admin        │
│         └───────────────────────────┘                           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ EasyTier VPN (qsnet)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      qs_admin (后端服务器)                        │
│                  192.168.0.156 / 10.144.144.2                    │
│                                                                 │
│    ┌──────────────┐         ┌──────────────┐                   │
│    │ wtb-backend  │◄────────│ wtb-postgres │                   │
│    │   :8080      │         │   :5432      │                   │
│    └──────────────┘         └──────────────┘                   │
│         │                                                       │
│    0.0.0.0:18080 ◄── hooper nginx 反向代理                      │
└─────────────────────────────────────────────────────────────────┘
```

### 关键说明

- **hooper 是公网入口**：固定 IP `10.144.144.1`，域名 `wtb.lqqnw.cn` 解析到这里
- **hooper 只做反向代理**：本身**不运行**后端容器（已清理），所有流量通过 EasyTier 内网转发到 qs_admin
- **qs_admin 是实际后端**：运行 `wtb-backend` + `wtb-postgres` Docker 容器
- **EasyTier 组网**：hooper 和 qs_admin 在同一虚拟内网 `qsnet`，qs_admin 的虚拟 IP 是 `10.144.144.2`

---

## 二、部署原则（重要）

| 场景 | 推荐方式 | 耗时 |
|------|---------|------|
| **日常代码更新（推荐）** | 本地编译 → 拷贝二进制 → 容器内替换 → 重启 | ~10 秒 |
| 首次部署 / 基础镜像变更 | 构建镜像 → docker save → 上传加载 | 3-5 分钟 |
| 前端 admin-web 更新 | 本地 build → rsync 到 hooper | ~10 秒 |

**⚠️ 不要每次更新都构建 Docker 镜像！** 镜像构建 + save + 上传非常慢，日常只传二进制即可。

---

## 三、快速部署（日常更新，推荐）

```bash
# 1. 本地交叉编译（Mac 编译 Linux amd64 二进制）
cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o backend-linux .

# 2. 上传到 qs_admin
rsync -avz backend-linux qs_admin@192.168.0.156:~/wtb/

# 3. 热更新容器（无需重新构建镜像）
ssh qs_admin@192.168.0.156 "
  docker cp ~/wtb/backend-linux wtb-backend:/app/backend-linux
  docker restart wtb-backend
  sleep 2
  docker ps | grep wtb-backend
  curl -s http://127.0.0.1:18080/health
"
```

**注意**：如果修改了 `docker-compose.yml`（如新增环境变量），需要用 `docker compose up -d` 重新加载配置：
```bash
ssh qs_admin@192.168.0.156 "cd ~/wtb && docker compose up -d"
```

---

## 四、首次部署（或需要重建镜像时）

如果服务器上没有镜像，或 Dockerfile / 基础镜像有变更，才需要走完整流程：

```bash
# 1. 本地交叉编译
cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o backend-linux .

# 2. 准备构建上下文
cd ..
cp backend/backend-linux backend-linux

# 3. 构建镜像（本地有外网，服务器无外网无法 build）
DOCKER_BUILDKIT=0 docker build --platform linux/amd64 \
  -f deploy/qs-admin/Dockerfile -t wtb-backend:deploy .

# 4. 导出镜像
docker save -o /tmp/wtb-images.tar postgres:16-alpine wtb-backend:deploy

# 5. 上传镜像包 + 配置
rsync -avz --progress /tmp/wtb-images.tar qs_admin@192.168.0.156:~/wtb/
rsync -avz deploy/qs-admin/docker-compose.yml qs_admin@192.168.0.156:~/wtb/
rsync -avz init.sql qs_admin@192.168.0.156:~/wtb/

# 6. 服务器加载并启动
ssh qs_admin@192.168.0.156 "
  cd ~/wtb
  docker load -i wtb-images.tar
  docker compose up -d
  sleep 10
  docker ps | grep wtb
  curl -s http://127.0.0.1:18080/health
"
```

---

## 五、配置文件说明

### 5.1 docker-compose.yml（qs_admin）

```yaml
services:
  backend:
    image: wtb-backend:deploy
    container_name: wtb-backend
    environment:
      DB_DSN: "host=postgres user=admin password=xxx port=5432 sslmode=disable TimeZone=Asia/Shanghai"
      IMAGE_PATH: "/app/images"
      WX_APPID: "wx1e4315d1974c72f6"
      WX_APPSECRET: "your_secret_here"
      WX_ENV_VERSION: "trial"  # release / trial / develop
    ports:
      - "18080:8080"
    volumes:
      - ./images:/app/images
      - ./uploads:/app/uploads
    networks:
      - wtb-network

  postgres:
    image: postgres:16-alpine
    container_name: wtb-postgres
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: xxx
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - wtb-network
```

### 5.2 hooper 上的 nginx 配置

```nginx
server {
    server_name wtb.lqqnw.cn;

    location /images/ {
        alias /var/www/wtb/images/;
        expires 30d;
    }

    location / {
        proxy_pass http://10.144.144.2:18080;   # ← 通过 EasyTier 转发到 qs_admin
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    listen 443 ssl;
    ssl_certificate /etc/letsencrypt/live/wtb.lqqnw.cn/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/wtb.lqqnw.cn/privkey.pem;
}
```

> **⚠️ 重要**：hooper 上**不要**运行 `wtb-backend` docker 容器，nginx 直接代理到 qs_admin 的 EasyTier IP `10.144.144.2:18080`。

---

## 六、小程序真机调试配置

小程序 `app.js` 中的 `apiBase`：
```javascript
globalData: {
  apiBase: 'https://wtb.lqqnw.cn',   // ← 走 hooper 公网入口
}
```

**如果还没正式发布小程序**：微信接口（如生成小程序码 `getwxacodeunlimited`）可能会返回 `40066 invalid url`，因为接口要求小程序有已发布的版本。此时可用**方案 B**：小程序内加"手动选桌"弹窗，体验版码即可使用。

---

## 七、目录结构（qs_admin 服务器）

```
/home/qs_admin/wtb/
├── docker-compose.yml       # 编排文件
├── Dockerfile               # 首次构建用
├── wtb-images.tar           # 镜像包（首次部署后可选删除）
├── init.sql                 # 数据库初始化
├── images/                  # 菜品图片（持久化目录）
├── uploads/                 # 生成的二维码图片
│   └── seats/
│       └── seat_6.png
└── backend-linux            # 预编译二进制（日常热更新用）
```

---

## 八、问题速查

| 现象 | 原因 | 解决 |
|------|------|------|
| `docker pull` timeout | 服务器无外网 | 本机构建镜像 + `docker save` + `docker load` |
| `exec format error` | 二进制架构不匹配 | 本地编译时加 `GOARCH=amd64` |
| 接口返回 `appid missing` | 容器未重新加载环境变量 | `docker compose up -d`（不是 `docker restart`） |
| 生成小程序码 `40066` | 小程序未正式发布/未认证 | 先发布正式版，或改用"手动选桌"方案 |
| 图片 404 | images 目录缺失 | 检查 `~/wtb/images/` 是否存在 |
| nginx 502 | qs_admin 后端未启动 | `ssh qs_admin@192.168.0.156 "docker ps \| grep wtb"` |
| hooper 本机 18080 被占用 | 旧容器未清理 | `docker stop wtb-backend wtb-postgres && docker rm wtb-backend wtb-postgres` |

---

## 九、配套部署

- **前端 admin-web**：部署在 `anzhitek` 服务器，详见 `deploy/anzhitek/README.md`
- **hooper nginx**：域名 `wtb.lqqnw.cn` / `wtbadm.lqqnw.cn`，SSL 证书由 Certbot 管理
