# UOS Squid Exporter Makefile
.PHONY: build clean test install run fmt vet lint help

# 变量定义
BINARY_NAME=uos-squid-exporter
VERSION=1.0.0
BUILD_DIR=build
BIN_DIR=bin
GO_FILES=$(shell find . -name "*.go" -type f -not -path "./vendor/*")

# 默认目标
all: build

# 帮助信息
help:
	@echo "可用的 Makefile 目标:"
	@echo "  build     - 构建二进制文件"
	@echo "  clean     - 清理构建文件"
	@echo "  test      - 运行测试"
	@echo "  install   - 安装依赖"
	@echo "  run       - 运行程序"
	@echo "  fmt       - 格式化代码"
	@echo "  vet       - 代码静态分析"
	@echo "  lint      - 代码lint检查"
	@echo "  help      - 显示此帮助信息"

# 安装依赖
install:
	@echo "安装依赖..."
	go mod download
	go mod tidy

# 格式化代码
fmt:
	@echo "格式化 Go 代码..."
	go fmt ./...

# 代码静态分析
vet:
	@echo "运行 go vet..."
	go vet ./...

# 构建
build: fmt vet
	@echo "构建 $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BIN_DIR)/$(BINARY_NAME) .

# 构建发布版本
build-release:
	@echo "构建发布版本..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags "-X main.Version=$(VERSION) -w -s" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
		-ldflags "-X main.Version=$(VERSION) -w -s" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 运行测试覆盖率
test-coverage:
	@echo "运行测试覆盖率..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 运行程序
run:
	@echo "运行 $(BINARY_NAME)..."
	go run .

# 清理
clean:
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html

# 安装到系统
install-system: build
	@echo "安装到系统..."
	sudo cp $(BIN_DIR)/$(BINARY_NAME) /usr/local/bin/

# 创建服务文件
service:
	@echo "创建 systemd 服务文件..."
	@echo "[Unit]" > $(BINARY_NAME).service
	@echo "Description=UOS Squid Exporter" >> $(BINARY_NAME).service
	@echo "After=network.target" >> $(BINARY_NAME).service
	@echo "" >> $(BINARY_NAME).service
	@echo "[Service]" >> $(BINARY_NAME).service
	@echo "Type=simple" >> $(BINARY_NAME).service
	@echo "User=prometheus" >> $(BINARY_NAME).service
	@echo "Group=prometheus" >> $(BINARY_NAME).service
	@echo "ExecStart=/usr/local/bin/$(BINARY_NAME)" >> $(BINARY_NAME).service
	@echo "Restart=always" >> $(BINARY_NAME).service
	@echo "RestartSec=5" >> $(BINARY_NAME).service
	@echo "" >> $(BINARY_NAME).service
	@echo "[Install]" >> $(BINARY_NAME).service
	@echo "WantedBy=multi-user.target" >> $(BINARY_NAME).service

# 打包
package: build-release
	@echo "打包发布文件..."
	@mkdir -p $(BUILD_DIR)/package
	@cp $(BUILD_DIR)/$(BINARY_NAME)-* $(BUILD_DIR)/package/
	@cp README.md LICENSE $(BUILD_DIR)/package/
	@cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)-$(VERSION).tar.gz package/*
