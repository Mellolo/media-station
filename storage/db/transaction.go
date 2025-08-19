package db

import (
	"context"
	"github.com/Mellolo/common/errors"
	"github.com/beego/beego/v2/client/orm"
)

func DoTransaction(f func(tx orm.TxOrmer)) {
	err := orm.NewOrm().DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		f(txOrm)
		return nil
	})
	if err != nil {
		panic(errors.WrapError(err, "transaction failed"))
	}
}
