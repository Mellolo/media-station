# Docker 部署方案

## ⚠️ 前置配置（首次必须完成）

Registry使用HTTP协议（http://192.168.5.178:5000/v2/），Docker默认要求HTTPS，需要配置insecure-registries.

### 步骤0: 配置NAS Docker客户端（必须）

**使用自动化脚本**（推荐）：
```bash
# 上传配置脚本
make config-nas-docker

# SSH连接并执行（需要sudo密码）
ssh mellolo@192.168.5.178
sudo ./config-nas-docker.sh
```

**手动配置**：
```bash
ssh mellolo@192.168.5.178

# 编辑Docker配置
sudo vi /etc/docker/daemon.json
```

添加内容（如果文件已有内容，合并配置）：
```json
{
  "insecure-registries": [
    "192.168.5.178:5000"
  ]
}
```

重启Docker：
```bash
sudo systemctl restart docker
# 或 sudo service docker restart
```

验证配置：
```bash
sudo docker info | grep -A 5 "Insecure Registries"
# 应该看到: 192.168.5.178:5000
```

### 本地Docker客户端配置

Docker Desktop → Settings → Docker Engine，添加:
```json
{
  "insecure-registries": [
    "192.168.5.178:5000"
  ]
}
```
点击 "Apply & Restart"

---

## 部署流程（前置配置完成后）

### 步骤1: 本地构建并推送镜像
```bash
make build-push
# 或 ./build-and-push.sh
```
将镜像推送到NAS Registry (http://192.168.5.178:5000)

### 步骤2: NAS上部署

#### 方式A: 上传脚本执行
```bash
make upload-deploy-script
ssh mellolo@192.168.5.178
chmod +x ~/deploy-on-nas.sh && sudo ./deploy-on-nas.sh
```

#### 方式B: 手动执行
```bash
ssh mellolo@192.168.5.178
sudo docker pull 192.168.5.178:5000/media-station:latest
sudo docker stop media-station || true && sudo docker rm media-station || true
sudo docker run -d -p 18080:8080 --name media-station --restart=always 192.168.5.178:5000/media-station:latest
```

---

## 完成后

访问: http://192.168.5.178:18080

查看日志: `sudo docker logs -f media-station` (在NAS上执行)