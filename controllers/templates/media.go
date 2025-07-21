package templates

import (
	"io"
	"net/http"
)

type PictureTemplate struct {
	Reader io.ReadCloser
	Header http.Header
}

type VideoTemplate struct {
	Reader io.ReadCloser
	Header http.Header
}
