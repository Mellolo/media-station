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
	Id              int64    `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Actors          []int64  `json:"actors"`
	Tags            []string `json:"tags"`
	PermissionLevel string   `json:"permissionLevel"`
}

// +k8s:deepcopy-gen=true
type VideoSearchDTO struct {
	Keyword string   `json:"keyword"`
	Actors  []int64  `json:"actors"`
	Tags    []string `json:"tags"`
}

// +k8s:deepcopy-gen=true
type VideoItemDTO struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	CoverUrl        string `json:"coverUrl"`
	PermissionLevel string `json:"permissionLevel"`
}

// +k8s:deepcopy-gen=true
type VideoPageDTO struct {
	Id              int64    `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Actors          []int64  `json:"actors"`
	Tags            []string `json:"tags"`
	Uploader        string   `json:"uploader"`
	CoverUrl        string   `json:"coverUrl"`
	VideoUrl        string   `json:"videoUrl"`
	PermissionLevel string   `json:"permissionLevel"`
}

type VideoFileDTO struct {
	Reader io.ReadCloser
	Header http.Header
}

type VideoCoverDTO struct {
	Reader io.ReadCloser
	Header http.Header
}
