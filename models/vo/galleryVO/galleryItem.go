package galleryVO

// +k8s:deepcopy-gen=true
type GalleryItemVO struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	CoverUrl        string `json:"coverUrl"`
	PermissionLevel string `json:"permissionLevel"`
}