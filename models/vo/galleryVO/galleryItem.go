package galleryVO

// +k8s:deepcopy-gen=true
type GalleryItemVO struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	PageCount       int    `json:"pageCount"`
	PermissionLevel string `json:"permissionLevel"`
}
