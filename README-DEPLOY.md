# Docker 部署方案

## 两步部署流程

### 步骤1: 本地构建并推送镜像
```bash
make build-push
# 或 ./build-and-push.sh
```
将镜像推送到NAS Registry (192.168.5.178:5000)

### 步骤2: NAS上部署

#### 方式A: 上传脚本执行
```bash
make upload-deploy-script
ssh mellolo@192.168.5.178
chmod +x ~/deploy-on-nas.sh && ./deploy-on-nas.sh
```

#### 方式B: 手动执行
```bash
ssh mellolo@192.168.5.178
docker pull 192.168.5.178:5000/media-station:latest
docker stop media-station || true && docker rm media-station || true
docker run -d -p 18080:8080 --name media-station --restart=always 192.168.5.178:5000/media-station:latest
```

---

## 完成后

访问: http://192.168.5.178:18080

查看日志: `docker logs -f media-station` (在NAS上执行)

---

## 前置配置

Docker Desktop → Settings → Docker Engine，添加:
```json
{
  "insecure-registries": [
    "192.168.5.178:5000"
  ]
}
```
