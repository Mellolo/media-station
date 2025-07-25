package actorDO

// +k8s:deepcopy-gen=true
// ActorDO 表示演员的数据对象
type ActorDO struct {
	Id          int64
	CreateAt    string
	Name        string
	Description string
	Creator     string
	CoverUrl    string
	Art         ActorArtDO
}

// +k8s:deepcopy-gen=true
type ActorArtDO struct {
	VideoIds   []int64 `json:"videoIds,omitempty"`
	GalleryIds []int64 `json:"galleryIds,omitempty"`
}
