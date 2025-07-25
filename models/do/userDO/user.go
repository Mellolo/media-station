package userDO

// +k8s:deepcopy-gen=true
type UserDO struct {
	Username    string
	Password    string
	PhoneNumber string
	WechatId    string
	Art         UserArt
}

// +k8s:deepcopy-gen=true
type UserArt struct {
	VideoIds   []int64 `json:"videoIds,omitempty"`
	GalleryIds []int64 `json:"galleryIds,omitempty"`
}
