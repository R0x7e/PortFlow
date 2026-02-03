# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制源代码
COPY . .

# 编译静态二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o portflow ./cmd/portflow

# Stage 2: Final Image
FROM alpine:latest

WORKDIR /app

# 复制二进制文件
COPY --from=builder /app/portflow .

# 暴露端口（示例，实际使用时可通过 -listen-port 指定）
EXPOSE 8080

# 启动命令
ENTRYPOINT ["./portflow"]
CMD ["--help"]
