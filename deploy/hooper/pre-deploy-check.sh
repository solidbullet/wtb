#!/bin/bash
set -e

PASS=0
FAIL=0

check() {
    if eval "$2" > /dev/null 2>&1; then
        echo "  ✅ PASS: $1"
        PASS=$((PASS+1))
    else
        echo "  ❌ FAIL: $1"
        FAIL=$((FAIL+1))
    fi
}

echo "========================================"
echo "  WTB 部署前强制检查清单"
echo "========================================"

echo ""
echo "【文件检查】"
check "init.sql 是普通文件（不是目录）" "[ -f /home/hooper/wtb/init.sql ]"
check "docker-compose.yml 存在" "[ -f /home/hooper/wtb/docker-compose.yml ]"
check "后端 Dockerfile 存在" "[ -f /home/hooper/wtb/Dockerfile ]"

echo ""
echo "【图片检查】"
check "nginx 图片目录存在" "[ -d /var/www/wtb/images ]"
check "图片目录有菜品图片" "ls /var/www/wtb/images/*.png > /dev/null 2>&1"
check "hongshao.png 存在" "[ -f /var/www/wtb/images/hongshao.png ]"

echo ""
echo "【配置检查】"
check "docker-compose.yml 格式正确" "cd /home/hooper/wtb && docker compose config > /dev/null 2>&1"
check "nginx 语法正确" "sudo nginx -t > /dev/null 2>&1"
check "wtbadm nginx 配置有 ^~ /images/" "grep -q 'location ^~ /images/' /etc/nginx/sites-enabled/wtbadm"
check "wtbadm nginx 配置 /images/ 在正则之前" "bash -c 'IMG_LINE=\$(grep -n \"location ^~ /images/\" /etc/nginx/sites-enabled/wtbadm | head -1 | cut -d: -f1); REG_LINE=\$(grep -n \"location ~*\" /etc/nginx/sites-enabled/wtbadm | head -1 | cut -d: -f1); [ -n \"\$IMG_LINE\" ] && [ -n \"\$REG_LINE\" ] && [ \"\$IMG_LINE\" -lt \"\$REG_LINE\" ]'"

echo ""
echo "【容器检查】"
check "postgres 容器在运行" "docker ps | grep -q wtb-postgres"
check "backend 容器在运行" "docker ps | grep -q wtb-backend"

echo ""
echo "【网络检查】"
check "wtb.lqqnw.cn API 可访问" "curl -s -o /dev/null -w '%{http_code}' https://wtb.lqqnw.cn/api/menu/categories | grep -q '200'"
check "wtb.lqqnw.cn 图片可访问" "curl -s -o /dev/null -w '%{http_code}' https://wtb.lqqnw.cn/images/hongshao.png | grep -q '200'"
check "wtbadm.lqqnw.cn 图片可访问" "curl -s -o /dev/null -w '%{http_code}' https://wtbadm.lqqnw.cn/images/hongshao.png | grep -q '200'"

echo ""
echo "========================================"
echo "  结果: $PASS 通过, $FAIL 失败"
echo "========================================"

if [ $FAIL -gt 0 ]; then
    echo "❌ 有检查项未通过，禁止部署！先修复上述问题。"
    exit 1
else
    echo "✅ 全部通过，可以安全部署。"
fi
