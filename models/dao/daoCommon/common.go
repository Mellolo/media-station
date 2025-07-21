package daoCommon

import "github.com/beego/beego/v2/client/orm"

type CommonColumn struct {
	Id        int64             `orm:"column(id)"`
	CreatedAt orm.DateTimeField `orm:"column(created_at);auto_now_add;type(datetime)"`
	UpdatedAt orm.DateTimeField `orm:"column(updated_at);auto_now;type(datetime)"`
}
