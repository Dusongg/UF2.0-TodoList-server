# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.22 AS builder

# 设置工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 复制到工作目录
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 将源代码复制到工作目录
COPY . .

# 构建 Go 应用
RUN go build -o myapp .

## 使用官方 Ubuntu 镜像作为基础镜像
#FROM ubuntu:22.04
# 使用官方 Debian 镜像作为运行阶段的基础镜像
FROM debian:stable-slim

#
## 安装 Redis 和 MySQL 客户端
#RUN apt-get update && apt-get install -y \
#    redis-tools \
#    mysql-client \
#    && rm -rf /var/lib/apt/lists/*
#
# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件到运行阶段
COPY --from=builder /app/myapp .

# 暴露应用使用的端口（根据需要修改）
EXPOSE 8080

# 启动应用
CMD ["./myapp"]
