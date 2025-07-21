package oss

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/google/uuid"
	"github.com/mellolo/common/errors"
	"github.com/mellolo/common/oss"
	"io"
	"media-station/models/do/videoDO"
	"media-station/util"
)

const (
	filePartSize int64 = 16 * 1024 * 1024
)

type VideoStorage interface {
	Upload(bucket, path string, file io.ReadCloser, size int64, ch chan string)
	Download(bucket, path string, rangeHeader ...string) videoDO.VideoFileDO
	Remove(bucket, path string)
}

type VideoStorageImpl struct {
	client oss.ObjectStorageClient
}

func NewVideoStorage() VideoStorage {
	client, err := oss.GetOss()
	if err != nil {
		panic(errors.WrapError(err, "get oss client failed"))
	}
	return &VideoStorageImpl{
		client: client,
	}
}

func (impl *VideoStorageImpl) Upload(bucket, path string, file io.ReadCloser, size int64, ch chan string) {
	defer func() {
		if ch != nil {
			close(ch)
		}
	}()

	panicContext := errors.CatchPanic(func() {
		err := impl.client.RemoveIncompleteUpload(bucket, path)
		if err != nil {
			panic(errors.WrapError(err, fmt.Sprintf("remove incomplete upload video [%s] failed", path)))
		}

		if size > filePartSize {
			partNum := int(size / filePartSize)
			if size%filePartSize != 0 {
				partNum++
			}
			var parts []oss.Part
			// 拆分上传
			uploadID, _ := impl.client.NewMultipartUpload(bucket, path)
			for partNumber := 1; partNumber <= partNum; partNumber++ {
				eTag, putErr := impl.client.PutObjectPart(bucket, path, uploadID, partNumber, file, filePartSize)
				if putErr != nil {
					panic(errors.WrapError(putErr, fmt.Sprintf("upload video [%s] part [%d] failed", path, partNumber)))
				}
				parts = append(parts, oss.Part{
					PartNumber: partNumber,
					ETag:       eTag,
				})
				logs.Info(fmt.Sprintf("video uploaded part %d, ETag: %s", partNumber, eTag))

				if ch != nil {
					ch <- util.GetProcessBarJsonString(float64(partNumber) / float64(partNum))
				}
			}
			err = impl.client.CompleteMultipartUpload(bucket, path, uploadID, parts)
			if err != nil {
				panic(errors.WrapError(err, fmt.Sprintf("complete upload video [%s] failed", path)))
			}
		} else {
			err = impl.client.PutObject(bucket, path, file, size)
			if err != nil {
				panic(errors.WrapError(err, fmt.Sprintf("upload video [%s] failed", path)))
			}
		}

		if ch != nil {
			ch <- util.GetDoneProcessBarJsonString()
		}
	})

	if panicContext.Err != nil {
		ch <- util.GetFailedProcessBarJsonString()
		uniqueId, _ := uuid.NewV7()
		logs.Error(
			fmt.Sprintf("upload video failed\n%s",
				util.FormatErrorLog(uniqueId.String(), panicContext.Err.Error(), panicContext.RecoverStack),
			))
	}
}

func (impl *VideoStorageImpl) Download(bucket, path string, rangeHeader ...string) videoDO.VideoFileDO {
	reader, header, err := impl.client.GetObjectReader(bucket, path, rangeHeader...)
	if err != nil {
		panic(errors.WrapError(err, "get video failed"))
	}

	return videoDO.VideoFileDO{
		Reader: reader,
		Header: header,
	}
}

func (impl *VideoStorageImpl) Remove(bucket, path string) {
	err := impl.client.DeleteObject(bucket, path)
	if err != nil {
		panic(errors.WrapError(err, "remove pic failed"))
	}
}
