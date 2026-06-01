#!/bin/bash

# 在NAS上配置Docker允许HTTP连接Registry
# 使用方法:
#   1. 上传此脚本到NAS: scp config-nas-docker.sh mellolo@192.168.5.178:~/config-nas-docker.sh
#   2. SSH连接到NAS: ssh mellolo@192.168.5.178
#   3. 执行此脚本: sudo ./config-nas-docker.sh

REGISTRY="192.168.5.178:5000"

echo "======================================"
echo "  配置NAS Docker客户端"
echo "======================================"
echo ""
echo "Registry: http://${REGISTRY}"
echo ""

# 创建或更新daemon.json
echo "配置Docker daemon.json..."
if [ -f /etc/docker/daemon.json ]; then
    # 文件存在，需要合并配置
    echo "检测到已有配置文件，将合并配置..."
    
    # 读取现有配置
    existing_config=$(cat /etc/docker/daemon.json)
    
    # 检查是否已有insecure-registries配置
    if echo "$existing_config" | grep -q "insecure-registries"; then
        # 已有配置，检查是否包含我们的Registry
        if echo "$existing_config" | grep -q "$REGISTRY"; then
            echo "✓ 配置已存在，无需修改"
        else
            # 需要添加Registry到现有列表
            echo "添加Registry到现有insecure-registries列表..."
            # 使用Python或jq来处理JSON（如果可用）
            if command -v python3 &> /dev/null; then
                python3 -c "
import json
import sys
with open('/etc/docker/daemon.json', 'r') as f:
    config = json.load(f)
if 'insecure-registries' in config:
    if '$REGISTRY' not in config['insecure-registries']:
        config['insecure-registries'].append('$REGISTRY')
else:
    config['insecure-registries'] = ['$REGISTRY']
with open('/etc/docker/daemon.json', 'w') as f:
    json.dump(config, f, indent=2)
"
            else
                echo "警告: 未找到python3，请手动编辑/etc/docker/daemon.json"
                echo "添加以下配置:"
                echo '  "insecure-registries": ["192.168.5.178:5000"]'
                exit 1
            fi
        fi
    else
        # 没有insecure-registries，需要添加
        if command -v python3 &> /dev/null; then
            python3 -c "
import json
with open('/etc/docker/daemon.json', 'r') as f:
    config = json.load(f)
config['insecure-registries'] = ['$REGISTRY']
with open('/etc/docker/daemon.json', 'w') as f:
    json.dump(config, f, indent=2)
"
        else
            echo "警告: 未找到python3，请手动编辑/etc/docker/daemon.json"
            echo "添加以下配置:"
            echo '  "insecure-registries": ["192.168.5.178:5000"]'
            exit 1
        fi
    fi
else
    # 文件不存在，创建新文件
    echo "创建新的配置文件..."
    echo '{"insecure-registries": ["192.168.5.178:5000"]}' > /etc/docker/daemon.json
fi

echo "✓ 配置完成"
echo ""

# 显示配置内容
echo "当前配置:"
cat /etc/docker/daemon.json
echo ""

# 重启Docker服务
echo "重启Docker服务..."
if command -v systemctl &> /dev/null; then
    systemctl restart docker
elif command -v service &> /dev/null; then
    service docker restart
else
    echo "警告: 无法找到systemctl或service命令"
    echo "请手动重启Docker"
fi

echo "✓ Docker服务已重启"
echo ""

# 验证配置
echo "验证配置..."
docker info | grep -A 5 "Insecure Registries"

echo ""
echo "======================================"
echo "  配置完成"
echo "======================================"
echo ""
echo "现在可以正常使用HTTP Registry了"
echo ""
echo "测试命令:"
echo "  docker pull 192.168.5.178:5000/media-station:latest"