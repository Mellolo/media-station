package actorVO

import (
	"io"
	"net/http"
)

type ActorCoverFileVO struct {
	Reader io.ReadCloser
	Header http.Header
}
