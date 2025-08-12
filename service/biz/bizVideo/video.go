package bizVideo

import (
	"bytes"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/mellolo/common/errors"
	"github.com/mellolo/common/utils/videoUtil"
	"io"
	"k8s.io/apimachinery/pkg/util/sets"
	"media-station/enum"
	"media-station/generator"
	"media-station/models/do/videoDO"
	"media-station/models/dto/fileDTO"
	"media-station/models/dto/videoDTO"
	"media-station/storage/cache"
	"media-station/storage/db"
	"media-station/storage/oss"
	"time"
)

const (
	videoIdGenerateKey      = "video"
	videoCoverIdGenerateKey = "videoCover"
	bucketVideo             = "video"
	videoDurationCacheKey   = "videoDuration"
)

type VideoBizService interface {
	GetVideoPage(id int64, tx ...orm.TxOrmer) videoDTO.VideoPageDTO
	GetVideoCover(id int64, tx ...orm.TxOrmer) videoDTO.VideoCoverDTO
	SearchVideo(searchDTO videoDTO.VideoSearchDTO, tx ...orm.TxOrmer) []videoDTO.VideoItemDTO
	CreateVideo(createDTO videoDTO.VideoCreateDTO, videoDTO fileDTO.FileDTO, tx ...orm.TxOrmer) int64
	UpdateVideo(id int64, updateDTO videoDTO.VideoUpdateDTO, tx ...orm.TxOrmer)
	DeleteVideo(id int64, tx ...orm.TxOrmer) (string, string)
	PlayVideo(id int64, rangeHeader ...string) videoDTO.VideoFileDTO

	RemoveVideoCover(path string)
	RemoveVideoFile(path string)
}

func NewVideoBizService() *VideoBizServiceImpl {
	return &VideoBizServiceImpl{
		videoMapper:      db.NewVideoMapper(),
		idGenerator:      generator.NewIdGenerator(),
		pictureStorage:   oss.NewPictureStorage(),
		videoStorage:     oss.NewVideoStorage(),
		distributedCache: cache.NewDistributedCache(),
	}
}

type VideoBizServiceImpl struct {
	videoMapper      db.VideoMapper
	idGenerator      generator.IdGenerator
	pictureStorage   oss.PictureStorage
	videoStorage     oss.VideoStorage
	distributedCache cache.DistributedCache
}

func (impl *VideoBizServiceImpl) GetVideoPage(id int64, tx ...orm.TxOrmer) videoDTO.VideoPageDTO {
	video, err := impl.videoMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("get video [%d] failed", id)))
	}
	return videoDTO.VideoPageDTO{
		Id:              video.Id,
		Name:            video.Name,
		Description:     video.Description,
		Actors:          video.Actors,
		Tags:            video.Tags,
		Uploader:        video.Uploader,
		CoverUrl:        video.CoverUrl,
		VideoUrl:        video.VideoUrl,
		PermissionLevel: video.PermissionLevel,
	}
}

func (impl *VideoBizServiceImpl) GetVideoCover(id int64, tx ...orm.TxOrmer) videoDTO.VideoCoverDTO {
	video, err := impl.videoMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("get video [%d] failed", id)))
	}

	cover := impl.videoStorage.GetCover(bucketVideo, video.CoverUrl)
	return videoDTO.VideoCoverDTO{
		Reader: cover.Reader,
		Header: cover.Header,
	}
}

func (impl *VideoBizServiceImpl) SearchVideo(searchDTO videoDTO.VideoSearchDTO, tx ...orm.TxOrmer) []videoDTO.VideoItemDTO {
	// 读取数据库
	var videoDOList []*videoDO.VideoDO
	if searchDTO.Keyword == "" {
		doList, err := impl.videoMapper.SelectAllLimit(200, tx...)
		if err != nil {
			panic(errors.WrapError(err, "select all video error"))
		}
		videoDOList = append(videoDOList, doList...)
	} else {
		doList, err := impl.videoMapper.SelectByKeyword(searchDTO.Keyword)
		if err != nil {
			panic(errors.WrapError(err, fmt.Sprintf("select video by keyword [%s] error", searchDTO.Keyword)))
		}
		videoDOList = append(videoDOList, doList...)
	}

	// 再筛选，并转化为DTO
	var items []videoDTO.VideoItemDTO
	for _, do := range videoDOList {
		if do.PermissionLevel == enum.PermissionForbidden || do.PermissionLevel == enum.PermissionPrivate {
			continue
		}
		if len(do.Actors) > 0 && !sets.NewInt64(do.Actors...).HasAll(searchDTO.Actors...) {
			continue
		}
		if len(do.Tags) > 0 && !sets.NewString(do.Tags...).HasAll(searchDTO.Tags...) {
			continue
		}

		items = append(items, videoDTO.VideoItemDTO{
			Id:              do.Id,
			Name:            do.Name,
			CoverUrl:        do.CoverUrl,
			Duration:        impl.getVideoDuration(do.Id, do.VideoUrl),
			PermissionLevel: do.PermissionLevel,
		})
	}

	return items
}

func (impl *VideoBizServiceImpl) CreateVideo(createDTO videoDTO.VideoCreateDTO, videoDTO fileDTO.FileDTO, tx ...orm.TxOrmer) int64 {
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

	// 上传视频
	filename := impl.idGenerator.GenerateId(videoIdGenerateKey)
	path := fmt.Sprintf("videos/%s.mp4", filename)
	impl.videoStorage.Upload(bucketVideo, path, videoDTO.File, videoDTO.Size)
	video.VideoUrl = path

	// 上传封面
	url := impl.videoStorage.GetStreamURL(bucketVideo, path, 3*time.Minute)
	filename = impl.idGenerator.GenerateId(videoCoverIdGenerateKey)
	path = fmt.Sprintf("covers/%s.jpg", filename)
	data, err := videoUtil.CaptureScreenshotFromStreamURL(url)
	impl.pictureStorage.Upload(bucketVideo, path, io.NopCloser(bytes.NewReader(data)), int64(len(data)))
	video.CoverUrl = path

	id, err := impl.videoMapper.Insert(video, tx...)
	if err != nil {
		panic(errors.WrapError(err, "create video failed"))
	}
	return id
}

func (impl *VideoBizServiceImpl) UpdateVideo(id int64, updateDTO videoDTO.VideoUpdateDTO, tx ...orm.TxOrmer) {
	video, err := impl.videoMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("video [%d] doesn't exist", id)))
	}

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

	// 更新写入数据库
	err = impl.videoMapper.Update(id, video, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("update video [%d] failed", id)))
	}
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

func (impl *VideoBizServiceImpl) getVideoDuration(id int64, videoPath string) float64 {
	val, getErr := impl.distributedCache.Get(fmt.Sprintf("%s|%d", videoDurationCacheKey, id))
	if getErr != nil {
		logs.Error(fmt.Sprintf("error get video duration: %v", getErr))
	}
	if seconds, ok := val.(float64); ok {
		return seconds
	}
	seconds, err := videoUtil.GetVideoDuration(impl.videoStorage.GetStreamURL(bucketVideo, videoPath, time.Minute))
	if err == nil {
		setErr := impl.distributedCache.Set(fmt.Sprintf("%s|%d", videoDurationCacheKey, id), seconds, time.Hour)
		if setErr != nil {
			logs.Error(fmt.Sprintf("error set video duration: %v", setErr))
		}
	}

	return seconds
}
