package facade

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/mellolo/common/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"media-station/enum"
	"media-station/models/dto/fileDTO"
	"media-station/models/dto/galleryDTO"
	"media-station/models/dto/performDTO"
	"media-station/models/dto/tagDTO"
	"media-station/models/vo/galleryVO"
	"media-station/service/biz/bizActor"
	"media-station/service/biz/bizGallery"
	"media-station/service/biz/bizPerform"
	"media-station/service/biz/bizTag"
	"media-station/storage/db"
)

type GalleryFacade struct {
	AbstractFacade
	galleryBizService bizGallery.GalleryBizService
	actorBizService   bizActor.ActorBizService
	tagBizService     bizTag.TagBizService
	performBizService bizPerform.PerformBizService
}

func NewGalleryFacade() *GalleryFacade {
	return &GalleryFacade{
		galleryBizService: bizGallery.NewGalleryBizService(),
		actorBizService:   bizActor.NewActorBizService(),
		tagBizService:     bizTag.NewTagBizService(),
		performBizService: bizPerform.NewPerformBizService(),
	}
}

func (impl *GalleryFacade) SearchGallery(c *web.Controller) []galleryVO.GalleryItemVO {
	// 上下文
	ctx := impl.GetContext(c)
	// 关键词
	keyword := c.GetString("keyword", "")
	// 演员
	actors := impl.GetStringAsInt64List(c, "actors")
	actors = sets.NewInt64(actors...).List()
	// tags
	tags := impl.GetStringAsStringList(c, "tags")
	tags = sets.NewString(tags...).List()

	var voList []galleryVO.GalleryItemVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		// 仅有关键词
		if len(actors) == 0 && len(tags) == 0 {
			items := impl.galleryBizService.SearchGalleryByKeyword(ctx, keyword, tx)
			for _, item := range items {
				voList = append(voList, galleryVO.GalleryItemVO{
					Id:              item.Id,
					Name:            item.Name,
					PageCount:       len(item.PicPaths),
					PermissionLevel: item.PermissionLevel,
				})
			}
			return
		}

		searchDTO := galleryDTO.GallerySearchDTO{
			Keyword: keyword,
		}

		if len(actors) > 0 && len(tags) > 0 {
			galleryIdSet := sets.NewInt64(impl.performBizService.SelectArtByActor(ctx, enum.ArtGallery, actors, tx)...)
			galleryIdSet = galleryIdSet.Intersection(sets.NewInt64(impl.tagBizService.SelectArtByTag(ctx, enum.ArtGallery, tags, tx)...))
			searchDTO.Ids = galleryIdSet.List()
		} else if len(actors) > 0 {
			searchDTO.Ids = impl.performBizService.SelectArtByActor(ctx, enum.ArtGallery, actors, tx)
		} else {
			searchDTO.Ids = impl.tagBizService.SelectArtByTag(ctx, enum.ArtGallery, tags, tx)
		}

		items := impl.galleryBizService.SearchGallery(ctx, searchDTO, tx)

		for _, item := range items {
			voList = append(voList, galleryVO.GalleryItemVO{
				Id:              item.Id,
				Name:            item.Name,
				PageCount:       len(item.PicPaths),
				PermissionLevel: item.PermissionLevel,
			})
		}
	})

	return voList
}

func (impl *GalleryFacade) SearchGalleryByTag(c *web.Controller) []galleryVO.GalleryItemVO {
	// 上下文
	ctx := impl.GetContext(c)
	// tag
	tagName := c.GetString("tag", "")

	var voList []galleryVO.GalleryItemVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		galleryIds := impl.tagBizService.SelectArtByTag(ctx, enum.ArtGallery, []string{tagName}, tx)
		if len(galleryIds) == 0 {
			return
		}

		items := impl.galleryBizService.SearchGallery(ctx, galleryDTO.GallerySearchDTO{
			Ids: galleryIds,
		}, tx)

		for _, item := range items {
			voList = append(voList, galleryVO.GalleryItemVO{
				Id:              item.Id,
				Name:            item.Name,
				PageCount:       len(item.PicPaths),
				PermissionLevel: item.PermissionLevel,
			})
		}
	})

	return voList
}

func (impl *GalleryFacade) GetGalleryPage(c *web.Controller) galleryVO.GalleryPageVO {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetRestfulParamInt64(c, ":id")

	var vo galleryVO.GalleryPageVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		gallery := impl.galleryBizService.GetGallery(ctx, id, tx)

		actorIds := impl.performBizService.SelectActorByArt(ctx, enum.ArtGallery, id, tx)
		var actors []galleryVO.GalleryActorVO
		for _, actorId := range actorIds {
			actor := impl.actorBizService.GetActor(ctx, actorId)
			actors = append(actors, galleryVO.GalleryActorVO{
				Id:   actor.Id,
				Name: actor.Name,
			})
		}

		tags := impl.tagBizService.SelectTagByArt(ctx, enum.ArtGallery, id, tx)

		vo = galleryVO.GalleryPageVO{
			Id:              gallery.Id,
			Name:            gallery.Name,
			Description:     gallery.Description,
			PageCount:       len(gallery.PicPaths),
			Actors:          actors,
			Tags:            tags,
			PermissionLevel: gallery.PermissionLevel,
		}
	})

	return vo
}

func (impl *GalleryFacade) GetGalleryCover(c *web.Controller) galleryVO.GalleryFileVO {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetRestfulParamInt64(c, ":id")

	dto := impl.galleryBizService.ShowGalleryPage(ctx, id, 1)

	return galleryVO.GalleryFileVO{
		Header: dto.Header,
		Reader: dto.Reader,
	}
}

func (impl *GalleryFacade) UploadGallery(c *web.Controller) {
	// 上下文
	ctx := impl.GetContext(c)
	// 名称
	name := impl.GetStringNotEmpty(c, "name")
	// 描述
	description := c.GetString("description", "")
	// 演员
	actors := impl.GetStringAsInt64List(c, "actorIds")
	actors = sets.NewInt64(actors...).List()
	// tags
	tags := impl.GetStringAsStringList(c, "tags")
	tags = sets.NewString(tags...).List()
	// 上传者
	uploader := ctx.UserClaim.Username
	// 权限
	permissionLevel := c.GetString("permissionLevel", "")

	// 图片文件
	headers, err := c.GetFiles("files")
	if err != nil {
		panic(errors.WrapError(err, "get gallery files failed"))
	}

	db.DoTransaction(func(tx orm.TxOrmer) {
		createDTO := galleryDTO.GalleryCreateDTO{
			Name:            name,
			Description:     description,
			Uploader:        uploader,
			PermissionLevel: permissionLevel,
		}

		var galleryFileDTOList []fileDTO.FileDTO
		for _, header := range headers {
			file, fileErr := header.Open()
			if fileErr != nil {
				panic(errors.WrapError(fileErr, "open gallery file failed"))
			}
			galleryFileDTOList = append(galleryFileDTOList, fileDTO.FileDTO{
				File: file,
				Size: header.Size,
			})
		}

		if len(galleryFileDTOList) == 0 {
			panic(errors.NewError("no images in this upload"))
		}

		// 创建图集
		id := impl.galleryBizService.CreateGallery(ctx, createDTO, galleryFileDTOList)
		// 更新作品actor出演关系
		impl.performBizService.InsertOrUpdateActorsOfArt(ctx, performDTO.ArtPerformDTO{
			ArtType:  enum.ArtGallery,
			ArtId:    id,
			ActorIds: actors,
		}, tx)
		// 更新作品tag
		impl.tagBizService.InsertOrUpdateTagsOfArt(ctx, tagDTO.ArtTagDTO{
			ArtType: enum.ArtGallery,
			ArtId:   id,
			Tags:    tags,
		}, tx)
	})
}

func (impl *GalleryFacade) UpdateGallery(c *web.Controller) {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetInt64NotInvalid(c, "id")
	// 名称
	name := impl.GetStringNotEmpty(c, "name")
	// 描述
	description := c.GetString("description", "")
	// 演员
	actors := impl.GetStringAsInt64List(c, "actorIds")
	actors = sets.NewInt64(actors...).List()
	// tags
	tags := impl.GetStringAsStringList(c, "tags")
	tags = sets.NewString(tags...).List()
	// 上传页码
	var pages []galleryDTO.GalleryUpdatePageDTO
	impl.GetStringAsStruct(c, "pages", &pages)
	if len(pages) == 0 {
		panic(errors.NewError("no images in this update"))
	}
	// 权限
	permissionLevel := c.GetString("permissionLevel", "")

	// 图片文件
	headers, err := c.GetFiles("files")
	if err != nil {
		panic(errors.WrapError(err, "get gallery files failed"))
	}

	var dir string
	var filenames []string
	db.DoTransaction(func(tx orm.TxOrmer) {
		createDTO := galleryDTO.GalleryUpdateDTO{
			Id:              id,
			Name:            name,
			Description:     description,
			Pages:           pages,
			PermissionLevel: permissionLevel,
		}

		var galleryFileDTOList []fileDTO.FileDTO
		for _, header := range headers {
			file, fileErr := header.Open()
			if fileErr != nil {
				panic(errors.WrapError(fileErr, "open gallery file failed"))
			}
			galleryFileDTOList = append(galleryFileDTOList, fileDTO.FileDTO{
				File: file,
				Size: header.Size,
			})
		}

		// 创建图集
		origin := impl.galleryBizService.UpdateGallery(ctx, createDTO, galleryFileDTOList)
		dir = origin.DirPath
		filenames = origin.PicPaths
		// 更新作品actor出演关系
		impl.performBizService.InsertOrUpdateActorsOfArt(ctx, performDTO.ArtPerformDTO{
			ArtType:  enum.ArtGallery,
			ArtId:    id,
			ActorIds: actors,
		}, tx)
		// 更新作品tag
		impl.tagBizService.InsertOrUpdateTagsOfArt(ctx, tagDTO.ArtTagDTO{
			ArtType: enum.ArtGallery,
			ArtId:   id,
			Tags:    tags,
		}, tx)
	})

	go func() {
		if dir != "" {
			impl.galleryBizService.RemoveGalleryDir(ctx, dir, filenames)
		}
	}()
}

func (impl *GalleryFacade) DeleteGallery(c *web.Controller) {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetRestfulParamInt64(c, ":id")

	var dir string
	var filenames []string
	db.DoTransaction(func(tx orm.TxOrmer) {
		// 更新视频
		gallery := impl.galleryBizService.DeleteGallery(ctx, id, tx)
		dir = gallery.DirPath
		filenames = gallery.PicPaths
		// 更新作品actor出演关系
		impl.performBizService.DeleteArt(ctx, enum.ArtGallery, id, tx)
		// 更新作品tag
		impl.tagBizService.DeleteArt(ctx, enum.ArtGallery, id, tx)
	})

	go func() {
		if dir != "" {
			impl.galleryBizService.RemoveGalleryDir(ctx, dir, filenames)
		}
	}()
}

func (impl *GalleryFacade) ShowPage(c *web.Controller) galleryVO.GalleryFileVO {
	// 上下文
	ctx := impl.GetContext(c)
	// id
	id := impl.GetRestfulParamInt64(c, ":id")
	// 页数
	page := impl.GetRestfulParamInt(c, ":page")

	dto := impl.galleryBizService.ShowGalleryPage(ctx, id, page)

	return galleryVO.GalleryFileVO{
		Header: dto.Header,
		Reader: dto.Reader,
	}
}
