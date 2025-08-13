package facade

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/mellolo/common/errors"
	"media-station/models/dto/actorDTO"
	"media-station/models/dto/fileDTO"
	"media-station/models/vo/actorVO"
	"media-station/models/vo/videoVO"
	"media-station/service/biz/bizActor"
	"media-station/service/biz/bizGallery"
	"media-station/service/biz/bizVideo"
	"media-station/storage/db"
)

type ActorFacade struct {
	AbstractFacade
	actorBizService   bizActor.ActorBizService
	videoBizService   bizVideo.VideoBizService
	galleryBizService bizGallery.GalleryBizService
}

func NewActorFacade() *ActorFacade {
	return &ActorFacade{
		actorBizService:   bizActor.NewActorBizService(),
		videoBizService:   bizVideo.NewVideoBizService(),
		galleryBizService: bizGallery.NewGalleryBizService(),
	}
}

func (impl *ActorFacade) GetActorPage(c *web.Controller) actorVO.ActorPageVO {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetRestfulParamInt64(c, ":id")

	var vo actorVO.ActorPageVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		actorPage := impl.actorBizService.GetActorPage(ctx, id)
		vo = actorVO.ActorPageVO{
			Id:          actorPage.Id,
			Name:        actorPage.Name,
			Description: actorPage.Description,
			Creator:     actorPage.Creator,
			GalleryIds:  actorPage.Art.GalleryIds,
		}

		items := impl.videoBizService.SearchVideoByActor(ctx, id, tx)

		for _, item := range items {
			vo.Videos = append(vo.Videos, videoVO.VideoItemVO{
				Id:              item.Id,
				Name:            item.Name,
				Duration:        item.Duration,
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

	createDTO := actorDTO.ActorCreateDTO{
		Name:        name,
		Description: description,
		Creator:     creator,
	}

	// 获取封面文件
	reader, size := impl.GetFileNotInvalid(c, "cover")

	var id int64
	db.DoTransaction(func(tx orm.TxOrmer) {
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

		lastCoverUrl = impl.actorBizService.UpdateActor(ctx, dto.Id, dto, coverDTO, tx)
	})

	if lastCoverUrl != "" {
		impl.actorBizService.RemoveLastCover(ctx, lastCoverUrl)
	}

}

func (impl *ActorFacade) DeleteActor(c *web.Controller) {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetRestfulParamInt64(c, ":id")

	var lastCoverUrl string
	db.DoTransaction(func(tx orm.TxOrmer) {
		page := impl.actorBizService.DeleteActor(ctx, id, tx)
		lastCoverUrl = page.CoverUrl

		for _, videoId := range page.Art.VideoIds {
			impl.videoBizService.RemoveActor(ctx, videoId, []int64{id}, tx)
		}

		for _, galleryId := range page.Art.GalleryIds {
			impl.galleryBizService.RemoveActor(ctx, galleryId, []int64{id}, tx)
		}
	})

	impl.actorBizService.RemoveLastCover(ctx, lastCoverUrl)

}
