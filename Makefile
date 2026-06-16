# Media Station - 本地Docker
# 使用方法: make build / make deploy

.PHONY: help build deploy clean stop logs

# 配置
IMAGE_NAME=media-station
TAG=latest
CONTAINER_NAME=media-station
APP_PORT=18080

# 帮助信息
help:
	@echo "用法:"
	@echo "  make build   - 构建Docker镜像"
	@echo "  make deploy  - 部署到本地Docker（停止旧容器+启动新容器）"
	@echo "  make stop    - 停止容器"
	@echo "  make logs    - 查看容器日志"
	@echo "  make clean   - 清理容器和镜像"

# 构建镜像
build:
	@echo "===================================="
	@echo "🔨 构建Docker镜像"
	@echo "===================================="
	docker build --build-arg GOPROXY=https://goproxy.cn,direct -t $(IMAGE_NAME):$(TAG) .
	@echo "✅ 构建完成"
	@echo ""

# 部署到本地Docker
deploy: build
	@echo "===================================="
	@echo "📦 部署到本地Docker"
	@echo "===================================="
	@echo "停止旧容器..."
	-docker stop $(CONTAINER_NAME) 2>/dev/null || true
	-docker rm $(CONTAINER_NAME) 2>/dev/null || true
	@echo "启动新容器..."
	docker run -d \
		-p $(APP_PORT):8080 \
		--name $(CONTAINER_NAME) \
		--restart=always \
		$(IMAGE_NAME):$(TAG)
	@echo ""
	@echo "===================================="
	@echo "✅ 部署成功！"
	@echo "===================================="
	@echo "访问: http://localhost:$(APP_PORT)"
	@echo ""

# 停止容器
stop:
	@echo "停止容器..."
	-docker stop $(CONTAINER_NAME) 2>/dev/null || true
	-docker rm $(CONTAINER_NAME) 2>/dev/null || true
	@echo "✅ 已停止"

# 查看日志
logs:
	docker logs -f $(CONTAINER_NAME)

# 清理容器和镜像
clean:
	@echo "清理容器和镜像..."
	-docker stop $(CONTAINER_NAME) 2>/dev/null || true
	-docker rm $(CONTAINER_NAME) 2>/dev/null || true
	-docker rmi $(IMAGE_NAME):$(TAG) 2>/dev/null || true
	@echo "✅ 清理完成"
