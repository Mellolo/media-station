package videoDAO

import (
	"github.com/beego/beego/v2/client/orm"
	"media-station/models/dao/daoCommon"
)

const (
	TableVideo = "video"
)

type VideoRecord struct {
	daoCommon.CommonColumn `orm:"inline"`

	Name            string `orm:"column(name)"`
	Description     string `orm:"column(description)"`
	Actors          string `orm:"column(actors)"`   // 演员列表
	Tags            string `orm:"column(tags)"`     // 标签
	Uploader        string `orm:"column(uploader)"` // 上传者用户名
	VideoUrl        string `orm:"column(video_url)"`
	CoverUrl        string `orm:"column(cover_url)"`
	PermissionLevel string `orm:"column(permission_level)"` // 权限等级
}

// 设置表名
func (v *VideoRecord) TableName() string {
	return TableVideo
}

// 自动注册 model
func init() {
	orm.RegisterModel(new(VideoRecord))
}
