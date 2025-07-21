package oss

import (
	"fmt"
	"github.com/mellolo/common/errors"
	"github.com/mellolo/common/oss"
	"io"
	"media-station/models/do/galleryDO"
)

type PictureStorage interface {
	Upload(bucket, path string, file io.ReadCloser, size int64)
	Download(bucket, path string) galleryDO.PictureDO
	Remove(bucket, path string)
}

type PictureStorageImpl struct {
	client oss.ObjectStorageClient
}

func NewPictureStorage() PictureStorage {
	client, err := oss.GetOss()
	if err != nil {
		panic(errors.WrapError(err, "get oss client failed"))
	}
	return &PictureStorageImpl{
		client: client,
	}
}

func (impl *PictureStorageImpl) Upload(bucket, path string, file io.ReadCloser, size int64) {
	err := impl.client.RemoveIncompleteUpload(bucket, path)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("remove incomplete upload pic [%s] failed", path)))
	}

	err = impl.client.PutObject(bucket, path, file, size)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("upload pic [%s] failed", path)))
	}
}

func (impl *PictureStorageImpl) Download(bucket, path string) galleryDO.PictureDO {
	reader, header, err := impl.client.GetObjectReader(bucket, path)
	if err != nil {
		panic(errors.WrapError(err, "get pic failed"))
	}

	return galleryDO.PictureDO{
		Reader: reader,
		Header: header,
	}
}

func (impl *PictureStorageImpl) Remove(bucket, path string) {
	err := impl.client.DeleteObject(bucket, path)
	if err != nil {
		panic(errors.WrapError(err, "remove pic failed"))
	}
}
