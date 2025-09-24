# 项目信息
PROJECT_NAME := uos-exporter
NAME := squid-exporter
VERSION := 1.0.0
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
SRC_DIR := .
MODULES_DIR := modules

# 编译选项
GO := go
GO_BUILD := $(GO) build
GO_CLEAN := $(GO) clean
GO_MOD := $(GO) mod
GO_FLAGS := -ldflags="-s -w -X main.Version=$(VERSION)"

# 安装路径（支持自定义，默认为 /usr/local/bin）
PREFIX ?= /usr/local
INSTALL_BIN_DIR := $(PREFIX)/bin
INSTALL_UNIT_DIR := /usr/lib/systemd/system
INSTALL_CONFIG_DIR := /etc/uos-exporter

# 默认目标
all: build

# 创建构建目录
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

before_build:
	@echo "Installing dependencies..."
	$(GO_MOD) tidy

# 构建主程序
build: $(BIN_DIR) before_build
	@echo "Building $(NAME)..."
	$(GO_BUILD) $(GO_FLAGS) -o $(BIN_DIR)/$(NAME) $(SRC_DIR)

# 安装到指定路径
install: build
	@echo "Installing $(NAME) to $(DESTDIR)$(INSTALL_BIN_DIR)..."
	install -D -p -m 0755 $(BIN_DIR)/$(NAME) $(DESTDIR)$(INSTALL_BIN_DIR)/$(NAME)
	install -D -p -m 0644 $(SRC_DIR)/$(NAME).service $(DESTDIR)$(INSTALL_UNIT_DIR)/uos-$(NAME).service
	install -D -p -m 0644 $(SRC_DIR)/config/$(NAME).yaml $(DESTDIR)$(INSTALL_CONFIG_DIR)/$(NAME).yaml

# 清理生成的文件
clean:
	@echo "Cleaning up..."
	$(GO_CLEAN)
	rm -rf $(BUILD_DIR)


# 提供帮助信息
help:
	@echo "Usage: make [target]"
	@echo "Available targets:"
	@echo "  all       - Build the project (default)"
	@echo "  build     - Compile the main program"
	@echo "  install   - Install the binary to /usr/local/bin"
	@echo "  clean     - Remove generated files"
	@echo "  help      - Show this help message"

.PHONY: all build install clean help