package actorDO

// +k8s:deepcopy-gen=true
type ActorDO struct {
	Id          int64
	CreateAt    string
	Name        string
	Description string
	Creator     string
	CoverUrl    string
}
