package performDAO

import (
	"github.com/beego/beego/v2/client/orm"
	"media-station/models/dao/daoCommon"
)

const (
	TablePerform = "perform"
)

type PerformRecord struct {
	daoCommon.CommonColumn `orm:"inline"`

	ArtType string `orm:"column(art_type)"`
	ArtId   int64  `orm:"column(art_id)"`
	ActorId int64  `orm:"column(actor_id)"`
}

func (a *PerformRecord) TableName() string {
	return TablePerform
}

func init() {
	orm.RegisterModel(new(PerformRecord))
}
