# 第一阶段：构建环境
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目代码
COPY . .

# 编译项目
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-connect .

# 第二阶段：运行环境
FROM alpine:latest

# 安装必要的依赖
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/go-connect .

# 复制配置文件
COPY config/app.yaml config/

# 暴露端口
EXPOSE 9000

# 设置环境变量
ENV GIN_MODE=release

# 启动命令
CMD ["./go-connect"]