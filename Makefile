.PHONY: help build run test clean install lint

BIN_DIR := bin

help: ## 显示帮助信息
	@echo "Prediction Aggregator - Makefile Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## 安装依赖
	go mod download
	go mod tidy

build: ## 编译所有程序
	@echo "🔨 Building..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/aggregator cmd/aggregator/main.go
	go build -o $(BIN_DIR)/scanner cmd/scanner/main.go
	go build -o $(BIN_DIR)/maker cmd/maker/main.go
	@echo "✅ Build complete: $(BIN_DIR)/"

build-aggregator: ## 编译 aggregator
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/aggregator cmd/aggregator/main.go

build-scanner: ## 编译 scanner
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/scanner cmd/scanner/main.go

build-maker: ## 编译 maker
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/maker cmd/maker/main.go

run: ## 运行 aggregator
	go run cmd/aggregator/main.go

run-scanner: ## 运行 scanner
	go run cmd/scanner/main.go

run-maker: ## 运行 maker
	go run cmd/maker/main.go

test: ## 运行测试
	@echo "🧪 Running tests..."
	go test -v -race -cover ./...

test-integration: ## 运行集成测试
	@echo "🧪 Running integration tests..."
	go test -v -tags=integration ./...

lint: ## 代码检查
	@echo "🔍 Linting code..."
	golangci-lint run ./...

fmt: ## 格式化代码
	@echo "✨ Formatting code..."
	go fmt ./...

clean: ## 清理构建文件
	@echo "🧹 Cleaning..."
	rm -rf $(BIN_DIR)/
	go clean

docker-build: ## 构建 Docker 镜像
	@echo "🐳 Building Docker image..."
	docker build -t prediction-aggregator:latest .

docker-run: ## 运行 Docker 容器
	@echo "🐳 Running Docker container..."
	docker run --env-file .env prediction-aggregator:latest

.DEFAULT_GOAL := help
