package videoVO

// +k8s:deepcopy-gen=true
type VideoPageVO struct {
	Id              int64          `json:"id"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	Actors          []VideoActorVO `json:"actors"`
	Tags            []string       `json:"tags"`
	PermissionLevel string         `json:"permissionLevel"`
}

// +k8s:deepcopy-gen=true
type VideoActorVO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
