package userVO

// +k8s:deepcopy-gen=true
type UserProfileVO struct {
	Username string `json:"username"`
}