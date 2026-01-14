#FROM ubuntu:latest
#LABEL authors="Wzx"
#
#ENTRYPOINT ["top", "-b"]
FROM golang:alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建服务端
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./server/main.go && ls -la ./server

# 最终镜像
FROM alpine:latest

# 安装 ca-certificates
#RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建镜像复制二进制文件
COPY --from=builder /app/server .

# 暴露端口
EXPOSE 8888

# 运行服务
CMD ["./server"]