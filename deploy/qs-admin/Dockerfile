# 汪托帮后端服务 Dockerfile（本地编译版）
# 直接使用本机交叉编译好的二进制，避免服务器下载 Go 依赖
# 适用场景：目标服务器无外网，无法 docker pull / go mod download

FROM alpine:latest

# 安装 CA 证书（HTTPS 请求需要）和时区数据
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 复制本地编译好的二进制文件（要求与 Dockerfile 同目录）
COPY backend-linux .

# 复制初始图片资源
COPY images /app/images

# 暴露服务端口
EXPOSE 8080

# 设置时区
ENV TZ=Asia/Shanghai

# 启动命令
CMD ["./backend-linux"]
