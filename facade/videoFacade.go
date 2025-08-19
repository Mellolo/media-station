package facade

import (
	"github.com/Mellolo/common/errors"
	"github.com/Mellolo/common/utils/jsonUtil"
	"github.com/Mellolo/media-station/enum"
	"github.com/Mellolo/media-station/models/dto/fileDTO"
	"github.com/Mellolo/media-station/models/dto/performDTO"
	"github.com/Mellolo/media-station/models/dto/tagDTO"
	"github.com/Mellolo/media-station/models/dto/videoDTO"
	"github.com/Mellolo/media-station/models/vo/videoVO"
	"github.com/Mellolo/media-station/service/biz/bizActor"
	"github.com/Mellolo/media-station/service/biz/bizPerform"
	"github.com/Mellolo/media-station/service/biz/bizTag"
	"github.com/Mellolo/media-station/service/biz/bizVideo"
	"github.com/Mellolo/media-station/storage/db"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"k8s.io/apimachinery/pkg/util/sets"
)

type VideoFacade struct {
	AbstractFacade
	videoBizService   bizVideo.VideoBizService
	actorBizService   bizActor.ActorBizService
	performBizService bizPerform.PerformBizService
	tagBizService     bizTag.TagBizService
}

func NewVideoFacade() *VideoFacade {
	return &VideoFacade{
		videoBizService:   bizVideo.NewVideoBizService(),
		actorBizService:   bizActor.NewActorBizService(),
		performBizService: bizPerform.NewPerformBizService(),
		tagBizService:     bizTag.NewTagBizService(),
	}
}

func (impl *VideoFacade) SearchVideo(c *web.Controller) []videoVO.VideoItemVO {
	// 上下文
	ctx := impl.GetContext(c)
	keyword := c.GetString("keyword", "")
	// 演员
	actors := impl.GetStringAsInt64List(c, "actors")
	actors = sets.NewInt64(actors...).List()
	// tags
	tags := impl.GetStringAsStringList(c, "tags")
	tags = sets.NewString(tags...).List()

	var voList []videoVO.VideoItemVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		// 仅有关键词
		if len(actors) == 0 && len(tags) == 0 {
			items := impl.videoBizService.SearchVideoByKeyword(ctx, keyword, tx)
			for _, item := range items {
				voList = append(voList, videoVO.VideoItemVO{
					Id:              item.Id,
					Name:            item.Name,
					Duration:        item.Duration,
					PermissionLevel: item.PermissionLevel,
				})
			}
			return
		}

		searchDTO := videoDTO.VideoSearchDTO{
			Keyword: keyword,
		}

		if len(actors) > 0 && len(tags) > 0 {
			videoIdSet := sets.NewInt64(impl.performBizService.SelectArtByActor(ctx, enum.ArtVideo, actors, tx)...)
			videoIdSet = videoIdSet.Intersection(sets.NewInt64(impl.tagBizService.SelectArtByTag(ctx, enum.ArtVideo, tags, tx)...))
			searchDTO.Ids = videoIdSet.List()
		} else if len(actors) > 0 {
			searchDTO.Ids = impl.performBizService.SelectArtByActor(ctx, enum.ArtVideo, actors, tx)
		} else {
			searchDTO.Ids = impl.tagBizService.SelectArtByTag(ctx, enum.ArtVideo, tags, tx)
		}

		items := impl.videoBizService.SearchVideo(ctx, searchDTO, tx)

		for _, item := range items {
			voList = append(voList, videoVO.VideoItemVO{
				Id:              item.Id,
				Name:            item.Name,
				Duration:        item.Duration,
				PermissionLevel: item.PermissionLevel,
			})
		}
	})

	return voList
}

func (impl *VideoFacade) SearchVideoByTag(c *web.Controller) []videoVO.VideoItemVO {
	// 上下文
	ctx := impl.GetContext(c)
	// tag
	tagName := c.GetString("tag", "")

	var voList []videoVO.VideoItemVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		videoIds := impl.tagBizService.SelectArtByTag(ctx, enum.ArtVideo, []string{tagName}, tx)
		if len(videoIds) == 0 {
			return
		}

		items := impl.videoBizService.SearchVideo(ctx, videoDTO.VideoSearchDTO{
			Ids: videoIds,
		}, tx)

		for _, item := range items {
			voList = append(voList, videoVO.VideoItemVO{
				Id:              item.Id,
				Name:            item.Name,
				Duration:        item.Duration,
				PermissionLevel: item.PermissionLevel,
			})
		}
	})

	return voList
}

func (impl *VideoFacade) RecommendVideo(c *web.Controller) []videoVO.VideoItemVO {
	// 上下文
	ctx := impl.GetContext(c)

	var voList []videoVO.VideoItemVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		items := impl.videoBizService.SearchVideoByKeyword(ctx, "", tx)

		for _, item := range items {
			voList = append(voList, videoVO.VideoItemVO{
				Id:              item.Id,
				Name:            item.Name,
				Duration:        item.Duration,
				PermissionLevel: item.PermissionLevel,
			})
		}
	})

	return voList
}

func (impl *VideoFacade) GetVideoPage(c *web.Controller) videoVO.VideoPageVO {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetRestfulParamInt64(c, ":id")

	var vo videoVO.VideoPageVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		video := impl.videoBizService.GetVideo(ctx, id, tx)

		actorIds := impl.performBizService.SelectActorByArt(ctx, enum.ArtVideo, id, tx)
		var actors []videoVO.VideoActorVO
		for _, actorId := range actorIds {
			actor := impl.actorBizService.GetActor(ctx, actorId)
			actors = append(actors, videoVO.VideoActorVO{
				Id:   actor.Id,
				Name: actor.Name,
			})
		}

		tags := impl.tagBizService.SelectTagByArt(ctx, enum.ArtVideo, id, tx)

		vo = videoVO.VideoPageVO{
			Id:              video.Id,
			Name:            video.Name,
			Description:     video.Description,
			Actors:          actors,
			Tags:            tags,
			PermissionLevel: video.PermissionLevel,
		}
	})

	return vo
}

func (impl *VideoFacade) GetVideoCover(c *web.Controller) videoVO.VideoCoverFileVO {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetRestfulParamInt64(c, ":id")

	var vo videoVO.VideoCoverFileVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		cover := impl.videoBizService.GetVideoCover(ctx, id)
		vo = videoVO.VideoCoverFileVO{
			Reader: cover.Reader,
			Header: cover.Header,
		}
	})

	return vo
}

func (impl *VideoFacade) UploadVideo(c *web.Controller) {
	// 上下文
	ctx := impl.GetContext(c)
	// 名称
	name := impl.GetStringNotEmpty(c, "name")
	// 描述
	description := c.GetString("description", "")
	// 演员
	actorIds := impl.GetStringAsInt64List(c, "actorIds")
	actorIds = sets.NewInt64(actorIds...).List()
	// tags
	tags := impl.GetStringAsStringList(c, "tags")
	tags = sets.NewString(tags...).List()
	// 上传者
	uploader := ctx.UserClaim.Username
	// 权限
	permissionLevel := c.GetString("permissionLevel", "")

	// 视频文件
	reader, header, err := c.GetFile("file")
	if err != nil {
		panic(errors.WrapError(err, "get video file failed"))
	}
	videoFileDTO := fileDTO.FileDTO{
		File: reader,
		Size: header.Size,
	}

	db.DoTransaction(func(tx orm.TxOrmer) {
		createDTO := videoDTO.VideoCreateDTO{
			Name:            name,
			Description:     description,
			Actors:          actorIds,
			Tags:            tags,
			Uploader:        uploader,
			PermissionLevel: permissionLevel,
		}
		// 创建视频
		id := impl.videoBizService.CreateVideo(ctx, createDTO, videoFileDTO)

		// 更新作品actor出演关系
		impl.performBizService.InsertOrUpdateActorsOfArt(ctx, performDTO.ArtPerformDTO{
			ArtType:  enum.ArtVideo,
			ArtId:    id,
			ActorIds: actorIds,
		}, tx)
		// 更新作品tag
		impl.tagBizService.InsertOrUpdateTagsOfArt(ctx, tagDTO.ArtTagDTO{
			ArtType: enum.ArtVideo,
			ArtId:   id,
			Tags:    tags,
		}, tx)
	})
}

func (impl *VideoFacade) UpdateVideo(c *web.Controller) {
	// 上下文
	ctx := impl.GetContext(c)

	var requestBody struct {
		videoDTO.VideoUpdateDTO
		ActorIds []int64  `json:"actorIds"`
		Tags     []string `json:"tags"`
	}
	jsonUtil.UnmarshalJsonString(string(c.Ctx.Input.RequestBody), &requestBody)

	db.DoTransaction(func(tx orm.TxOrmer) {
		updateDTO := requestBody.VideoUpdateDTO
		// 更新视频
		_ = impl.videoBizService.UpdateVideo(ctx, updateDTO, tx)
		// 更新作品actor出演关系
		impl.performBizService.InsertOrUpdateActorsOfArt(ctx, performDTO.ArtPerformDTO{
			ArtType:  enum.ArtVideo,
			ArtId:    updateDTO.Id,
			ActorIds: requestBody.ActorIds,
		}, tx)
		// 更新作品tag
		impl.tagBizService.InsertOrUpdateTagsOfArt(ctx, tagDTO.ArtTagDTO{
			ArtType: enum.ArtVideo,
			ArtId:   updateDTO.Id,
			Tags:    requestBody.Tags,
		}, tx)
	})
}

func (impl *VideoFacade) DeleteVideo(c *web.Controller) {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetRestfulParamInt64(c, ":id")

	var coverUrl, videoUrl string
	db.DoTransaction(func(tx orm.TxOrmer) {
		// 更新视频
		video := impl.videoBizService.DeleteVideo(ctx, id, tx)
		coverUrl, videoUrl = video.CoverUrl, video.VideoUrl
		// 更新作品actor出演关系
		impl.performBizService.DeleteArt(ctx, enum.ArtVideo, id, tx)
		// 更新作品tag
		impl.tagBizService.DeleteArt(ctx, enum.ArtVideo, id, tx)
	})

	go func() {
		if coverUrl != "" {
			impl.videoBizService.RemoveVideoCover(ctx, coverUrl)
		}
		if videoUrl != "" {
			impl.videoBizService.RemoveVideoFile(ctx, videoUrl)
		}
	}()
}

func (impl *VideoFacade) PlayVideo(c *web.Controller) videoVO.VideoFileVO {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetRestfulParamInt64(c, ":id")

	dto := impl.videoBizService.PlayVideo(ctx, id, c.Ctx.Request.Header["Range"]...)
	return videoVO.VideoFileVO{
		Header: dto.Header,
		Reader: dto.Reader,
	}
}
