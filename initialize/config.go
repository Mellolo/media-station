package initialize

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/mellolo/common/config"
	"github.com/mellolo/common/errors"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func InitConfig() {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(
			web.AppConfig.DefaultString("nacos::nacos_addr", "localhost"),
			uint64(web.AppConfig.DefaultInt64("nacos::nacos_port", 8848))),
	}

	cc := constant.ClientConfig{
		NamespaceId:         web.AppConfig.DefaultString("nacos::nacos_namespace", ""),
		NotLoadCacheAtStart: true,
		LogDir:              web.AppConfig.DefaultString("nacos::nacos_log_dir", "./logs/nacos/"),
		CacheDir:            web.AppConfig.DefaultString("nacos::nacos_cache_dir", "./cache/nacos"),
		LogLevel:            web.AppConfig.DefaultString("nacos::nacos_log_level", ""),
		Username:            web.AppConfig.DefaultString("nacos::nacos_username", "nacos"),
		Password:            web.AppConfig.DefaultString("nacos::nacos_password", "nacos"),
	}

	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  &cc,
		ServerConfigs: sc,
	})
	if err != nil {
		panic(errors.WrapError(err, "new configuration client failed"))
	}

	// 添加nacos配置客户端的初始化
	err = config.InitConfigClient(config.NewNacosConfigClient(client))
	if err != nil {
		panic(errors.WrapError(err, "init configuration client failed"))
	}
}
