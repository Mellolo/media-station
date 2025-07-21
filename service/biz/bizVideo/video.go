package bizVideo

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/mellolo/common/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"media-station/enum"
	"media-station/generator"
	"media-station/models/do/videoDO"
	"media-station/models/dto/fileDTO"
	"media-station/models/dto/videoDTO"
	"media-station/storage/db"
	"media-station/storage/oss"
)

const (
	videoIdGenerateKey      = "video"
	videoCoverIdGenerateKey = "videoCover"
	bucketVideo             = "video"
)

type VideoBizService interface {
	SearchVideo(createDTO videoDTO.VideoSearchDTO, tx ...orm.TxOrmer) []videoDTO.VideoItemDTO
	CreateVideo(createDTO videoDTO.VideoCreateDTO, videoDTO, coverDTO fileDTO.FileDTO, ch chan string, tx ...orm.TxOrmer) int64
	UpdateVideo(id int64, updateDTO videoDTO.VideoUpdateDTO, coverDTO fileDTO.FileDTO, tx ...orm.TxOrmer) string
	DeleteVideo(id int64, tx ...orm.TxOrmer) (string, string)
	PlayVideo(id int64, rangeHeader ...string) videoDTO.VideoFileDTO

	RemoveVideoCover(path string)
	RemoveVideoFile(path string)
}

func NewVideoBizService() *VideoBizServiceImpl {
	return &VideoBizServiceImpl{
		videoMapper:    db.NewVideoMapper(),
		idGenerator:    generator.NewIdGenerator(),
		pictureStorage: oss.NewPictureStorage(),
		videoStorage:   oss.NewVideoStorage(),
	}
}

type VideoBizServiceImpl struct {
	videoMapper    db.VideoMapper
	idGenerator    generator.IdGenerator
	pictureStorage oss.PictureStorage
	videoStorage   oss.VideoStorage
}

func (impl *VideoBizServiceImpl) SearchVideo(createDTO videoDTO.VideoSearchDTO, tx ...orm.TxOrmer) []videoDTO.VideoItemDTO {
	// 读取数据库
	var videoDOList []*videoDO.VideoDO
	if createDTO.Keyword == "" {
		doList, err := impl.videoMapper.SelectAllLimit(200, tx...)
		if err != nil {
			panic(errors.WrapError(err, "select all video error"))
		}
		videoDOList = append(videoDOList, doList...)
	} else {
		doList, err := impl.videoMapper.SelectByKeyword(createDTO.Keyword)
		if err != nil {
			panic(errors.WrapError(err, fmt.Sprintf("select video by keyword [%s] error", createDTO.Keyword)))
		}
		videoDOList = append(videoDOList, doList...)
	}

	// 再筛选，并转化为DTO
	var items []videoDTO.VideoItemDTO
	for _, do := range videoDOList {
		if do.PermissionLevel == enum.PermissionForbidden || do.PermissionLevel == enum.PermissionPrivate {
			continue
		}
		if !sets.NewInt64(do.Actors...).HasAll(createDTO.Actors...) {
			continue
		}
		if !sets.NewString(do.Tags...).HasAll(createDTO.Tags...) {
			continue
		}

		items = append(items, videoDTO.VideoItemDTO{
			Id:              do.Id,
			Name:            do.Name,
			CoverUrl:        do.CoverUrl,
			PermissionLevel: do.PermissionLevel,
		})
	}

	return items
}

func (impl *VideoBizServiceImpl) CreateVideo(createDTO videoDTO.VideoCreateDTO, videoDTO, coverDTO fileDTO.FileDTO, ch chan string, tx ...orm.TxOrmer) int64 {
	if videoDTO.File == nil {
		panic(errors.NewError("video file is empty"))
	}

	video := &videoDO.VideoDO{
		Name:        createDTO.Name,
		Description: createDTO.Description,
		Actors:      createDTO.Actors,
		Tags:        createDTO.Tags,
		Uploader:    createDTO.Uploader,
	}
	if sets.NewString(enum.PermissionLevels...).Has(createDTO.PermissionLevel) {
		video.PermissionLevel = createDTO.PermissionLevel
	} else {
		video.PermissionLevel = enum.PermissionPublic
	}
	// 上传封面
	if coverDTO.File != nil {
		filename := impl.idGenerator.GenerateId(videoCoverIdGenerateKey)
		path := fmt.Sprintf("covers/%s.jpg", filename)
		impl.pictureStorage.Upload(bucketVideo, path, coverDTO.File, coverDTO.Size)
		video.CoverUrl = path
	}
	// 上传视频
	filename := impl.idGenerator.GenerateId(videoIdGenerateKey)
	path := fmt.Sprintf("videos/%s.mp4", filename)
	impl.videoStorage.Upload(bucketVideo, path, videoDTO.File, videoDTO.Size, ch)
	video.VideoUrl = path

	id, err := impl.videoMapper.Insert(video, tx...)
	if err != nil {
		panic(errors.WrapError(err, "create video failed"))
	}
	return id
}

func (impl *VideoBizServiceImpl) UpdateVideo(id int64, updateDTO videoDTO.VideoUpdateDTO, coverDTO fileDTO.FileDTO, tx ...orm.TxOrmer) string {
	video, err := impl.videoMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("video [%d] doesn't exist", id)))
	}
	lastCoverUrl := video.CoverUrl

	// 更新video
	if updateDTO.Name != "" {
		video.Name = updateDTO.Name
	}
	if updateDTO.Description != "" {
		video.Description = updateDTO.Description
	}
	if len(updateDTO.Actors) > 0 {
		video.Actors = updateDTO.Actors
	}
	if len(updateDTO.Tags) > 0 {
		video.Tags = updateDTO.Tags
	}
	if sets.NewString(enum.PermissionLevels...).Has(updateDTO.PermissionLevel) {
		video.PermissionLevel = updateDTO.PermissionLevel
	}
	// 更新封面
	if coverDTO.File != nil {
		filename := impl.idGenerator.GenerateId(videoCoverIdGenerateKey)
		path := fmt.Sprintf("covers/%s.jpg", filename)
		impl.pictureStorage.Upload(bucketVideo, path, coverDTO.File, coverDTO.Size)
		video.CoverUrl = path
	}

	// 更新写入数据库
	err = impl.videoMapper.Update(id, video, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("update video [%d] failed", id)))
	}

	return lastCoverUrl
}

func (impl *VideoBizServiceImpl) DeleteVideo(id int64, tx ...orm.TxOrmer) (string, string) {
	video, err := impl.videoMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("video [%d] doesn't exist", id)))
	}
	err = impl.videoMapper.DeleteById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("delete video [%d] failed", id)))
	}

	return video.CoverUrl, video.VideoUrl
}

func (impl *VideoBizServiceImpl) PlayVideo(id int64, rangeHeader ...string) videoDTO.VideoFileDTO {
	video, err := impl.videoMapper.SelectById(id)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("video [%d] doesn't exist", id)))
	}
	do := impl.videoStorage.Download(bucketVideo, video.VideoUrl, rangeHeader...)
	return videoDTO.VideoFileDTO{
		Header: do.Header,
		Reader: do.Reader,
	}
}

func (impl *VideoBizServiceImpl) RemoveVideoCover(path string) {
	impl.pictureStorage.Remove(bucketVideo, path)
}

func (impl *VideoBizServiceImpl) RemoveVideoFile(path string) {
	impl.videoStorage.Remove(bucketVideo, path)
}
