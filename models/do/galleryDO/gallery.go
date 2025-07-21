package galleryDO

import (
	"github.com/beego/beego/v2/client/orm"
	"io"
	"net/http"
)

// GalleryDO 表示画廊的数据对象
type GalleryDO struct {
	Id              int64
	CreateAt        orm.DateTimeField
	Name            string
	Description     string
	PageCount       int
	Actors          []int64
	Tags            []string
	Uploader        string
	CoverUrl        string
	GalleryUrl      string
	PermissionLevel string
}

type PictureDO struct {
	Reader io.ReadCloser
	Header http.Header
}
