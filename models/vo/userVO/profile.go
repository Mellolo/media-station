package userVO

// +k8s:deepcopy-gen=true
type UserStatusProfileVO struct {
	Username string `json:"username"`
}

// +k8s:deepcopy-gen=true
type UserProfileVO struct {
	Username string `json:"username"`
}
