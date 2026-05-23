# 本地开发环境文档

> 生成日期：2026年
> 适用场景：Go 后台 + 微信小程序 + PostgreSQL 开发

---

## 系统信息

| 项目 | 值 |
|------|-----|
| **操作系统** | macOS 26.3 (25D125) |
| **内核** | Darwin 25.3.0 (xnu-12377.81.4~5/RELEASE_ARM64_T6020) |
| **架构** | arm64 (Apple Silicon) |
| **主机名** | QSs-Mac-mini |
| **工作区** | /Users/admin/workspace/jyq/dp |
| **Shell** | zsh |

---

## 工具链版本

### 后端开发

| 工具 | 版本 | 安装方式 | 路径 |
|------|------|----------|------|
| **Go** | 1.24.1 | Homebrew | /opt/homebrew/Cellar/go/1.24.1/libexec |
| **GOPATH** | /Users/admin/go | — | — |
| **GNU Make** | 3.81 | 系统自带 | /usr/bin/make |

### 数据库

| 工具 | 版本 | 安装方式 | 备注 |
|------|------|----------|------|
| **PostgreSQL** | 18.3 | Homebrew | 本地运行中（主数据库） |
| **MySQL Server** | 5.7.24 | conda-forge | 本地运行中（遗留） |

### 前端开发

| 工具 | 版本 | 安装方式 |
|------|------|----------|
| **Node.js** | 22.16.0 | — |
| **npm** | 10.9.2 | — |
| **微信开发者工具** | 已安装 | /Applications/wechatwebdevtools.app |

### 容器

| 工具 | 版本 | 状态 |
|------|------|------|
| **Docker** | 28.4.0 | 已安装，daemon 未启动 |
| **Docker Compose** | v2.39.4-desktop.1 | — |

### 工具链

| 工具 | 版本 |
|------|------|
| **Git** | 2.48.1 |
| **Homebrew** | 5.1.9 |
| **VS Code** | 已安装 (/usr/local/bin/code) |
| **curl** | 8.2.1 |

---

## 数据库连接信息

### PostgreSQL（主数据库）

| 项目 | 值 |
|------|-----|
| **地址** | /tmp (Unix socket) 或 localhost:5432 |
| **用户** | admin |
| **密码** | 无（peer 认证） |
| **连接命令** | `psql -h /tmp -U admin -d postgres` |

### GORM 连接串 (DSN)

```
host=/tmp user=admin dbname=wtb_user sslmode=disable TimeZone=Asia/Shanghai
```

> 每个微服务替换 `dbname=wtb_user` 为对应的数据库名。

### MySQL（遗留）

| 项目 | 值 |
|------|-----|
| **地址** | 127.0.0.1:3306 |
| **用户** | root |
| **密码** | `dev123456` |

---

## 已有数据库

### PostgreSQL

- postgres
- template0
- template1
- testdb

### MySQL

- information_schema
- mysql
- performance_schema
- sys
- wtb_ordering（点餐系统遗留库）

---

## 缺失但可能需要安装的工具

| 工具 | 用途 | 安装命令 |
|------|------|----------|
| **air** | Go 热重载开发 | `go install github.com/air-verse/air@latest` |
| **Redis** | 缓存/购物车/分布式锁 | `brew install redis` |

---

## 环境变量建议

在 `~/.zshrc` 中确认或添加：

```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

---

## 快速验证命令

```bash
# Go
go version

# PostgreSQL
psql -h /tmp -U admin -d postgres -c "SELECT VERSION();"

# MySQL（遗留）
mysql -u root -pdev123456 -e "SELECT VERSION();"

# Node
node --version && npm --version

# Git
git --version

# Docker (需要先启动 Docker Desktop)
docker --version
```

---

## AI 编程工具使用指引

将此文档提供给 AI 编程助手（如 DeepSeek TUI、Cursor、Claude Code 等），它们即可直接使用上述信息进行开发，无需反复确认环境配置。

关键信息速查卡片：

```
系统:       macOS 26.3 arm64   工作区: /Users/admin/workspace/jyq/dp
Go:         1.24.1             PostgreSQL: 18.3 (/tmp, admin, 无密码)
Node:       22.16.0            npm:        10.9.2
Git:        2.48.1             Docker:     28.4.0 (未启动)
```
