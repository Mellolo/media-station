package facade

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/mellolo/common/errors"
	"github.com/mellolo/common/utils/jsonUtil"
	"k8s.io/apimachinery/pkg/util/sets"
	"media-station/models/dto/actorDTO"
	"media-station/models/dto/fileDTO"
	"media-station/models/dto/tagDTO"
	"media-station/models/dto/videoDTO"
	"media-station/models/vo/videoVO"
	"media-station/service/biz/bizActor"
	"media-station/service/biz/bizTag"
	"media-station/service/biz/bizVideo"
	"media-station/storage/db"
)

type VideoFacade struct {
	AbstractFacade
	videoBizService bizVideo.VideoBizService
	actorBizService bizActor.ActorBizService
	tagBizService   bizTag.TagBizService
}

func NewVideoFacade() *VideoFacade {
	return &VideoFacade{
		videoBizService: bizVideo.NewVideoBizService(),
		actorBizService: bizActor.NewActorBizService(),
		tagBizService:   bizTag.NewTagBizService(),
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
		searchDTO := videoDTO.VideoSearchDTO{
			Keyword: keyword,
			Actors:  actors,
			Tags:    tags,
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
		items := impl.videoBizService.SearchVideoByTag(ctx, tagName, tx)

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
		page := impl.videoBizService.GetVideoPage(ctx, id, tx)
		var actors []videoVO.VideoActorVO
		for _, actorId := range page.Actors {
			actorPage := impl.actorBizService.GetActorPage(ctx, actorId)
			actors = append(actors, videoVO.VideoActorVO{
				Id:   actorPage.Id,
				Name: actorPage.Name,
			})
		}
		vo = videoVO.VideoPageVO{
			Id:              page.Id,
			Name:            page.Name,
			Description:     page.Description,
			Actors:          actors,
			Tags:            page.Tags,
			PermissionLevel: page.PermissionLevel,
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
	actors := impl.GetStringAsInt64List(c, "actors")
	actors = sets.NewInt64(actors...).List()
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
			Actors:          actors,
			Tags:            tags,
			Uploader:        uploader,
			PermissionLevel: permissionLevel,
		}
		// 创建视频
		id := impl.videoBizService.CreateVideo(ctx, createDTO, videoFileDTO)
		// 更新actor作品
		for _, actorId := range createDTO.Actors {
			updateDTO := actorDTO.ActorArtDTO{
				VideoIds: []int64{id},
			}
			impl.actorBizService.AddArt(ctx, actorId, updateDTO, tx)
		}
		// 更新tag作品
		for _, tagName := range createDTO.Tags {
			tag := tagDTO.TagCreateOrUpdateDTO{
				Name:    tagName,
				Creator: uploader,
				Details: tagDTO.TagDetailsDTO{
					VideoIds: []int64{id},
				},
			}
			impl.tagBizService.AddArt(ctx, tag, tx)
		}
	})
}

func (impl *VideoFacade) UpdateVideo(c *web.Controller) {
	// 上下文
	ctx := impl.GetContext(c)

	var dto videoDTO.VideoUpdateDTO
	jsonUtil.UnmarshalJsonString(string(c.Ctx.Input.RequestBody), &dto)

	db.DoTransaction(func(tx orm.TxOrmer) {
		id := dto.Id
		// 更新视频
		impl.videoBizService.UpdateVideo(ctx, id, dto, tx)
		// 更新actor作品
		for _, actorId := range dto.Actors {
			artDTO := actorDTO.ActorArtDTO{
				VideoIds: []int64{id},
			}
			impl.actorBizService.AddArt(ctx, actorId, artDTO, tx)
		}
		// 更新tag作品
		for _, tagName := range dto.Tags {
			updateDTO := tagDTO.TagCreateOrUpdateDTO{
				Name: tagName,
				Details: tagDTO.TagDetailsDTO{
					VideoIds: []int64{id},
				},
			}
			impl.tagBizService.AddArt(ctx, updateDTO, tx)
		}
	})
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
