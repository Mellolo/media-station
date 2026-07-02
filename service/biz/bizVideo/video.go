package bizVideo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

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
	"k8s.io/apimachinery/pkg/util/sets"
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

	// 流式转码相关方法
	StreamVideoToHLS(ctx contextDTO.ContextDTO, id int64) map[string]string
	ServeHLSSegment(sessionID string, filename string) (string, error)

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

	// 1. 保存上传文件到临时位置用于格式检测
	tempFile, err := ioutil.TempFile("", "upload_*")
	if err != nil {
		panic(errors.WrapError(err, "create temp file failed"))
	}
	tempPath := tempFile.Name()

	// 将上传文件内容复制到临时文件
	_, err = io.Copy(tempFile, videoDTO.File)
	if err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		panic(errors.WrapError(err, "save upload file to temp failed"))
	}
	tempFile.Close()

	// 2. 使用 ffprobe 检测真实格式
	formatInfo, err := detectVideoFormat(tempPath)
	if err != nil {
		logs.Warn(fmt.Sprintf("格式检测失败: %v, 使用默认扩展名 mp4", err))
		formatInfo = map[string]interface{}{"format_name": "mp4"}
	}

	// 3. 确定文件扩展名
	var formatName string
	if fn, ok := formatInfo["format_name"].(string); ok {
		formatName = fn
	}
	originalExt := GetExtensionFromFormat(formatName)

	// 4. 生成唯一文件名并保存到 MinIO（保留原始扩展名）
	filename := impl.idGenerator.GenerateId(videoIdGenerateKey)
	path := fmt.Sprintf("videos/%s.%s", filename, originalExt)

	// 重新打开临时文件用于上传
	fileReader, err := os.Open(tempPath)
	if err != nil {
		os.Remove(tempPath)
		panic(errors.WrapError(err, "open temp file failed"))
	}
	fileInfo, err := os.Stat(tempPath)
	if err != nil {
		fileReader.Close()
		os.Remove(tempPath)
		panic(errors.WrapError(err, "get temp file info failed"))
	}

	impl.videoStorage.Upload(bucketVideo, path, fileReader, fileInfo.Size())
	fileReader.Close()
	video.VideoUrl = path

	// 5. 生成封面和获取时长
	url := impl.videoStorage.GetStreamURL(bucketVideo, path, 3*time.Minute)
	filename = impl.idGenerator.GenerateId(videoCoverIdGenerateKey)
	path = fmt.Sprintf("covers/%s.jpg", filename)
	data, err := videoUtil.CaptureScreenshotFromStreamURL(url)
	if err != nil {
		os.Remove(tempPath)
		panic(errors.WrapError(err, "capture video thumbnail failed"))
	}
	impl.pictureStorage.Upload(bucketVideo, path, io.NopCloser(bytes.NewReader(data)), int64(len(data)))
	video.CoverUrl = path

	// 6. 获取时长
	seconds, err := videoUtil.GetVideoDuration(impl.videoStorage.GetStreamURL(bucketVideo, video.VideoUrl, time.Minute))
	if err != nil {
		os.Remove(tempPath)
		panic(errors.WrapError(err, "get video duration failed"))
	}
	video.Duration = seconds

	// 7. 保存到数据库
	id, err := impl.videoMapper.Insert(video, tx...)
	if err != nil {
		os.Remove(tempPath)
		panic(errors.WrapError(err, "create video failed"))
	}

	// 8. 清理临时文件
	os.Remove(tempPath)

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

	// MinIO 可能返回 application/octet-stream，根据文件扩展名设置正确的 Content-Type
	contentType := do.Header.Get("Content-Type")
	if contentType == "" || contentType == "application/octet-stream" {
		lowerUrl := strings.ToLower(video.VideoUrl)
		switch {
		case strings.HasSuffix(lowerUrl, ".mp4") || strings.HasSuffix(lowerUrl, ".m4v"):
			do.Header.Set("Content-Type", "video/mp4")
		case strings.HasSuffix(lowerUrl, ".webm"):
			do.Header.Set("Content-Type", "video/webm")
		case strings.HasSuffix(lowerUrl, ".mkv"):
			do.Header.Set("Content-Type", "video/x-matroska")
		case strings.HasSuffix(lowerUrl, ".avi"):
			do.Header.Set("Content-Type", "video/x-msvideo")
		case strings.HasSuffix(lowerUrl, ".mov"):
			do.Header.Set("Content-Type", "video/quicktime")
		case strings.HasSuffix(lowerUrl, ".flv"):
			do.Header.Set("Content-Type", "video/x-flv")
		case strings.HasSuffix(lowerUrl, ".wmv"):
			do.Header.Set("Content-Type", "video/x-ms-wmv")
		case strings.HasSuffix(lowerUrl, ".ts"):
			do.Header.Set("Content-Type", "video/mp2t")
		}
	}

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

// detectVideoFormat 使用 ffprobe 检测视频格式
func detectVideoFormat(filePath string) (map[string]interface{}, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		filePath)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe执行失败: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	// 返回 format 子对象，而不是根级 JSON
	if format, ok := result["format"].(map[string]interface{}); ok {
		return format, nil
	}

	return result, nil
}
