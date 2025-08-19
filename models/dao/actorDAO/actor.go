package actorDAO

import (
	"github.com/Mellolo/media-station/models/dao/daoCommon"
	"github.com/beego/beego/v2/client/orm"
)

const (
	TableActor = "actor"
)

type ActorRecord struct {
	daoCommon.CommonColumn `orm:"inline"`

	Name        string `orm:"column(name)"`
	Description string `orm:"column(description)"`
	Creator     string `orm:"column(creator)"`
	CoverUrl    string `orm:"column(cover_url)"`
	Details     string `orm:"column(details)"`
}

func (a *ActorRecord) TableName() string {
	return TableActor
}

func init() {
	orm.RegisterModel(new(ActorRecord))
}
