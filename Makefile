.PHONY: help db-create run-user run-gateway test-user test-all

help:
	@echo "make db-create     # 创建所有数据库"
	@echo "make run-user      # 启动用户服务"
	@echo "make run-gateway   # 启动网关"
	@echo "make test-user     # 测试用户服务"
	@echo "make test-all      # 测试所有服务"

db-create:
	@for db in user seat menu order payment points activity pricing analytics; do \
		psql -h /tmp -U admin -d postgres -c "CREATE DATABASE wtb_$${db} OWNER admin ENCODING 'UTF8';" 2>/dev/null; \
		echo "created wtb_$${db}"; \
	done

run-user:
	cd services/user && go run main.go

run-gateway:
	cd gateway && go run main.go

test-user:
	cd services/user && go test ./... -v

test-all:
	cd pkg/jwt && go test ./... -v
	cd pkg/response && go test ./... -v
	cd pkg/httpclient && go test ./... -v
	@for svc in user seat menu order payment points activity pricing analytics admin; do \
		echo "=== testing services/$$svc ==="; \
		cd services/$$svc && go test ./... -v && cd ../..; \
	done
	cd gateway && go test ./... -v
