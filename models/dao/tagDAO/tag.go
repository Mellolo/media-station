package tagDAO

import (
	"github.com/Mellolo/media-station/models/dao/daoCommon"
	"github.com/beego/beego/v2/client/orm"
)

const (
	TableTag = "tag"
)

type TagRecord struct {
	daoCommon.CommonColumn `orm:"inline"`

	ArtType string `orm:"column(art_type)"`
	ArtId   int64  `orm:"column(art_id)"`
	Tag     string `orm:"column(tag)"`
}

func (t *TagRecord) TableName() string {
	return TableTag
}

func init() {
	orm.RegisterModel(new(TagRecord))
}
