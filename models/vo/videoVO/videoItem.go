package videoVO

// +k8s:deepcopy-gen=true
type VideoItemVO struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	CoverUrl        string `json:"coverUrl"`
	PermissionLevel string `json:"permissionLevel"`
}