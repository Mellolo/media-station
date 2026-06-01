#!/bin/bash

# 本地构建镜像并推送到NAS Registry
# 使用方法: ./build-and-push.sh

set -e

NAS_HOST="192.168.5.178"
REGISTRY_PORT="5000"
PROJECT_NAME="media-station"
IMAGE_NAME="${NAS_HOST}:${REGISTRY_PORT}/${PROJECT_NAME}:latest"

echo "======================================"
echo "  构建并推送镜像到Registry"
echo "======================================"
echo ""
echo "Registry: http://${NAS_HOST}:${REGISTRY_PORT}"
echo "镜像: ${IMAGE_NAME}"
echo ""

# 构建镜像
echo "构建镜像..."
docker build -t ${IMAGE_NAME} .
echo "✓ 构建完成"
echo ""

# 推送镜像
echo "推送镜像..."
docker push ${IMAGE_NAME}
echo "✓ 推送完成"
echo ""

echo "======================================"
echo "完成！镜像已推送到Registry"
echo "======================================"
echo ""
echo "查看镜像:"
echo "  curl http://${NAS_HOST}:${REGISTRY_PORT}/v2/_catalog"
echo ""
echo "下一步: 执行 ./deploy-on-nas.sh 部署应用"