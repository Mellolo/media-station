package galleryDO

import (
	"io"
	"net/http"
)

// +k8s:deepcopy-gen=true
// GalleryDO 表示画廊的数据对象
type GalleryDO struct {
	Id              int64
	CreateAt        string
	Name            string
	Description     string
	Uploader        string
	DirPath         string
	PicPaths        []string
	PermissionLevel string
}

type PictureDO struct {
	Reader io.ReadCloser
	Header http.Header
}
