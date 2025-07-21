package fileDTO

import (
	"io"
)

type FileDTO struct {
	File io.ReadCloser
	Size int64
}
