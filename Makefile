.PHONY: help run-backend test-backend

help:
	@echo "make run-backend   # 启动单体后端服务"
	@echo "make test-backend  # 测试后端及共享包"

run-backend:
	cd backend && go run main.go

test-backend:
	cd pkg/jwt && go test ./... -v
	cd pkg/response && go test ./... -v
	cd backend && go test ./... -v
