package bizVideo

import (
	"bytes"
	"fmt"
	"github.com/Mellolo/common/errors"
	"github.com/Mellolo/common/utils/videoUtil"
	"github.com/Mellolo/media-station/enum"
	"github.com/Mellolo/media-station/generator"
	"github.com/Mellolo/media-station/models/do/userDO"
	"github.com/Mellolo/media-station/models/do/videoDO"
	"github.com/Mellolo/media-station/models/dto/contextDTO"
	"github.com/Mellolo/media-station/models/dto/fileDTO"
	"github.com/Mellolo/media-station/models/dto/videoDTO"
	"github.com/Mellolo/media-station/service/domain/domainPermission"
	"github.com/Mellolo/media-station/storage/cache"
	"github.com/Mellolo/media-station/storage/db"
	"github.com/Mellolo/media-station/storage/oss"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"io"
	"k8s.io/apimachinery/pkg/util/sets"
	"strconv"
	"strings"
	"time"
)

const (
	videoIdGenerateKey      = "video"
	videoCoverIdGenerateKey = "videoCover"
	bucketVideo             = "video"
	videoDurationCacheKey   = "videoDuration"
)

type VideoBizService interface {
	GetVideo(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) videoDTO.VideoDTO
	GetVideoCover(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) videoDTO.VideoCoverDTO
	SearchVideo(ctx contextDTO.ContextDTO, searchDTO videoDTO.VideoSearchDTO, tx ...orm.TxOrmer) []videoDTO.VideoDTO
	SearchVideoByKeyword(ctx contextDTO.ContextDTO, keyword string, tx ...orm.TxOrmer) []videoDTO.VideoDTO
	CreateVideo(ctx contextDTO.ContextDTO, createDTO videoDTO.VideoCreateDTO, videoDTO fileDTO.FileDTO, tx ...orm.TxOrmer) int64
	UpdateVideo(ctx contextDTO.ContextDTO, updateDTO videoDTO.VideoUpdateDTO, tx ...orm.TxOrmer) videoDTO.VideoDTO
	DeleteVideo(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) videoDTO.VideoDTO
	PlayVideo(ctx contextDTO.ContextDTO, id int64, rangeHeader ...string) videoDTO.VideoFileDTO

	RemoveVideoCover(ctx contextDTO.ContextDTO, path string)
	RemoveVideoFile(ctx contextDTO.ContextDTO, path string)
}

func NewVideoBizService() *VideoBizServiceImpl {
	return &VideoBizServiceImpl{
		userMapper:              db.NewUserMapper(),
		videoMapper:             db.NewVideoMapper(),
		idGenerator:             generator.NewIdGenerator(),
		pictureStorage:          oss.NewPictureStorage(),
		videoStorage:            oss.NewVideoStorage(),
		distributedCache:        cache.NewDistributedCache(),
		permissionDomainService: domainPermission.NewPermissionDomainService(),
	}
}

type VideoBizServiceImpl struct {
	userMapper       db.UserMapper
	videoMapper      db.VideoMapper
	idGenerator      generator.IdGenerator
	pictureStorage   oss.PictureStorage
	videoStorage     oss.VideoStorage
	distributedCache cache.DistributedCache

	permissionDomainService domainPermission.PermissionDomainService
}

func (impl *VideoBizServiceImpl) GetVideo(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) videoDTO.VideoDTO {
	video, err := impl.videoMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("get video [%d] failed", id)))
	}
	return impl.convertVideoDO2VideoDTO(video)
}

func (impl *VideoBizServiceImpl) GetVideoCover(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) videoDTO.VideoCoverDTO {
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

func (impl *VideoBizServiceImpl) SearchVideo(ctx contextDTO.ContextDTO, searchDTO videoDTO.VideoSearchDTO, tx ...orm.TxOrmer) []videoDTO.VideoDTO {
	var videoDOList []videoDO.VideoDO
	for _, id := range searchDTO.Ids {
		do, err := impl.videoMapper.SelectById(id, tx...)
		if err != nil {
			panic(errors.WrapError(err, fmt.Sprintf("select video by id [%d] error", id)))
		}
		if searchDTO.Keyword != "" && !strings.Contains(do.Name, searchDTO.Keyword) && !strings.Contains(do.Description, searchDTO.Keyword) {
			continue
		}
		videoDOList = append(videoDOList, do)
	}

	// 再筛选，并转化为DTO
	var user userDO.UserDO
	if ctx.UserClaim.Username != "" {
		user, _ = impl.userMapper.SelectByUsername(ctx.UserClaim.Username, tx...)
	}
	var items []videoDTO.VideoDTO
	for _, do := range videoDOList {
		if !impl.permissionDomainService.IsVisible(user, do.Uploader, do.PermissionLevel) {
			continue
		}

		item := impl.convertVideoDO2VideoDTO(do)
		if item.Duration == 0 {
			item.Duration = impl.getVideoDuration(do.Id, do.VideoUrl)
		}
		items = append(items, item)
	}

	return items
}

func (impl *VideoBizServiceImpl) SearchVideoByKeyword(ctx contextDTO.ContextDTO, keyword string, tx ...orm.TxOrmer) []videoDTO.VideoDTO {
	var videoDOList []videoDO.VideoDO
	if keyword != "" {
		doList, err := impl.videoMapper.SelectByKeyword(keyword)
		if err != nil {
			panic(errors.WrapError(err, fmt.Sprintf("select video by keyword [%s] error", keyword)))
		}
		videoDOList = append(videoDOList, doList...)
	} else {
		doList, err := impl.videoMapper.SelectAllLimit(200, tx...)
		if err != nil {
			panic(errors.WrapError(err, "select all video error"))
		}
		videoDOList = append(videoDOList, doList...)
	}

	// 再筛选，并转化为DTO
	var user userDO.UserDO
	if ctx.UserClaim.Username != "" {
		user, _ = impl.userMapper.SelectByUsername(ctx.UserClaim.Username, tx...)
	}
	var items []videoDTO.VideoDTO
	for _, do := range videoDOList {
		if !impl.permissionDomainService.IsVisible(user, do.Uploader, do.PermissionLevel) {
			continue
		}

		item := impl.convertVideoDO2VideoDTO(do)
		if item.Duration == 0 {
			item.Duration = impl.getVideoDuration(do.Id, do.VideoUrl)
		}
		items = append(items, item)
	}

	return items
}

func (impl *VideoBizServiceImpl) CreateVideo(ctx contextDTO.ContextDTO, createDTO videoDTO.VideoCreateDTO, videoDTO fileDTO.FileDTO, tx ...orm.TxOrmer) int64 {
	if videoDTO.File == nil {
		panic(errors.NewError("video file is empty"))
	}

	video := videoDO.VideoDO{
		Name:        createDTO.Name,
		Description: createDTO.Description,
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
	if err != nil {
		panic(errors.WrapError(err, "capture video thumbnail failed"))
	}
	impl.pictureStorage.Upload(bucketVideo, path, io.NopCloser(bytes.NewReader(data)), int64(len(data)))
	video.CoverUrl = path

	// 时长
	seconds, err := videoUtil.GetVideoDuration(impl.videoStorage.GetStreamURL(bucketVideo, video.VideoUrl, time.Minute))
	if err != nil {
		panic(errors.WrapError(err, "get video duration failed"))
	}
	video.Duration = seconds

	id, err := impl.videoMapper.Insert(video, tx...)
	if err != nil {
		panic(errors.WrapError(err, "create video failed"))
	}
	return id
}

func (impl *VideoBizServiceImpl) UpdateVideo(ctx contextDTO.ContextDTO, updateDTO videoDTO.VideoUpdateDTO, tx ...orm.TxOrmer) videoDTO.VideoDTO {
	video, err := impl.videoMapper.SelectById(updateDTO.Id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("video [%d] doesn't exist", updateDTO.Id)))
	}
	origin := *video.DeepCopy()

	// 更新video
	video.Name = updateDTO.Name
	video.Description = updateDTO.Description
	if sets.NewString(enum.PermissionLevels...).Has(updateDTO.PermissionLevel) {
		video.PermissionLevel = updateDTO.PermissionLevel
	}

	// 更新写入数据库
	err = impl.videoMapper.Update(updateDTO.Id, video, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("update video [%d] failed", updateDTO.Id)))
	}

	return impl.convertVideoDO2VideoDTO(origin)
}

func (impl *VideoBizServiceImpl) DeleteVideo(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) videoDTO.VideoDTO {
	video, err := impl.videoMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("video [%d] doesn't exist", id)))
	}
	err = impl.videoMapper.DeleteById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("delete video [%d] failed", id)))
	}

	return impl.convertVideoDO2VideoDTO(video)
}

func (impl *VideoBizServiceImpl) PlayVideo(ctx contextDTO.ContextDTO, id int64, rangeHeader ...string) videoDTO.VideoFileDTO {
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

func (impl *VideoBizServiceImpl) RemoveVideoCover(ctx contextDTO.ContextDTO, path string) {
	impl.pictureStorage.Remove(bucketVideo, path)
}

func (impl *VideoBizServiceImpl) RemoveVideoFile(ctx contextDTO.ContextDTO, path string) {
	impl.videoStorage.Remove(bucketVideo, path)
}

func (impl *VideoBizServiceImpl) getVideoDuration(id int64, videoPath string) float64 {
	val, getErr := impl.distributedCache.Get(fmt.Sprintf("%s|%d", videoDurationCacheKey, id))
	if getErr != nil {
		logs.Error(fmt.Sprintf("error get video duration: %v", getErr))
	}
	if durationStr, ok := val.(string); ok {
		if seconds, err := strconv.ParseFloat(durationStr, 64); err == nil {
			return seconds
		}
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

func (impl *VideoBizServiceImpl) convertVideoDO2VideoDTO(do videoDO.VideoDO) videoDTO.VideoDTO {
	return videoDTO.VideoDTO{
		Id:              do.Id,
		Name:            do.Name,
		Description:     do.Description,
		Uploader:        do.Uploader,
		CoverUrl:        do.CoverUrl,
		VideoUrl:        do.VideoUrl,
		Duration:        do.Duration,
		PermissionLevel: do.PermissionLevel,
	}
}
