package performDO

// +k8s:deepcopy-gen=true
type PerformDO struct {
	ArtType string
	ArtId   int64
	ActorId int64
}
