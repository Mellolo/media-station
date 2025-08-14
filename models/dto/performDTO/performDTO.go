package performDTO

type ArtPerformDTO struct {
	ArtType  string  `json:"artType"`
	ArtId    int64   `json:"artId"`
	ActorIds []int64 `json:"actorIds"`
}
