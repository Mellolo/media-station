package actorDTO

// +k8s:deepcopy-gen=true
// ActorCreateDTO 表示创建演员的数据传输对象
type ActorCreateDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
}

// +k8s:deepcopy-gen=true
type ActorUpdateDTO struct {
	Id          int64       `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Art         ActorArtDTO `json:"art"`
}

// +k8s:deepcopy-gen=true
type ActorRemoveArtDTO struct {
	Id  int64       `json:"id"`
	Art ActorArtDTO `json:"art"`
}

// +k8s:deepcopy-gen=true
type ActorPageDTO struct {
	Id          int64       `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Creator     string      `json:"creator"`
	CoverUrl    string      `json:"coverUrl"`
	Art         ActorArtDTO `json:"art"`
}

// +k8s:deepcopy-gen=true
type ActorArtDTO struct {
	VideoIds   []int64 `json:"videoIds"`
	GalleryIds []int64 `json:"galleryIds"`
}

// +k8s:deepcopy-gen=true
type ActorItemDTO struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CoverUrl    string `json:"coverUrl"`
}
