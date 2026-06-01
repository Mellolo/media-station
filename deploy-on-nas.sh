#!/bin/bash

# 在NAS上拉取镜像并运行容器（直接在NAS服务器上执行）
# 使用方法: 
#   1. SSH连接到NAS: ssh mellolo@192.168.5.178
#   2. 执行此脚本: ./deploy-on-nas.sh

NAS_HOST="192.168.5.178"
REGISTRY_PORT="5000"
APP_PORT="18080"
PROJECT_NAME="media-station"
IMAGE_NAME="${NAS_HOST}:${REGISTRY_PORT}/${PROJECT_NAME}:latest"

echo "======================================"
echo "  在NAS上部署应用"
echo "======================================"
echo ""
echo "镜像: ${IMAGE_NAME}"
echo "端口: ${APP_PORT}:8080"
echo ""

# 拉取镜像
echo "拉取镜像..."
docker pull ${IMAGE_NAME}
echo "✓ 拉取完成"
echo ""

# 停止并删除旧容器
echo "停止旧容器..."
docker stop ${PROJECT_NAME} || true
docker rm ${PROJECT_NAME} || true
echo "✓ 清理完成"
echo ""

# 启动新容器
echo "启动新容器..."
docker run -d \
  -p ${APP_PORT}:8080 \
  --name ${PROJECT_NAME} \
  --restart=always \
  ${IMAGE_NAME}
echo "✓ 启动完成"
echo ""

# 验证部署
echo "验证部署..."
docker ps | grep ${PROJECT_NAME}
echo ""

echo "======================================"
echo "  部署完成"
echo "======================================"
echo ""
echo "访问: http://${NAS_HOST}:${APP_PORT}"
echo ""
echo "查看日志:"
echo "  docker logs -f ${PROJECT_NAME}"
echo ""
echo "管理命令:"
echo "  停止: docker stop ${PROJECT_NAME}"
echo "  重启: docker restart ${PROJECT_NAME}"
echo "  删除: docker rm -f ${PROJECT_NAME}"
