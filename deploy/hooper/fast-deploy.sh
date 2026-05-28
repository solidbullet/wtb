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
