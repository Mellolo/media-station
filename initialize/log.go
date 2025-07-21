package initialize

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/mellolo/common/errors"
)

func InitLog() {
	// 日志配置
	logConf := map[string]interface{}{
		"filename": web.AppConfig.DefaultString("log::logfile", "./logs/app.log"),
		"maxdays":  web.AppConfig.DefaultInt("log::max_days", 7),
		"maxSize":  web.AppConfig.DefaultInt("log::maxsize", 16),
		"loglevel": web.AppConfig.DefaultInt("loglevel", 7),
		"perm":     web.AppConfig.DefaultString("log::perm", "0666"),
		"separate": web.AppConfig.DefaultStrings("log::separate", []string{"error"}),
	}
	logConfBytes, marshalErr := json.Marshal(logConf)
	if marshalErr != nil {
		panic(errors.WrapError(marshalErr, "log configuration JSON serialization failed"))
	}
	if err := logs.SetLogger(logs.AdapterFile, string(logConfBytes)); err != nil {
		panic(errors.WrapError(err, "log configuration failed"))
	}

	// 日志显示函数调用行
	logFuncCall := web.AppConfig.DefaultBool("logfunc", false)
	logs.SetLogFuncCall(logFuncCall)
}
