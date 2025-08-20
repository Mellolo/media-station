# 构建阶段
FROM golang:latest AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 和 sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 运行阶段
FROM linuxserver/ffmpeg:latest

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件
RUN mkdir -p ./conf
COPY conf/app.prod.conf ./conf/app.conf

# 创建必要的目录
RUN mkdir -p ./logs ./cache

# 暴露端口
EXPOSE 8080

# 启动命令
ENTRYPOINT ["./main"]