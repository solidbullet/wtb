#!/bin/bash
# 启动管理后台所需的所有服务

set -e

echo "=== 启动管理后台服务 ==="

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 记录所有 PID
PID_FILE="/tmp/wtb_pids.txt"
> "$PID_FILE"

# 启动业务微服务（后台运行）
services=(
  "user:8081"
  "seat:8082"
  "menu:8083"
  "order:8084"
  "payment:8085"
  "points:8086"
  "activity:8087"
  "pricing:8088"
  "analytics:8089"
)

for svc in "${services[@]}"; do
  name=${svc%%:*}
  port=${svc##*:}
  echo -n "启动 services/$name (:$port) ... "
  cd "services/$name"
  go run main.go > "/tmp/wtb_$name.log" 2>&1 &
  echo $! >> "$PID_FILE"
  cd ../..
  echo "PID $!"
done

# 启动 admin BFF 服务
echo -n "启动 services/admin (:8090) ... "
cd services/admin
go run main.go > /tmp/wtb_admin.log 2>&1 &
echo $! >> "$PID_FILE"
echo "PID $!"
cd ../..

# 启动 gateway
echo -n "启动 gateway (:8080) ... "
cd gateway
go run main.go > /tmp/wtb_gateway.log 2>&1 &
echo $! >> "$PID_FILE"
echo "PID $!"
cd ..

echo ""
echo "等待服务就绪（最多30秒）..."

ports=(8081 8082 8083 8084 8085 8086 8087 8088 8089 8090 8080)
for port in "${ports[@]}"; do
  for i in {1..30}; do
    if lsof -Pi :$port -sTCP:LISTEN > /dev/null 2>&1; then
      echo -e "${GREEN}✓${NC} 端口 $port 就绪"
      break
    fi
    sleep 1
    if [ $i -eq 30 ]; then
      echo -e "${RED}✗${NC} 端口 $port 未就绪，请检查日志 /tmp/wtb_*.log"
    fi
  done
done

echo ""
echo "启动前端静态服务器 (Vite React)..."
cd admin-web
npm run dev -- --port 3000 > /tmp/wtb_frontend.log 2>&1 &
echo $! >> "$PID_FILE"
echo "前端 PID $! 运行于 http://localhost:3000 (Vite)"
cd ..

echo ""
echo "============================================"
echo "管理后台访问地址: http://localhost:3000"
echo "Admin API 地址: http://localhost:8090"
echo "Gateway 地址:  http://localhost:8080"
echo ""
echo "登录账号: admin"
echo "登录密码: admin123"
echo "============================================"
echo ""
echo "查看日志: tail -f /tmp/wtb_*.log"
echo "停止所有服务: kill \$(cat $PID_FILE)"
