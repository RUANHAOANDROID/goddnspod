
# 第一阶段：使用 golang 镜像作为基础镜像进行编译
FROM golang:1.20-alpine AS builder
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo '$TZ' > /etc/timezone
# 设置工作目录
WORKDIR /app

# 将当前目录的 go.mod 和 go.sum 复制到工作目录
COPY go.mod .
COPY go.sum .

# 下载依赖
RUN go mod download

# 将其他源代码文件复制到工作目录
COPY . .

# 编译 Go 应用，使用静态链接
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o runner

# 第二阶段：使用轻量级的 alpine 镜像作为基础镜像
FROM alpine:3.19.0 as runner

# 设置工作目录
WORKDIR /app

# 从第一阶段复制编译好的二进制文件到当前阶段
COPY --from=builder /app/runner .
COPY --from=builder /app/config.yml .

# 运行应用
CMD ["./runner"]
