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
	"media-station/models/do/userDO"
	"media-station/models/do/videoDO"
	"media-station/models/dto/contextDTO"
	"media-station/models/dto/fileDTO"
	"media-station/models/dto/videoDTO"
	"media-station/service/domain/domainPermission"
	"media-station/storage/cache"
	"media-station/storage/db"
	"media-station/storage/oss"
	"strconv"
	"time"
)

const (
	videoIdGenerateKey      = "video"
	videoCoverIdGenerateKey = "videoCover"
	bucketVideo             = "video"
	videoDurationCacheKey   = "videoDuration"
)

type VideoBizService interface {
	GetVideoPage(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) videoDTO.VideoPageDTO
	GetVideoCover(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) videoDTO.VideoCoverDTO
	SearchVideo(ctx contextDTO.ContextDTO, searchDTO videoDTO.VideoSearchDTO, tx ...orm.TxOrmer) []videoDTO.VideoItemDTO
	SearchVideoByActor(ctx contextDTO.ContextDTO, actorId int64, tx ...orm.TxOrmer) []videoDTO.VideoItemDTO
	SearchVideoByTag(ctx contextDTO.ContextDTO, tagName string, tx ...orm.TxOrmer) []videoDTO.VideoItemDTO
	CreateVideo(ctx contextDTO.ContextDTO, createDTO videoDTO.VideoCreateDTO, videoDTO fileDTO.FileDTO, tx ...orm.TxOrmer) int64
	UpdateVideo(ctx contextDTO.ContextDTO, id int64, updateDTO videoDTO.VideoUpdateDTO, tx ...orm.TxOrmer) videoDTO.VideoPageDTO
	RemoveActor(ctx contextDTO.ContextDTO, id int64, actorIds []int64, tx ...orm.TxOrmer)
	RemoveTag(ctx contextDTO.ContextDTO, id int64, tags []string, tx ...orm.TxOrmer)
	DeleteVideo(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) videoDTO.VideoPageDTO
	PlayVideo(ctx contextDTO.ContextDTO, id int64, rangeHeader ...string) videoDTO.VideoFileDTO

	RemoveVideoCover(ctx contextDTO.ContextDTO, path string)
	RemoveVideoFile(ctx contextDTO.ContextDTO, path string)
}

func NewVideoBizService() *VideoBizServiceImpl {
	return &VideoBizServiceImpl{
		userMapper:              db.NewUserMapper(),
		actorMapper:             db.NewActorMapper(),
		tagMapper:               db.NewTagMapper(),
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
	actorMapper      db.ActorMapper
	tagMapper        db.TagMapper
	videoMapper      db.VideoMapper
	idGenerator      generator.IdGenerator
	pictureStorage   oss.PictureStorage
	videoStorage     oss.VideoStorage
	distributedCache cache.DistributedCache

	permissionDomainService domainPermission.PermissionDomainService
}

func (impl *VideoBizServiceImpl) GetVideoPage(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) videoDTO.VideoPageDTO {
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

func (impl *VideoBizServiceImpl) SearchVideo(ctx contextDTO.ContextDTO, searchDTO videoDTO.VideoSearchDTO, tx ...orm.TxOrmer) []videoDTO.VideoItemDTO {
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
	var user = new(userDO.UserDO)
	if ctx.UserClaim.Username != "" {
		user, _ = impl.userMapper.SelectByUsername(ctx.UserClaim.Username, tx...)
	}
	var items []videoDTO.VideoItemDTO
	for _, do := range videoDOList {
		if !impl.permissionDomainService.IsAccessAllowed(*user, do.Uploader, do.PermissionLevel) {
			continue
		}
		if len(searchDTO.Actors) > 0 && !sets.NewInt64(do.Actors...).HasAll(searchDTO.Actors...) {
			continue
		}
		if len(searchDTO.Tags) > 0 && !sets.NewString(do.Tags...).HasAll(searchDTO.Tags...) {
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

func (impl *VideoBizServiceImpl) SearchVideoByActor(ctx contextDTO.ContextDTO, actorId int64, tx ...orm.TxOrmer) []videoDTO.VideoItemDTO {
	actor, actorErr := impl.actorMapper.SelectById(actorId, tx...)
	if actorErr != nil || actor == nil {
		logs.Error("get actor [%d] failed", actorId) // 不报错，记录日志
		return []videoDTO.VideoItemDTO{}
	}

	var user = new(userDO.UserDO)
	if ctx.UserClaim.Username != "" {
		user, _ = impl.userMapper.SelectByUsername(ctx.UserClaim.Username, tx...)
	}
	// 筛选，并转化为DTO
	var items []videoDTO.VideoItemDTO
	for _, videoId := range actor.Art.VideoIds {
		do, videoErr := impl.videoMapper.SelectById(videoId, tx...)
		if videoErr != nil {
			logs.Error("select video [%d] error", videoId) // 不报错，记录日志
		}
		if do == nil {
			continue
		}

		if !impl.permissionDomainService.IsAccessAllowed(*user, do.Uploader, do.PermissionLevel) {
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

func (impl *VideoBizServiceImpl) SearchVideoByTag(ctx contextDTO.ContextDTO, tagName string, tx ...orm.TxOrmer) []videoDTO.VideoItemDTO {
	tag, tagErr := impl.tagMapper.SelectByName(tagName, tx...)
	if tagErr != nil || tag == nil {
		logs.Error("get tag [%d] failed", tagName) // 不报错，记录日志
		return []videoDTO.VideoItemDTO{}
	}

	var user = new(userDO.UserDO)
	if ctx.UserClaim.Username != "" {
		user, _ = impl.userMapper.SelectByUsername(ctx.UserClaim.Username, tx...)
	}
	// 筛选，并转化为DTO
	var items []videoDTO.VideoItemDTO
	for _, videoId := range tag.Art.VideoIds {
		do, videoErr := impl.videoMapper.SelectById(videoId, tx...)
		if videoErr != nil {
			logs.Error("select video [%d] error", videoId) // 不报错，记录日志
		}
		if do == nil {
			continue
		}

		if !impl.permissionDomainService.IsAccessAllowed(*user, do.Uploader, do.PermissionLevel) {
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

func (impl *VideoBizServiceImpl) CreateVideo(ctx contextDTO.ContextDTO, createDTO videoDTO.VideoCreateDTO, videoDTO fileDTO.FileDTO, tx ...orm.TxOrmer) int64 {
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

func (impl *VideoBizServiceImpl) UpdateVideo(ctx contextDTO.ContextDTO, id int64, updateDTO videoDTO.VideoUpdateDTO, tx ...orm.TxOrmer) videoDTO.VideoPageDTO {
	video, err := impl.videoMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("video [%d] doesn't exist", id)))
	}
	origin := video.DeepCopy()

	// 更新video
	video.Name = updateDTO.Name
	video.Description = updateDTO.Description
	video.Actors = updateDTO.ActorIds
	video.Tags = updateDTO.Tags
	if sets.NewString(enum.PermissionLevels...).Has(updateDTO.PermissionLevel) {
		video.PermissionLevel = updateDTO.PermissionLevel
	}

	// 更新写入数据库
	err = impl.videoMapper.Update(id, video, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("update video [%d] failed", id)))
	}

	return videoDTO.VideoPageDTO{
		Id:              origin.Id,
		Name:            origin.Name,
		Description:     origin.Description,
		Actors:          origin.Actors,
		Tags:            origin.Tags,
		Uploader:        origin.Uploader,
		CoverUrl:        origin.CoverUrl,
		VideoUrl:        origin.VideoUrl,
		PermissionLevel: origin.PermissionLevel,
	}
}

func (impl *VideoBizServiceImpl) RemoveActor(ctx contextDTO.ContextDTO, id int64, actorIds []int64, tx ...orm.TxOrmer) {
	video, err := impl.videoMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("video [%d] doesn't exist", id)))
	}

	// 更新video
	video.Actors = sets.NewInt64(video.Actors...).Delete(actorIds...).List()

	// 更新写入数据库
	err = impl.videoMapper.Update(id, video, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("update video [%d] failed", id)))
	}
}

func (impl *VideoBizServiceImpl) RemoveTag(ctx contextDTO.ContextDTO, id int64, tags []string, tx ...orm.TxOrmer) {
	video, err := impl.videoMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("video [%d] doesn't exist", id)))
	}

	// 更新video
	video.Tags = sets.NewString(video.Tags...).Delete(tags...).List()

	// 更新写入数据库
	err = impl.videoMapper.Update(id, video, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("update video [%d] failed", id)))
	}
}

func (impl *VideoBizServiceImpl) DeleteVideo(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) videoDTO.VideoPageDTO {
	video, err := impl.videoMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("video [%d] doesn't exist", id)))
	}
	err = impl.videoMapper.DeleteById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("delete video [%d] failed", id)))
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
