package tagDO

// +k8s:deepcopy-gen=true
type TagDO struct {
	Id       int64
	CreateAt string
	Name     string
	Creator  string
	Art      TagArt
}

// +k8s:deepcopy-gen=true
type TagArt struct {
	VideoIds   []int64 `json:"videoIds,omitempty"`
	GalleryIds []int64 `json:"galleryIds,omitempty"`
}
