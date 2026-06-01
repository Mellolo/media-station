#!/bin/bash

# 一键部署脚本
# 使用方法: ./deploy.sh

set -e

PROJECT_NAME="media-station"
IMAGE_NAME="media-station:latest"

echo "开始部署 ${PROJECT_NAME}..."

# 构建Docker镜像
echo "1. 构建Docker镜像..."
docker build -t ${IMAGE_NAME} .

# 停止并删除旧容器
echo "2. 停止并删除旧容器..."
docker stop ${PROJECT_NAME} || true
docker rm ${PROJECT_NAME} || true

# 启动新容器
echo "3. 启动新容器..."
docker run -d \
  -p 18080:8080 \
  --name ${PROJECT_NAME} \
  --restart=always \
  ${IMAGE_NAME}

echo "部署完成！"
echo "访问地址: http://localhost:18080"
echo "查看日志: docker logs -f ${PROJECT_NAME}"