package videoVO

import (
	"io"
	"net/http"
)

type VideoFileVO struct {
	Reader io.ReadCloser
	Header http.Header
}

type VideoCoverFileVO struct {
	Reader io.ReadCloser
	Header http.Header
}
