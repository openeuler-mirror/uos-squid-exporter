# SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
# SPDX-License-Identifier: MIT

# 多阶段构建 - 构建阶段
FROM golang:1.22.8-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w -X main.Version=1.0.0" -o squid-exporter .

# 运行阶段
FROM alpine:latest

# 安装必要的包
RUN apk --no-cache add ca-certificates tzdata

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/squid-exporter .

# 复制配置文件
COPY --from=builder /app/config/squid-exporter.yaml ./config/

# 创建日志目录
RUN mkdir -p /var/log/uos-exporter && \
    chown -R appuser:appgroup /var/log/uos-exporter

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 9109

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9109/healthz || exit 1

# 启动命令
CMD ["./squid-exporter"]
