package facade

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/mellolo/common/errors"
	"media-station/enum"
	"media-station/models/dto/actorDTO"
	"media-station/models/dto/fileDTO"
	"media-station/models/dto/galleryDTO"
	"media-station/models/dto/videoDTO"
	"media-station/models/vo/actorVO"
	"media-station/models/vo/galleryVO"
	"media-station/models/vo/videoVO"
	"media-station/service/biz/bizActor"
	"media-station/service/biz/bizGallery"
	"media-station/service/biz/bizPerform"
	"media-station/service/biz/bizVideo"
	"media-station/storage/db"
)

type ActorFacade struct {
	AbstractFacade
	actorBizService   bizActor.ActorBizService
	videoBizService   bizVideo.VideoBizService
	galleryBizService bizGallery.GalleryBizService
	performBizService bizPerform.PerformBizService
}

func NewActorFacade() *ActorFacade {
	return &ActorFacade{
		actorBizService:   bizActor.NewActorBizService(),
		videoBizService:   bizVideo.NewVideoBizService(),
		galleryBizService: bizGallery.NewGalleryBizService(),
		performBizService: bizPerform.NewPerformBizService(),
	}
}

func (impl *ActorFacade) GetActorPage(c *web.Controller) actorVO.ActorPageVO {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetRestfulParamInt64(c, ":id")

	var vo actorVO.ActorPageVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		actor := impl.actorBizService.GetActor(ctx, id)
		vo = actorVO.ActorPageVO{
			Id:          actor.Id,
			Name:        actor.Name,
			Description: actor.Description,
			Creator:     actor.Creator,
		}

		videoIds := impl.performBizService.SelectArtByActor(ctx, enum.ArtVideo, []int64{id}, tx)

		videos := impl.videoBizService.SearchVideo(ctx, videoDTO.VideoSearchDTO{
			Ids: videoIds,
		}, tx)

		for _, item := range videos {
			vo.Videos = append(vo.Videos, videoVO.VideoItemVO{
				Id:              item.Id,
				Name:            item.Name,
				Duration:        item.Duration,
				PermissionLevel: item.PermissionLevel,
			})
		}

		galleryIds := impl.performBizService.SelectArtByActor(ctx, enum.ArtGallery, []int64{id}, tx)

		galleries := impl.galleryBizService.SearchGallery(ctx, galleryDTO.GallerySearchDTO{
			Ids: galleryIds,
		}, tx)

		for _, item := range galleries {
			vo.Galleries = append(vo.Galleries, galleryVO.GalleryItemVO{
				Id:              item.Id,
				Name:            item.Name,
				PageCount:       len(item.PicPaths),
				PermissionLevel: item.PermissionLevel,
			})
		}
	})

	return vo
}

func (impl *ActorFacade) GetActorCover(c *web.Controller) actorVO.ActorCoverFileVO {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetRestfulParamInt64(c, ":id")

	var vo actorVO.ActorCoverFileVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		cover := impl.actorBizService.GetActorCover(ctx, id)
		vo = actorVO.ActorCoverFileVO{
			Reader: cover.Reader,
			Header: cover.Header,
		}
	})

	return vo
}

func (impl *ActorFacade) CreateActor(c *web.Controller) int64 {
	// 上下文
	ctx := impl.GetContext(c)
	// 名称
	name := c.GetString("name", "")
	if name == "" {
		panic(errors.NewError("actor name is empty"))
	}
	// 描述
	description := c.GetString("description", "")

	// 创建者
	creator := ctx.UserClaim.Username
	if creator == "" {
		panic(errors.NewError("creator is empty"))
	}

	// 获取封面文件
	reader, size := impl.GetFileNotInvalid(c, "cover")

	var id int64
	db.DoTransaction(func(tx orm.TxOrmer) {
		createDTO := actorDTO.ActorCreateDTO{
			Name:        name,
			Description: description,
			Creator:     creator,
		}

		coverDTO := fileDTO.FileDTO{
			File: reader,
			Size: size,
		}
		id = impl.actorBizService.CreateActor(ctx, createDTO, coverDTO, tx)
	})

	return id
}

func (impl *ActorFacade) SearchActor(c *web.Controller) []actorVO.ActorItemVO {
	// 上下文
	ctx := impl.GetContext(c)

	var searchDTO actorDTO.ActorSearchDTO
	searchDTO.Keyword = c.GetString("keyword", "")

	var voList []actorVO.ActorItemVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		items := impl.actorBizService.SearchActor(ctx, searchDTO, tx)

		for _, item := range items {
			voList = append(voList, actorVO.ActorItemVO{
				Id:   item.Id,
				Name: item.Name,
			})
		}
	})

	return voList
}

func (impl *ActorFacade) UpdateActor(c *web.Controller) {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetInt64NotInvalid(c, "id")
	// 名称
	name := impl.GetStringNotEmpty(c, "name")
	// 描述
	description := c.GetString("description", "")
	// 获取封面文件
	reader, size := impl.GetFile(c, "cover")

	var lastCoverUrl string
	db.DoTransaction(func(tx orm.TxOrmer) {
		dto := actorDTO.ActorUpdateDTO{
			Id:          id,
			Name:        name,
			Description: description,
		}

		coverDTO := fileDTO.FileDTO{
			File: reader,
			Size: size,
		}

		origin := impl.actorBizService.UpdateActor(ctx, dto.Id, dto, coverDTO, tx)
		if coverDTO.File != nil {
			lastCoverUrl = origin.CoverUrl
		}
	})

	go func() {
		if lastCoverUrl != "" {
			impl.actorBizService.RemoveLastCover(ctx, lastCoverUrl)
		}
	}()
}

func (impl *ActorFacade) DeleteActor(c *web.Controller) {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetRestfulParamInt64(c, ":id")

	var lastCoverUrl string
	db.DoTransaction(func(tx orm.TxOrmer) {
		origin := impl.actorBizService.DeleteActor(ctx, id, tx)
		lastCoverUrl = origin.CoverUrl

		impl.performBizService.DeleteActor(ctx, id, tx)
	})

	go func() {
		if lastCoverUrl != "" {
			impl.actorBizService.RemoveLastCover(ctx, lastCoverUrl)
		}
	}()
}
