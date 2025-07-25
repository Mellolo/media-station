package actorVO

// +k8s:deepcopy-gen=true
type ActorPageVO struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`

	VideoIds   []int64 `json:"videoIds"`
	GalleryIds []int64 `json:"galleryIds"`
}
