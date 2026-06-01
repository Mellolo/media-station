# Docker 部署方案

## 两步部署流程

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

---

## 前置配置（重要）

### 1. 本地Docker客户端配置
Docker Desktop → Settings → Docker Engine，添加:
```json
{
  "insecure-registries": [
    "192.168.5.178:5000"
  ]
}
```
点击 "Apply & Restart"

### 2. NAS Docker配置（必须）
Registry使用HTTP协议，需要在NAS上配置Docker允许HTTP连接：

```bash
ssh mellolo@192.168.5.178

# 配置Docker daemon
sudo vi /etc/docker/daemon.json
```

添加内容：
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
```

**配置完成后**，才能正常拉取镜像：http://192.168.5.178:5000/v2/