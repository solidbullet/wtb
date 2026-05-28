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
