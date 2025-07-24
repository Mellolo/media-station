package videoVO

// +k8s:deepcopy-gen=true
type VideoPageVO struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	CoverUrl        string `json:"coverUrl"`
	PermissionLevel string `json:"permissionLevel"`
}