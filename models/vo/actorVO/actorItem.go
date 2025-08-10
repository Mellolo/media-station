package actorVO

// +k8s:deepcopy-gen=true
type ActorItemVO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
