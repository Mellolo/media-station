package tagDO

import "github.com/beego/beego/v2/client/orm"

type TagDO struct {
	Id       int64
	CreateAt orm.DateTimeField
	Name     string
	Creator  string
	Details  TagDetails
}

type TagDetails struct {
	VideoIds   []int64 `json:"videoIds,omitempty"`
	GalleryIds []int64 `json:"galleryIds,omitempty"`
}
