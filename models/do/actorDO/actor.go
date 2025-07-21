package actorDO

import "github.com/beego/beego/v2/client/orm"

// ActorDO 表示演员的数据对象
type ActorDO struct {
	Id          int64
	CreateAt    orm.DateTimeField
	Name        string
	Description string
	Creator     string
	CoverUrl    string
	Details     ActorDetailsDO
}

type ActorDetailsDO struct {
	VideoIds   []int64 `json:"videoIds,omitempty"`
	GalleryIds []int64 `json:"galleryIds,omitempty"`
}
