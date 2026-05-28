# WTB 生产环境部署指南（hooper）

> **目标服务器**：`hooper@10.144.144.1`（EasyTier 内网）  
> **域名**：`wtb.lqqnw.cn` / `wtbadm.lqqnw.cn`  
> **哲学**：本地编译 → 只传二进制 → 热更新。服务器无 Go 缓存，`docker build` 每次重新下载依赖极慢。

---

## 一、架构总览

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

## 二、服务器环境准备（首次）

```bash
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

---

## 三、文件传输

> ⚠️ **重要**：公网 IP `118.145.193.23` 带宽小、易超时。务必使用 EasyTier 虚拟内网 `10.144.144.1` 传输。

| 本地路径 | 服务器路径 | 说明 |
|----------|-----------|------|
| `backend/Dockerfile` | `~/wtb/Dockerfile` | 后端镜像构建文件 |
| `backend/router.go` 等源码 | `~/wtb/` 对应目录 | 后端源码（用于 Docker build） |
| `docker-compose.yml` | `~/wtb/docker-compose.yml` | 编排文件 |
| `init.sql` | `~/wtb/init.sql` | 数据库初始化 |
| `miniprogram/images/*.png` | `/var/www/wtb/images/` | 菜品图片（必须传！） |
| `admin-web/dist/` | `/var/www/wtbadm/` | 管理后台前端（如有更新） |

**图片传输示例**：

```bash
# 1. 本地压缩图片
tar czf /tmp/images.tar.gz -C miniprogram/images .

# 2. 通过 EasyTier 内网传
rsync -avz -e "ssh" /tmp/images.tar.gz hooper@10.144.144.1:/tmp/

# 3. 服务器上解压到统一目录
ssh hooper@10.144.144.1 "sudo tar xzf /tmp/images.tar.gz -C /var/www/wtb/images/"
```

---

## 四、配置文件

### 4.1 docker-compose.yml

见本目录 `docker-compose.yml`。核心要点：
- 无 caddy（使用系统 nginx）
- backend 端口映射 `127.0.0.1:18080:8080`
- 图片 volume：`/var/www/wtb/images:/miniprogram/images`

### 4.2 nginx 配置

见本目录 `nginx-wtb.conf` 和 `nginx-wtbadm.conf`，分别放置到：
- `/etc/nginx/sites-enabled/wtb.lqqnw.cn`
- `/etc/nginx/sites-enabled/wtbadm`

> ⚠️ **关键**：`wtbadm` 配置中 `location ^~ /images/` 必须在正则 `location ~* \.(png|...)$` **之前**。

---

## 五、部署方案

### 5.1 超快速部署（日常代码改动，推荐）

> 全程 10-20 秒。服务器上没有 Go 模块缓存，`docker build` 每次都会重新下载依赖，非常慢。

见本目录 `fast-deploy.sh`：

```bash
# 本地交叉编译 → 只传二进制 → docker cp → 重启
./fast-deploy.sh
```

**适用场景**：
- ✅ 后端 Go 代码有改动（API、业务逻辑等）
- ✅ 只改动了一两个文件
- ✅ 需要秒级上线

**不适用场景**：
- ❌ Dockerfile 本身有改动
- ❌ 新增/修改了静态资源且 Dockerfile 的 `COPY` 层需要更新
- ❌ 依赖了 `miniprogram/images` 目录的新增图片（需要重新 build 镜像）

### 5.2 完整重建部署（首次/镜像更新）

见本目录 `full-deploy.sh`：

```bash
# 强制检查 → docker compose build → 启动
./full-deploy.sh
```

> ⚠️ 服务器上 `go build` 下载依赖很慢（10-30 分钟），日常小改动请用「超快速部署」。

### 5.3 更新 admin-web 前端

```bash
# 本地构建
cd admin-web
npm run build

# 上传到服务器（通过 EasyTier 内网）
rsync -avz --delete -e "ssh" dist/ hooper@10.144.144.1:/var/www/wtbadm/

# 不需要重启任何服务，nginx 直接服务静态文件
```

---

## 六、部署前强制检查

见本目录 `pre-deploy-check.sh`。每次部署前执行：

```bash
chmod +x pre-deploy-check.sh
./pre-deploy-check.sh
```

不通过禁止部署。

---

## 七、回滚操作

```bash
cd /home/hooper/wtb
docker stop wtb-backend
docker rm wtb-backend
docker compose up -d backend
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
