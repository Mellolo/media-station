package galleryVO

// +k8s:deepcopy-gen=true
type GalleryPageVO struct {
	Id              int64            `json:"id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	PageCount       int              `json:"pageCount"`
	Actors          []GalleryActorVO `json:"actors"`
	Tags            []string         `json:"tags"`
	PermissionLevel string           `json:"permissionLevel"`
}

// +k8s:deepcopy-gen=true
type GalleryActorVO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
