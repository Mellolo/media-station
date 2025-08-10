package galleryVO

// +k8s:deepcopy-gen=true
type GalleryItemVO struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	PermissionLevel string `json:"permissionLevel"`
}
