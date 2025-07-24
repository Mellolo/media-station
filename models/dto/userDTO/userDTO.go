package userDTO

type UserRegisterDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserClaimDTO struct {
	Username string `json:"username"`
}

type UserProfileDTO struct {
	Username string         `json:"username"`
	Details  UserDetailsDTO `json:"details"`
}

type UserDetailsDTO struct {
	VideoIds   []int64 `json:"videoIds"`
	GalleryIds []int64 `json:"galleryIds"`
}
