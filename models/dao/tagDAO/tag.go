package tagDAO

import (
	"github.com/beego/beego/v2/client/orm"
	"media-station/models/dao/daoCommon"
)

const (
	TableTag = "tag"
)

type TagRecord struct {
	daoCommon.CommonColumn `orm:"inline"`

	Name    string `orm:"column(name)"`
	Creator string `orm:"column(creator)"` // 创建者用户名或ID

	Details string `orm:"column(details)"`
}

func (t *TagRecord) TableName() string {
	return TableTag
}

func init() {
	orm.RegisterModel(new(TagRecord))
}
