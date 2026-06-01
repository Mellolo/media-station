# 一键部署说明

## 方案选择

### 方案一：本地构建 + 一键推送到NAS（推荐）
**优势**: 本地构建速度快，NAS上无需编译

#### 使用方法
```bash
make deploy-nas
# 或者
./deploy-to-nas.sh
```

#### 部署流程
1. 本地构建Docker镜像
2. 导出镜像为tar文件
3. 通过scp上传到NAS
4. NAS上导入镜像并运行
5. 自动清理临时文件

### 方案二：NAS上直接构建部署
**适用场景**: 无法SSH连接NAS时

#### 部署流程

##### 1. 推送代码到Git
```bash
git add .
git commit -m "update"
git push origin main
```

##### 2. NAS上克隆项目
```bash
ssh mellolo@192.168.5.178
cd ~
git clone https://github.com/Mellolo/media-station.git
cd media-station
```

##### 3. 一键部署
```bash
make deploy
```
或者
```bash
./deploy.sh
```

## 部署完成后

- 访问地址: `http://192.168.5.178:18080`
- 查看日志: `docker logs -f media-station`
- 重启服务: `docker restart media-station`
- 停止服务: `docker stop media-station`

## 注意事项

- 确保NAS上的MySQL、Redis、Nacos、MinIO服务正常运行
- 配置文件使用 `conf/app.prod.conf`（Docker网关地址172.17.0.1）
- 容器会自动重启（`--restart=always`）