package galleryDTO

import (
	"io"
	"net/http"
)

type GalleryCreateDTO struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Actors          []int64  `json:"actors"`
	Tags            []string `json:"tags"`
	Uploader        string   `json:"uploader"`
	PermissionLevel string   `json:"permissionLevel"`
}

type GalleryUpdateDTO struct {
	Id              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Actors          []int64  `json:"actors"`
	Tags            []string `json:"tags"`
	Uploader        string   `json:"uploader"`
	PermissionLevel string   `json:"permissionLevel"`
}

type GallerySearchDTO struct {
	Keyword string   `json:"keyword"`
	Actors  []int64  `json:"actors"`
	Tags    []string `json:"tags"`
}

type GalleryItemDTO struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	CoverUrl        string `json:"coverUrl"`
	PermissionLevel string `json:"permissionLevel"`
}

type GalleryPageDTO struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Actors          []int64  `json:"actors"`
	Tags            []string `json:"tags"`
	Uploader        string   `json:"uploader"`
	CoverUrl        string   `json:"coverUrl"`
	GalleryUrl      string   `json:"galleryUrl"`
	PermissionLevel string   `json:"permissionLevel"`
}

type PictureDTO struct {
	Reader io.ReadCloser
	Header http.Header
}
