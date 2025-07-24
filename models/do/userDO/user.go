package userDO

// +k8s:deepcopy-gen=true
type UserDO struct {
	Username    string
	Password    string
	PhoneNumber string
	WechatId    string
	Details     UserDetails
}

// +k8s:deepcopy-gen=true
type UserDetails struct {
	VideoIds   []int64 `json:"videoIds,omitempty"`
	GalleryIds []int64 `json:"galleryIds,omitempty"`
}