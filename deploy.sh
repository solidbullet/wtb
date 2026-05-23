#!/bin/bash
set -e

echo "========================================"
echo "  汪托帮点餐系统 - Docker 部署脚本"
echo "  域名: wtb.lqqnw.cn"
echo "========================================"

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ 错误: Docker 未安装"
    echo "   请先安装 Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "❌ 错误: docker-compose 未安装"
    echo "   请先安装 docker-compose"
    exit 1
fi

# 判断使用 docker-compose 还是 docker compose
if docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
else
    COMPOSE_CMD="docker-compose"
fi

echo ""
echo "📦 步骤 1/5: 停止旧容器（如果存在）..."
$COMPOSE_CMD down 2>/dev/null || true

echo ""
echo "🔨 步骤 2/5: 构建后端镜像..."
$COMPOSE_CMD build --no-cache backend

echo ""
echo "🚀 步骤 3/5: 启动服务..."
$COMPOSE_CMD up -d

echo ""
echo "⏳ 步骤 4/5: 等待数据库初始化..."
sleep 8

echo ""
echo "🔒 步骤 5/5: 检查 HTTPS 证书..."
sleep 2

# 检查服务状态
echo ""
echo "📋 服务状态:"
$COMPOSE_CMD ps

echo ""
echo "========================================"
echo "  ✅ 部署完成！"
echo "========================================"
echo ""
echo "🌐 HTTP 地址:  http://wtb.lqqnw.cn"
echo "🔒 HTTPS 地址: https://wtb.lqqnw.cn  ← 推荐使用"
echo ""
echo "Caddy 会自动为 wtb.lqqnw.cn 申请 Let's Encrypt 证书"
echo "首次申请可能需要 10-30 秒，请稍等后访问 HTTPS"
echo ""
echo "常用命令:"
echo "  查看日志: $COMPOSE_CMD logs -f backend"
echo "  查看 Caddy: $COMPOSE_CMD logs -f caddy"
echo "  停止服务: $COMPOSE_CMD down"
echo "  重启服务: $COMPOSE_CMD restart backend"
echo ""
echo "⚠️  下一步:"
echo "  1. 将小程序 API_BASE 改为 https://wtb.lqqnw.cn"
echo "  2. 在微信公众平台配置服务器域名: https://wtb.lqqnw.cn"
echo "  3. 上传小程序代码"
echo ""
