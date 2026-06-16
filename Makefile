# Media Station - 部署到NAS
# 使用方法: make deploy

.PHONY: help deploy build push clean

# 配置
NAS_HOST=192.168.5.178
NAS_USER=mellolo
REGISTRY_PORT=5000
APP_PORT=18080
PROJECT_NAME=media-station

# 帮助信息
help:
	@echo "用法:"
	@echo "  make deploy  - 一键部署到NAS"
	@echo "  make build   - 构建Docker镜像"
	@echo "  make push    - 推送镜像到NAS Registry"
	@echo "  make clean   - 清理本地构建缓存"

# 一键部署
deploy: build push
	@echo ""
	@echo "===================================="
	@echo "📦 部署到NAS"
	@echo "===================================="
	@ssh $(NAS_USER)@$(NAS_HOST) "\
		sudo docker pull $(NAS_HOST):$(REGISTRY_PORT)/$(PROJECT_NAME):latest && \
		sudo docker stop $(PROJECT_NAME) 2>/dev/null || true && \
		sudo docker rm $(PROJECT_NAME) 2>/dev/null || true && \
		sudo docker run -d -p $(APP_PORT):8080 --name $(PROJECT_NAME) --restart=always $(NAS_HOST):$(REGISTRY_PORT)/$(PROJECT_NAME):latest && \
		echo '' && \
		echo '✅ 部署成功！' && \
		echo '访问: http://$(NAS_HOST):$(APP_PORT)' \
	"
	@echo ""

# 构建镜像
build:
	@echo "===================================="
	@echo "🔨 构建Docker镜像"
	@echo "===================================="
	docker build -t $(NAS_HOST):$(REGISTRY_PORT)/$(PROJECT_NAME):latest .
	@echo "✅ 构建完成"
	@echo ""

# 推送镜像
push:
	@echo "===================================="
	@echo "📤 推送镜像到NAS"
	@echo "===================================="
	docker push $(NAS_HOST):$(REGISTRY_PORT)/$(PROJECT_NAME):latest
	@echo "✅ 推送完成"
	@echo ""

# 清理本地构建缓存
clean:
	@echo "清理本地Docker构建缓存..."
	docker system prune -f
	@echo "✅ 清理完成"
