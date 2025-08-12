package videoDO

import (
	"io"
	"net/http"
)

// +k8s:deepcopy-gen=true
type VideoDO struct {
	Id              int64
	CreateAt        string
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

type VideoCoverFileDO struct {
	Reader io.ReadCloser
	Header http.Header
}
