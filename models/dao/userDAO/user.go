package userDAO

import (
	"github.com/beego/beego/v2/client/orm"
	"media-station/models/dao/daoCommon"
)

const (
	TableUser = "user"
)

type UserRecord struct {
	daoCommon.CommonColumn `orm:"inline"`

	Username    string `orm:"column(username)"`
	Password    string `orm:"column(password)"`
	PhoneNumber string `orm:"column(phone_number)"`
	WechatId    string `orm:"column(wechat_id)"`

	Details string `orm:"column(details)"`
}

func (g *UserRecord) TableName() string {
	return TableUser
}

func init() {
	orm.RegisterModel(new(UserRecord))
}
