package actorDTO

import (
	"io"
	"net/http"
)

// +k8s:deepcopy-gen=true
// ActorCreateDTO 表示创建演员的数据传输对象
type ActorCreateDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
}

// +k8s:deepcopy-gen=true
type ActorSearchDTO struct {
	Keyword string `json:"keyword"`
}

// +k8s:deepcopy-gen=true
type ActorUpdateDTO struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// +k8s:deepcopy-gen=true
type ActorDTO struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	CoverUrl    string `json:"coverUrl"`
}

type ActorCoverFileDTO struct {
	Reader io.ReadCloser
	Header http.Header
}
