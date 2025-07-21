package videoDO

import (
	"github.com/beego/beego/v2/client/orm"
	"io"
	"net/http"
)

type VideoDO struct {
	Id              int64
	CreateAt        orm.DateTimeField
	Name            string
	Description     string
	Actors          []int64
	Tags            []string // 标签
	Uploader        string   // 上传者用户名
	VideoUrl        string
	CoverUrl        string
	PermissionLevel string // 权限等级
}

type VideoFileDO struct {
	Reader io.ReadCloser
	Header http.Header
}
