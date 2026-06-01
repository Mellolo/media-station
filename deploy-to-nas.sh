#!/bin/bash

# 一键部署到NAS脚本
# 使用方法: ./deploy-to-nas.sh

set -e

NAS_HOST="192.168.5.178"
NAS_USER="mellolo"
PROJECT_NAME="media-station"
IMAGE_NAME="media-station:latest"
IMAGE_FILE="media-station-image.tar"

echo "=== 开始一键部署到NAS ==="

# 1. 本地构建Docker镜像
echo "1. 本地构建Docker镜像..."
docker build -t ${IMAGE_NAME} .

# 2. 导出镜像为tar文件
echo "2. 导出镜像为tar文件..."
docker save -o ${IMAGE_FILE} ${IMAGE_NAME}

# 3. 上传镜像到NAS
echo "3. 上传镜像到NAS..."
scp ${IMAGE_FILE} ${NAS_USER}@${NAS_HOST}:~/

# 4. 在NAS上导入并运行镜像
echo "4. 在NAS上导入并运行镜像..."
ssh ${NAS_USER}@${NAS_HOST} << 'ENDSSH'
# 导入镜像
docker load -i ~/media-station-image.tar

# 停止并删除旧容器
docker stop media-station || true
docker rm media-station || true

# 启动新容器
docker run -d \
  -p 18080:8080 \
  --name media-station \
  --restart=always \
  media-station:latest

# 清理临时文件
rm ~/media-station-image.tar

echo "NAS部署完成！"
ENDSSH

# 5. 清理本地临时文件
echo "5. 清理本地临时文件..."
rm ${IMAGE_FILE}

echo "=== 部署成功 ==="
echo "访问地址: http://${NAS_HOST}:18080"
echo "查看日志: ssh ${NAS_USER}@${NAS_HOST} 'docker logs -f media-station'"