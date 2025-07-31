package galleryDTO

import (
	"io"
	"net/http"
)

// +k8s:deepcopy-gen=true
type GalleryCreateDTO struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Actors          []int64  `json:"actors"`
	Tags            []string `json:"tags"`
	Uploader        string   `json:"uploader"`
	PermissionLevel string   `json:"permissionLevel"`
}

// +k8s:deepcopy-gen=true
type GalleryUpdateDTO struct {
	Id              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Actors          []int64  `json:"actors"`
	Tags            []string `json:"tags"`
	Uploader        string   `json:"uploader"`
	PermissionLevel string   `json:"permissionLevel"`
}

// +k8s:deepcopy-gen=true
type GallerySearchDTO struct {
	Keyword string   `json:"keyword"`
	Actors  []int64  `json:"actors"`
	Tags    []string `json:"tags"`
}

// +k8s:deepcopy-gen=true
type GalleryItemDTO struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	CoverUrl        string `json:"coverUrl"`
	PermissionLevel string `json:"permissionLevel"`
}

// +k8s:deepcopy-gen=true
type GalleryPageDTO struct {
	Id              int64    `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Actors          []int64  `json:"actors"`
	Tags            []string `json:"tags"`
	Uploader        string   `json:"uploader"`
	CoverUrl        string   `json:"coverUrl"`
	GalleryUrl      string   `json:"galleryUrl"`
	PermissionLevel string   `json:"permissionLevel"`
}

type PictureFileDTO struct {
	Reader io.ReadCloser
	Header http.Header
}
