package videoDTO

import (
	"io"
	"net/http"
)

// +k8s:deepcopy-gen=true
type VideoCreateDTO struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Actors          []int64  `json:"actors"`
	Tags            []string `json:"tags"`
	Uploader        string   `json:"uploader"`
	PermissionLevel string   `json:"permissionLevel"`
}

// +k8s:deepcopy-gen=true
type VideoUpdateDTO struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	PermissionLevel string `json:"permissionLevel"`
}

// +k8s:deepcopy-gen=true
type VideoSearchDTO struct {
	Ids     []int64 `json:"ids"`
	Keyword string  `json:"keyword"`
}

// +k8s:deepcopy-gen=true
type VideoDTO struct {
	Id              int64   `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Uploader        string  `json:"uploader"`
	CoverUrl        string  `json:"coverUrl"`
	VideoUrl        string  `json:"videoUrl"`
	Duration        float64 `json:"duration"`
	PermissionLevel string  `json:"permissionLevel"`
}

type VideoFileDTO struct {
	Reader io.ReadCloser
	Header http.Header
}

type VideoCoverDTO struct {
	Reader io.ReadCloser
	Header http.Header
}
