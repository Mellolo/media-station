package galleryDTO

import (
	"io"
	"net/http"
)

// +k8s:deepcopy-gen=true
type GalleryCreateDTO struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Uploader        string `json:"uploader"`
	PermissionLevel string `json:"permissionLevel"`
}

// +k8s:deepcopy-gen=true
type GalleryUpdateDTO struct {
	Id              int64                  `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Pages           []GalleryUpdatePageDTO `json:"pages"`
	PermissionLevel string                 `json:"permissionLevel"`
}

// +k8s:deepcopy-gen=true
type GalleryUpdatePageDTO struct {
	IsNewUploaded bool `json:"isNewUploaded"`
	Index         int  `json:"index"`
}

// +k8s:deepcopy-gen=true
type GallerySearchDTO struct {
	Ids     []int64 `json:"ids"`
	Keyword string  `json:"keyword"`
}

// +k8s:deepcopy-gen=true
type GalleryDTO struct {
	Id              int64    `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Uploader        string   `json:"uploader"`
	DirPath         string   `json:"dirPath"`
	PicPaths        []string `json:"picPaths"`
	PermissionLevel string   `json:"permissionLevel"`
}

type PictureFileDTO struct {
	Reader io.ReadCloser
	Header http.Header
}
