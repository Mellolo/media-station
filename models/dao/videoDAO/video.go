package videoDAO

import (
	"github.com/Mellolo/media-station/models/dao/daoCommon"
	"github.com/beego/beego/v2/client/orm"
)

const (
	TableVideo = "video"
)

type VideoRecord struct {
	daoCommon.CommonColumn `orm:"inline"`

	Name            string  `orm:"column(name)"`
	Description     string  `orm:"column(description)"`
	Uploader        string  `orm:"column(uploader)"` // 上传者用户名
	VideoUrl        string  `orm:"column(video_url)"`
	CoverUrl        string  `orm:"column(cover_url)"`
	Duration        float64 `orm:"column(duration)"`
	PermissionLevel string  `orm:"column(permission_level)"` // 权限等级
}

// 设置表名
func (v *VideoRecord) TableName() string {
	return TableVideo
}

// 自动注册 model
func init() {
	orm.RegisterModel(new(VideoRecord))
}
