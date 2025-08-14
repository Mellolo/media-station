package tagDTO

// +k8s:deepcopy-gen=true
type ArtTagDTO struct {
	ArtType string   `json:"artType"`
	ArtId   int64    `json:"artId"`
	Tags    []string `json:"tags"`
}
