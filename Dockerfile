# 编译 web
FROM node:18 as frontend-build

WORKDIR /app/frontend

COPY frontend/package*.json ./frontend

RUN npm install

COPY frontend/ ./frontend

RUN npm run build

# 编译go
FROM golang:1.20-alpine AS builder
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo '$TZ' > /etc/timezone
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
# 安装UPX
RUN apk add --no-cache upx
# 将其他源代码文件复制到工作目录
COPY . .
# 编译 Go 应用，使用静态链接
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o runner
# 使用UPX进行压缩
RUN upx --best --lzma runner
# 第二阶段：使用轻量级的 alpine 镜像作为基础镜像
FROM alpine:3.20.3 as runner

# 设置工作目录
WORKDIR /app

# 从第一阶段复制编译好的二进制文件到当前阶段
COPY --from=builder /app/runner .
COPY --from=builder /app/config.yml .
COPY --from=builder /app/frontend/build/client ./frontend/build/client

# 暴露端口 6565
EXPOSE 6565

# 运行应用
CMD ["./runner"]
