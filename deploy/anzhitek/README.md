# WTB 前端与反向代理部署指南（anzhitek）

> **目标服务器**：`ubuntu@192.168.192.122`  
> **服务角色**：nginx 反向代理 + admin-web 静态文件服务  
> **域名**：`wtbadm.anzhitek.com` / `wtb.anzhitek.com`

---

## 一、环境特点

| 项目 | 状态 | 说明 |
|------|------|------|
| nginx | ✅ Docker 容器运行 | 非系统 nginx，通过 docker-proxy 监听 80/443 |
| SSL 证书 | ✅ Let's Encrypt | `/etc/letsencrypt/live/wtb.anzhitek.com/` |
| 静态文件 | ✅ 已配置 | `/home/ubuntu/html/wtbadm/` 对应域名根目录 |
| 后端代理 | ✅ 已配置 | 指向 `192.168.192.75:18080`（qs_admin） |

---

## 二、架构关系

```
用户浏览器
     │
     ▼
wtbadm.anzhitek.com / wtb.anzhitek.com
     │
     ▼
┌────────────────────────────────────────────┐
│  anzhitek (192.168.192.122)                │
│  nginx Docker 容器                          │
│                                              │
│  ┌──────────────────────────────────────┐  │
│  │  wtbadm.anzhitek.com                 │  │
│  │  ├── /  → 静态文件 (/home/html/wtbadm)│  │
│  │  ├── /api → proxy_pass 后端          │  │
│  │  ├── /admin → proxy_pass 后端        │  │
│  │  └── /images → proxy_pass 后端       │  │
│  └──────────────────────────────────────┘  │
│                                              │
│  ┌──────────────────────────────────────┐  │
│  │  wtb.anzhitek.com                    │  │
│  │  └── / → proxy_pass 后端             │  │
│  └──────────────────────────────────────┘  │
└────────────────────────────────────────────┘
     │
     ▼
192.168.192.75:18080 (qs_admin backend)
```

---

## 三、目录映射（宿主机 ↔ 容器）

| 宿主机路径 | 容器内路径 | 用途 |
|-----------|-----------|------|
| `/home/ubuntu/html` | `/home/html` | 静态文件（admin-web、首页等） |
| `/home/ubuntu/nginx/conf.d` | `/etc/nginx/conf.d` | nginx 站点配置 |
| `/etc/letsencrypt` | `/etc/letsencrypt` | SSL 证书（只读） |

---

## 四、admin-web 更新部署

### 4.1 构建前端

```bash
cd admin-web
npm run build
```

### 4.2 上传到 anzhitek

```bash
rsync -avz --delete dist/ ubuntu@192.168.192.122:/home/ubuntu/html/wtbadm/
```

> ⚠️ 注意：nginx 容器挂载的是 `/home/ubuntu/html`，不是 `/home/html`。上传时必须写到宿主机的 `/home/ubuntu/html/wtbadm/`。

### 4.3 重启 nginx

```bash
ssh ubuntu@192.168.192.122 "docker restart nginx"
```

---

## 五、修改后端代理地址

如果后端服务器 IP 或端口变更，需要修改 nginx 配置：

```bash
# 编辑配置文件
ssh ubuntu@192.168.192.122
sudo vim /home/ubuntu/nginx/conf.d/wtb.conf

# 修改 proxy_pass 行，例如：
proxy_pass http://192.168.192.75:18080;

# 验证配置语法并重启
docker exec nginx nginx -t
docker restart nginx
```

---

## 六、踩坑记录

### ❌ 坑 1：找不到 nginx 配置文件

**现象**：`/etc/nginx/` 目录不存在，`which nginx` 无结果。

**根因**：nginx 运行在 Docker 容器内，系统层面没有安装 nginx。

**解决**：
```bash
docker ps | grep nginx          # 确认容器在运行
docker inspect nginx | grep Mounts  # 查看 volume 映射
```

---

### ❌ 坑 2：rsync 目标路径错误

**现象**：
```
mkdir "/home/html/wtbadm" failed: No such file or directory
```

**根因**：误把容器内的路径 `/home/html/wtbadm` 当作 rsync 目标，但 ssh 登录后的文件系统看到的是宿主机路径。

**解决**：必须使用宿主机的绝对路径：
```bash
rsync -avz --delete dist/ ubuntu@192.168.192.122:/home/ubuntu/html/wtbadm/
```

---

### ❌ 坑 3：HTTP 301 重定向到 HTTPS

**现象**：curl 测试返回 `301`。

**根因**：nginx 配置中 `server_name wtbadm.anzhitek.com` 的 80 端口配置了 `return 301 https://$host$request_uri;`。

**解决**：测试时必须使用 HTTPS：
```bash
curl -sk -H 'Host: wtbadm.anzhitek.com' https://127.0.0.1/index.html
```

---

## 七、验证清单

```bash
# 在 anzhitek 服务器上执行

# 1. nginx 配置语法
docker exec nginx nginx -t

# 2. 静态文件访问
curl -sk -o /dev/null -w "index.html: %{http_code}\n" \
  -H 'Host: wtbadm.anzhitek.com' https://127.0.0.1/index.html

# 3. API 代理
curl -sk -o /dev/null -w "health: %{http_code}\n" \
  -H 'Host: wtbadm.anzhitek.com' https://127.0.0.1/health

# 4. 图片代理
curl -sk -o /dev/null -w "image: %{http_code}\n" \
  -H 'Host: wtbadm.anzhitek.com' https://127.0.0.1/images/hongshao.png
```

---

## 八、相关配置参考

见本目录 `nginx-wtb.conf`，为当前线上运行的配置备份。
