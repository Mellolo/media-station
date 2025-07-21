package actorDTO

// ActorCreateDTO 表示创建演员的数据传输对象
type ActorCreateDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
}

type ActorUpdateDTO struct {
	Id          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Creator     string          `json:"creator"`
	Details     ActorDetailsDTO `json:"details"`
}

type ActorRemoveArtDTO struct {
	Id      int64           `json:"id"`
	Details ActorDetailsDTO `json:"details"`
}

type ActorDetailsDTO struct {
	VideoIds   []int64 `json:"videoIds"`
	GalleryIds []int64 `json:"galleryIds"`
}
