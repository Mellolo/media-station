package videoVO

// +k8s:deepcopy-gen=true
type VideoItemVO struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	PermissionLevel string `json:"permissionLevel"`
}
