.PHONY: all proto build run clean test

# 默认目标
all: proto build

# 生成proto代码
proto:
	@echo "生成proto代码..."
	cd proto && $(MAKE)

# 编译项目
build: proto
	@echo "编译项目..."
	go build -o elysia-llm-tool main.go

# 运行服务
run: build
	@echo "启动服务..."
	./elysia-llm-tool

# 开发模式运行（不编译）
dev:
	@echo "开发模式启动..."
	go run main.go

# 清理生成的文件
clean:
	@echo "清理文件..."
	rm -f elysia-llm-tool
	cd proto && $(MAKE) clean
	rm -rf logs/

# 安装依赖
deps:
	@echo "安装依赖..."
	go mod download
	go mod tidy

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run

# 帮助信息
help:
	@echo "可用命令："
	@echo "  make proto  - 生成proto代码"
	@echo "  make build  - 编译项目"
	@echo "  make run    - 编译并运行服务"
	@echo "  make dev    - 开发模式运行（不编译）"
	@echo "  make clean  - 清理生成的文件"
	@echo "  make deps   - 安装依赖"
	@echo "  make test   - 运行测试"
	@echo "  make fmt    - 格式化代码"
	@echo "  make lint   - 代码检查"
