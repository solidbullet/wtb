# WTB 部署指南（qs_admin @ 192.168.0.156）

> **目标服务器**：`qs_admin@192.168.0.156`（内网 SSH 免密）  
> **访问方式**：内网穿透（无需 nginx/caddy）  
> **核心挑战**：服务器无外网，无法拉取 Docker Hub 镜像  
> **解决方案**：**本机构建镜像 → 导出 tar → 上传加载**

---

## 一、环境特点

| 项目 | 状态 | 影响 |
|------|------|------|
| 公网访问 | ❌ 无外网 | 无法 `docker pull`，无法服务器上 `go build` 下载依赖 |
| SSH 登录 | ✅ 免密密钥 | `rsync` + `ssh` 直接可用，无需密码 |
| Docker | ✅ 已安装 | 可运行容器，但无法拉取新镜像 |
| nginx/caddy | ❌ 不需要 | 通过内网穿透暴露服务，无需反向代理 |
| Go 环境 | ❌ 未安装 | 必须在本地交叉编译 |

---

## 二、部署原理

```
┌─────────────┐      rsync       ┌─────────────────────────────┐
│   本机      │ ───────────────> │  192.168.0.156 (qs_admin)   │
│  (有外网)   │                  │                             │
│             │  1. 本地编译      │  docker load -i wtb-images.tar
│  Go编译器   │  2. docker build  │  docker compose up -d       │
│  Docker     │  3. docker save   │                             │
│  Clash代理  │  4. rsync上传     │  端口 18080 暴露给内网穿透    │
└─────────────┘                  └─────────────────────────────┘
```

**为什么不用服务器直接 build？**
1. 服务器无外网，`docker pull postgres:16-alpine` 直接 timeout
2. 服务器无 Go 环境，`go mod download` 无法执行
3. 即使传源码上去，Dockerfile 里的多阶段构建也会因网络问题失败

**为什么不用二进制直接运行？**
- 还需要 PostgreSQL，单独安装离线 PG 更复杂
- Docker Compose 统一管理更干净

---

## 三、首次部署（一键脚本）

```bash
# 直接执行
./deploy.sh
```

脚本会自动完成以下步骤：

| 步骤 | 说明 | 耗时 |
|------|------|------|
| 1. 本地交叉编译 | `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build` | ~5s |
| 2. 构建 backend 镜像 | 基于 alpine，COPY 预编译二进制 | ~10s |
| 3. 拉取 postgres 镜像 | `docker pull --platform linux/amd64` | ~30s（视网络） |
| 4. 导出镜像 tar | `docker save -o wtb-images.tar` | ~10s |
| 5. 上传 tar + 配置 | `rsync` 到服务器 | ~30s-2min（视 tar 大小） |
| 6. 服务器加载运行 | `docker load` + `docker compose up` | ~20s |

**总计**：首次约 3-5 分钟，后续只传二进制约 10 秒。

---

## 四、配置文件说明

### 4.1 docker-compose.yml

- **无 caddy**：内网穿透直接访问 backend 端口
- **backend 端口映射**：`18080:8080`（内网穿透目标地址）
- **图片持久化**：`./images:/app/images`，上传的图片保存在宿主机 `~/wtb/images/`

### 4.2 Dockerfile

- **非多阶段构建**：直接使用本机编译好的 `backend-linux`（amd64）
- **基础镜像**：`alpine:latest`
- **内置图片**：`COPY images /app/images`（首次构建时嵌入）
- **运行时挂载**：同名目录被 volume 覆盖，上传的图片写入宿主机

---

## 五、踩坑记录（本次部署真实经验）

### ❌ 坑 1：服务器 docker pull timeout

**现象**：
```
failed to resolve reference "postgres:16-alpine":
dial tcp 103.228.130.61:443: i/o timeout
```

**根因**：服务器无外网，无法访问 Docker Hub。

**解决**：本机拉取镜像 → `docker save` → `rsync` 上传 → 服务器 `docker load`。

---

### ❌ 坑 2：本机构建时镜像源 401 Unauthorized

**现象**：
```
failed to solve: alpine:latest: unexpected status from HEAD request
to https://docker.m.daocloud.io/v2/library/alpine/manifests/latest: 401 Unauthorized
```

**根因**：Docker Desktop 配置了 daocloud 镜像加速，但 buildkit 解析元数据时认证失败。

**解决**：
```bash
# 关闭 buildkit，走传统构建
DOCKER_BUILDKIT=0 docker build ...

# 同时配置代理（本机 Clash）
HTTP_PROXY=http://127.0.0.1:7897 HTTPS_PROXY=http://127.0.0.1:7897 docker build ...
```

---

### ❌ 坑 3：平台架构不匹配（ARM64 vs AMD64）

**现象**：
```
image with reference sha256:xxx was found but its platform
(linux/arm64/v8) does not match the specified platform (linux/amd64)
```

**根因**：本机是 Apple Silicon（ARM64），交叉编译的二进制是 AMD64，但 Docker 默认拉取 ARM64 基础镜像。

**解决**：
```bash
# 显式指定平台拉取基础镜像
docker pull --platform linux/amd64 alpine:latest

# 构建时指定平台
docker build --platform linux/amd64 ...
```

> ⚠️ 注意：`docker build --platform` 需要配合 `--pull` 或预先 `docker pull --platform` 才能确保基础镜像正确。

---

### ❌ 坑 4：build context 路径错误

**现象**：
```
COPY failed: file not found in build context: stat backend-linux: file does not exist
```

**根因**：`Dockerfile` 里的 `COPY backend-linux .` 要求文件在 build context 根目录，但实际在 `backend/backend-linux`。

**解决**：构建前复制到根目录，或调整 Dockerfile 路径。

---

### ❌ 坑 5：图片目录不一致

**现象**：上传图片后容器重启丢失。

**根因**：docker-compose 里图片路径配置与实际代码期望不符。

**解决**：
- 代码读取 `IMAGE_PATH` 环境变量 → 设置为 `/app/images`
- docker-compose volume 挂载 `./images:/app/images`
- 上传图片写入 `/app/images`（实际落到宿主机 `./images/`）

---

## 六、日常更新（代码改动后）

如果只有 Go 代码改动，不需要重新传 300MB 的镜像 tar！

```bash
# 1. 本地交叉编译
cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o backend-linux .

# 2. 只传二进制（通常 20-30MB，秒级）
rsync -avz backend-linux qs_admin@192.168.0.156:~/wtb/

# 3. 热更新容器
ssh qs_admin@192.168.0.156 "
  docker cp ~/wtb/backend-linux wtb-backend:/app/backend-linux
  docker restart wtb-backend
  docker ps | grep wtb-backend
"

# 4. 验证
curl -s http://192.168.0.156:18080/health
```

---

## 七、目录结构（服务器）

```
/home/qs_admin/wtb/
├── docker-compose.yml      # 编排文件
├── Dockerfile              # 本地二进制版 Dockerfile
├── wtb-images.tar          # Docker 镜像包（首次部署后可选删除）
├── init.sql                # 数据库初始化
├── images/                 # 菜品图片（持久化目录）
│   ├── hongshao.png
│   └── ...
└── backend-linux           # 预编译二进制（热更新时使用）
```

---

## 八、验证清单

部署完成后，在服务器上执行：

```bash
curl -s -o /dev/null -w "health: %{http_code}\n" http://127.0.0.1:18080/health
curl -s -o /dev/null -w "categories: %{http_code}\n" http://127.0.0.1:18080/api/menu/categories
curl -s -o /dev/null -w "image: %{http_code}\n" http://127.0.0.1:18080/images/hongshao.png
```

期望全部返回 `200`。

---

## 九、内网穿透配置

在你的内网穿透工具中填写：

```
目标地址: 192.168.0.156:18080
```

无需在服务器配置任何域名或 SSL，内网穿透工具负责处理。

---

## 十、问题速查

| 现象 | 原因 | 解决 |
|------|------|------|
| `docker pull` timeout | 服务器无外网 | 本机构建 + `docker save` + `docker load` |
| `exec format error` | 二进制架构不匹配 | 本地编译时加 `GOARCH=amd64` |
| `401 Unauthorized` (build) | daocloud 镜像源认证失败 | `DOCKER_BUILDKIT=0` + 代理 |
| 图片 404 | images 目录缺失 | 检查 `~/wtb/images/` 是否存在 |
| 容器重启后图片丢失 | 未使用 volume 持久化 | 确认 docker-compose 有 `./images:/app/images` |
| 数据库连接失败 | postgres 未就绪 | backend 设置了 `depends_on condition: service_healthy`，稍等 10 秒 |

---

## 十一、配套前端部署（anzhitek 反向代理）

qs_admin 只跑了后端 API 容器，**后台管理界面前端（admin-web）部署在另一台服务器 anzhitek 上**。

### 架构关系

```
浏览器 → wtbadm.anzhitek.com → nginx(anzhitek) → 192.168.192.75:18080(qs_admin)
                │                     │
                │                     ├── / → 静态文件 (admin-web dist)
                │                     ├── /api → proxy_pass 后端
                │                     ├── /admin → proxy_pass 后端
                │                     └── /images → proxy_pass 后端
                │
                └── SSL 终止 + 域名解析
```

### 更新 admin-web 流程

```bash
# 1. 本地构建前端
cd admin-web
npm run build

# 2. 上传 dist 到 anzhitek（注意是宿主机路径 /home/ubuntu/html/）
rsync -avz --delete dist/ ubuntu@192.168.192.122:/home/ubuntu/html/wtbadm/

# 3. 重启 nginx 容器
ssh ubuntu@192.168.192.122 "docker restart nginx"
```

### 修改后端代理地址

如果 qs_admin IP 变更，在 anzhitek 上修改 nginx 配置：

```bash
ssh ubuntu@192.168.192.122
sed -i 's/旧IP:旧端口/192.168.192.75:18080/g' /home/ubuntu/nginx/conf.d/wtb.conf
docker exec nginx nginx -t && docker restart nginx
```

详细文档见 `deploy/anzhitek/README.md`。
