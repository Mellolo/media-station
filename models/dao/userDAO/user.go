package userDAO

import (
	"github.com/Mellolo/media-station/models/dao/daoCommon"
	"github.com/beego/beego/v2/client/orm"
)

const (
	TableUser = "user"
)

type UserRecord struct {
	daoCommon.CommonColumn `orm:"inline"`

	Username string `orm:"column(username)"`
	Password string `orm:"column(password)"`
	Details  string `orm:"column(details)"`
}

func (g *UserRecord) TableName() string {
	return TableUser
}

func init() {
	orm.RegisterModel(new(UserRecord))
}
