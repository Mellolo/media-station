package userDTO

// +k8s:deepcopy-gen=true
type UserRegisterDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// +k8s:deepcopy-gen=true
type UserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// +k8s:deepcopy-gen=true
type UserClaimDTO struct {
	Username string `json:"username"`
}

// +k8s:deepcopy-gen=true
type UserProfileDTO struct {
	Username string         `json:"username"`
	Details  UserDetailsDTO `json:"details"`
}

// +k8s:deepcopy-gen=true
type UserDetailsDTO struct {
	VideoIds   []int64 `json:"videoIds"`
	GalleryIds []int64 `json:"galleryIds"`
}