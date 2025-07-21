package galleryVO

import (
	"io"
	"net/http"
)

type GalleryFileVO struct {
	Reader io.ReadCloser
	Header http.Header
}
