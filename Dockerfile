# 汪托帮后端服务 Dockerfile
# 多阶段构建，最终镜像体积小

# ========== 构建阶段 ==========
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 先复制依赖文件，利用 Docker 缓存层
COPY go.mod go.sum ./
RUN go mod download

# 复制全部源码
COPY . .

# 编译后端二进制（静态链接，无 CGO）
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN cd backend && go build -ldflags="-s -w" -o backend .

# ========== 运行阶段 ==========
FROM alpine:latest

# 安装 CA 证书（HTTPS 请求需要）和时区数据
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/backend/backend .

# 暴露服务端口
EXPOSE 8080

# 设置时区
ENV TZ=Asia/Shanghai

# 启动命令
CMD ["./backend"]
