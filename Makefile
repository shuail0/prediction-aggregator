.PHONY: help build run test clean install lint

help: ## 显示帮助信息
	@echo "Prediction Aggregator - Makefile Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## 安装依赖
	go mod download
	go mod tidy

build: ## 编译项目
	@echo "🔨 Building..."
	go build -o bin/aggregator cmd/aggregator/main.go
	@echo "✅ Build complete: bin/aggregator"

run: ## 运行程序
	@echo "🚀 Running aggregator..."
	go run cmd/aggregator/main.go

test: ## 运行测试
	@echo "🧪 Running tests..."
	go test -v -race -cover ./...

test-integration: ## 运行集成测试
	@echo "🧪 Running integration tests..."
	go test -v -tags=integration ./...

bench: ## 运行性能测试
	@echo "⚡ Running benchmarks..."
	go test -bench=. -benchmem ./pkg/orderbook

lint: ## 代码检查
	@echo "🔍 Linting code..."
	golangci-lint run ./...

fmt: ## 格式化代码
	@echo "✨ Formatting code..."
	go fmt ./...

clean: ## 清理构建文件
	@echo "🧹 Cleaning..."
	rm -rf bin/
	go clean

docker-build: ## 构建 Docker 镜像
	@echo "🐳 Building Docker image..."
	docker build -t prediction-aggregator:latest .

docker-run: ## 运行 Docker 容器
	@echo "🐳 Running Docker container..."
	docker run --env-file .env prediction-aggregator:latest

dev: ## 开发模式（热重载）
	@echo "🔥 Starting dev mode..."
	air

.DEFAULT_GOAL := help
