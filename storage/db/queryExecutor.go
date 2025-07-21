package db

import "github.com/beego/beego/v2/client/orm"

func getQueryExecutor(tx ...orm.TxOrmer) orm.QueryExecutor {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0]
	} else {
		return orm.NewOrm()
	}
}
