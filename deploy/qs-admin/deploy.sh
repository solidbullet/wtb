#!/bin/bash
set -e

TARGET_HOST="qs_admin@192.168.0.156"
REMOTE_DIR="/home/qs_admin/wtb"
LOCAL_DIR="$(cd "$(dirname "$0")/../.." && pwd)"

echo "========================================"
echo "  WTB 部署到 192.168.0.156 (qs_admin)"
echo "========================================"
echo ""
echo "  部署原理："
echo "  1. 本机交叉编译 backend-linux（amd64）"
echo "  2. 本机构建 Docker 镜像（利用外网拉取基础镜像）"
echo "  3. docker save 导出 postgres + backend 镜像"
echo "  4. rsync 上传镜像 tar + 配置到服务器"
echo "  5. 服务器 docker load + docker compose up"
echo ""

# 1. 本地交叉编译
echo "🔨 步骤 1/6：本地交叉编译 backend-linux..."
cd "$LOCAL_DIR/backend"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o backend-linux .
cd "$LOCAL_DIR"

# 2. 准备 build context（Dockerfile 要求文件在同级目录）
echo ""
echo "📦 步骤 2/6：准备 Docker build context..."
cp "$LOCAL_DIR/backend/backend-linux" "$LOCAL_DIR/backend-linux"
rsync -a "$LOCAL_DIR/miniprogram/images/" "$LOCAL_DIR/images/"

# 3. 构建镜像
echo ""
echo "🔨 步骤 3/6：本机构建 Docker 镜像（linux/amd64）..."
# 需要代理？取消下面注释并配置你的代理地址
# export HTTP_PROXY=http://127.0.0.1:7897
# export HTTPS_PROXY=http://127.0.0.1:7897

# 关闭 buildkit 避免某些镜像源 401 问题
DOCKER_BUILDKIT=0 docker build --platform linux/amd64 --pull -f "$LOCAL_DIR/deploy/qs-admin/Dockerfile" -t wtb-backend:deploy "$LOCAL_DIR"

# 4. 导出镜像
echo ""
echo "💾 步骤 4/6：导出镜像为 tar..."
docker save -o /tmp/wtb-images.tar postgres:16-alpine wtb-backend:deploy
ls -lh /tmp/wtb-images.tar

# 5. 上传文件
echo ""
echo "📤 步骤 5/6：上传文件到服务器..."
ssh "$TARGET_HOST" "mkdir -p $REMOTE_DIR/images"
rsync -avz --progress /tmp/wtb-images.tar "$TARGET_HOST:$REMOTE_DIR/"
rsync -avz "$LOCAL_DIR/deploy/qs-admin/docker-compose.yml" "$TARGET_HOST:$REMOTE_DIR/"
rsync -avz "$LOCAL_DIR/init.sql" "$TARGET_HOST:$REMOTE_DIR/"
rsync -avz "$LOCAL_DIR/images/" "$TARGET_HOST:$REMOTE_DIR/images/"

# 6. 服务器部署
echo ""
echo "🚀 步骤 6/6：服务器加载镜像并启动..."
ssh "$TARGET_HOST" "
  cd $REMOTE_DIR

  echo '  [服务器] 加载 Docker 镜像...'
  docker load -i wtb-images.tar

  echo '  [服务器] 停止旧容器...'
  docker compose down 2>/dev/null || true

  echo '  [服务器] 启动服务...'
  docker compose up -d

  echo '  [服务器] 等待就绪...'
  sleep 10

  echo '  [服务器] 容器状态:'
  docker ps | grep wtb || true
"

# 7. 验证
echo ""
echo "🧪 验证服务..."
ssh "$TARGET_HOST" "
  echo '  测试健康检查:'
  curl -s -o /dev/null -w '  health: %{http_code}\n' http://127.0.0.1:18080/health || echo '  health: failed'

  echo '  测试 API:'
  curl -s -o /dev/null -w '  categories: %{http_code}\n' http://127.0.0.1:18080/api/menu/categories || echo '  categories: failed'

  echo '  测试图片:'
  curl -s -o /dev/null -w '  hongshao.png: %{http_code}\n' http://127.0.0.1:18080/images/hongshao.png || echo '  hongshao.png: failed'
"

echo ""
echo "========================================"
echo "  ✅ 部署完成"
echo "========================================"
echo ""
echo "  后端地址: http://192.168.0.156:18080"
echo "  健康检查: http://192.168.0.156:18080/health"
echo ""
echo "  内网穿透目标地址: 192.168.0.156:18080"
echo ""
echo "  💡 提示：首次部署后，/tmp/wtb-images.tar 和本地 backend-linux"
echo "     可以删除以节省空间。日常更新只需重新编译二进制并 docker cp。"
echo ""
