# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装必要的系统依赖
RUN apk --no-cache add ca-certificates tzdata git

# 先拷贝依赖文件利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 拷贝源代码
COPY . .

# 编译（静态链接，最小化镜像）
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server main.go

# ─────────────────────────────────
# Runtime stage
FROM alpine:3.19

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从 build stage 拷贝可执行文件和迁移脚本
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

# 非特权用户运行
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

EXPOSE 8092

# 健康检查（依赖 TASK-09 的 /health 端点）
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8092/health || exit 1

CMD ["./server"]
