package tagDTO

// +k8s:deepcopy-gen=true
type TagCreateOrUpdateDTO struct {
	Name    string        `json:"name"`
	Creator string        `json:"creator"`
	Details TagDetailsDTO `json:"details"`
}

// +k8s:deepcopy-gen=true
type TagRemoveArtDTO struct {
	Name    string        `json:"name"`
	Details TagDetailsDTO `json:"details"`
}

// +k8s:deepcopy-gen=true
type TagDetailsDTO struct {
	VideoIds   []int64 `json:"videoIds"`
	GalleryIds []int64 `json:"galleryIds"`
}