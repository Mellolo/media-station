#!/bin/bash

# 媒体站 Docker 镜像构建脚本
# 使用项目根目录的 Dockerfile

set -e

# 默认参数
PROJECT_NAME="media-station"
VERSION="latest"
IMAGE_NAME="media-station"
VERBOSE=false

# 显示帮助信息
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -n, --name NAME        镜像名称 (默认: media-station)"
    echo "  -v, --version VERSION  镜像版本 (默认: latest)"
    echo "  --verbose              显示详细信息"
    echo "  -h, --help             显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0                                    # 使用默认参数构建"
    echo "  $0 -n my-media-station -v 1.0        # 指定名称和版本"
    echo "  $0 --name media-app --version v1.2   # 使用完整参数"
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--name)
            IMAGE_NAME="$2"
            shift 2
            ;;
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "未知参数: $1"
            show_help
            exit 1
            ;;
    esac
done

# 组合完整镜像名
FULL_IMAGE_NAME="${IMAGE_NAME}:${VERSION}"

echo "开始构建媒体站 Docker 镜像..."
echo "镜像名称: ${FULL_IMAGE_NAME}"
if [ "$VERBOSE" = true ]; then
    echo "项目目录: $(pwd)"
fi

# 构建 Docker 镜像
echo "构建 Docker 镜像: ${FULL_IMAGE_NAME}"
BUILD_CMD="docker build -t ${FULL_IMAGE_NAME} ."
if [ "$VERBOSE" = true ]; then
    echo "执行命令: $BUILD_CMD"
    echo "构建上下文: $(pwd)"
    echo "Dockerfile 路径: ./Dockerfile"
fi

# 执行构建命令
eval $BUILD_CMD

echo ""
echo "Docker 镜像构建完成: ${FULL_IMAGE_NAME}"
echo ""
echo "使用以下命令运行容器:"
echo "  docker run -d -p 8080:8080 --name ${IMAGE_NAME} ${FULL_IMAGE_NAME}"
echo ""
echo "使用以下命令查看运行日志:"
echo "  docker logs -f ${IMAGE_NAME}"
echo ""