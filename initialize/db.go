package initialize

import (
	"github.com/Mellolo/common/errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func InitDB() {
	dataSource, err := web.AppConfig.String("mysql::data_source")
	if err != nil || dataSource == "" {
		if dataSource == "" {
			panic(errors.NewError("invalid data source"))
		}
		panic(errors.WrapError(err, "invalid data source"))
	}

	err = orm.RegisterDataBase("default", "mysql", dataSource)
	if err != nil {
		panic(errors.WrapError(err, "register database failed"))
	}
}
