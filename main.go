package main

import (
	"github.com/Mellolo/common/errors"
	"github.com/beego/beego/v2/server/web"
	"media-station/initialize"
	_ "media-station/routers"
)

const (
	MaxMemory int64 = 10 * 1024 * 1024 * 1024
)

func init() {
	// 加载配置文件
	if err := web.LoadAppConfig("ini", "conf/app.conf"); err != nil {
		panic(errors.WrapError(err, "load app configuration failed"))
	}

	web.BConfig.MaxMemory = MaxMemory
	web.BConfig.MaxUploadSize = MaxMemory

	// 初始化日志配置
	initialize.InitLog()

	// 初始化配置中心
	initialize.InitConfig()

	// 初始化数据库配置
	initialize.InitDB()

	// 初始化对象存储配置
	initialize.InitOss()

	// 初始化缓存
	initialize.InitCache()

	// 初始化消息队列配置
	initialize.InitMessageQueue()
}

func main() {
	web.Run()
}
