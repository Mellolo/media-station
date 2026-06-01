# 一键部署说明

## 部署流程

### 1. 推送代码到Git
```bash
git add .
git commit -m "update"
git push origin main
```

### 2. NAS上克隆项目
```bash
ssh mellolo@192.168.5.178
cd ~
git clone <your-git-repo-url> media-station
cd media-station
```

### 3. 一键部署
```bash
make deploy
```
或者
```bash
./deploy.sh
```

## 部署完成后

- 访问地址: `http://192.168.5.178:8080`
- 查看日志: `docker logs -f media-station`
- 重启服务: `docker restart media-station`
- 停止服务: `docker stop media-station`

## 注意事项

- 确保NAS上的MySQL、Redis、Nacos、MinIO服务正常运行
- 配置文件使用 `conf/app.prod.conf`（Docker网关地址172.17.0.1）
- 容器会自动重启（`--restart=always`）