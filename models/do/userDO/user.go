package userDO

type UserDO struct {
	Username    string
	Password    string
	PhoneNumber string
	WechatId    string
	Details     UserDetails
}

type UserDetails struct {
	VideoIds   []int `json:"videoIds,omitempty"`
	GalleryIds []int `json:"galleryIds,omitempty"`
}
