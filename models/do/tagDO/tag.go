package tagDO

// +k8s:deepcopy-gen=true
type TagDO struct {
	Id       int64
	CreateAt string
	Name     string
	Creator  string
	Details  TagDetails
}

// +k8s:deepcopy-gen=true
type TagDetails struct {
	VideoIds   []int64 `json:"videoIds,omitempty"`
	GalleryIds []int64 `json:"galleryIds,omitempty"`
}
