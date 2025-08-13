package facade

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/mellolo/common/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"media-station/models/dto/actorDTO"
	"media-station/models/dto/fileDTO"
	"media-station/models/dto/galleryDTO"
	"media-station/models/dto/tagDTO"
	"media-station/models/vo/galleryVO"
	"media-station/service/biz/bizActor"
	"media-station/service/biz/bizGallery"
	"media-station/service/biz/bizTag"
	"media-station/storage/db"
)

type GalleryFacade struct {
	AbstractFacade
	galleryBizService bizGallery.GalleryBizService
	actorBizService   bizActor.ActorBizService
	tagBizService     bizTag.TagBizService
}

func NewGalleryFacade() *GalleryFacade {
	return &GalleryFacade{
		galleryBizService: bizGallery.NewGalleryBizService(),
		actorBizService:   bizActor.NewActorBizService(),
		tagBizService:     bizTag.NewTagBizService(),
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
		searchDTO := galleryDTO.GallerySearchDTO{
			Keyword: keyword,
			Actors:  actors,
			Tags:    tags,
		}

		items := impl.galleryBizService.SearchGallery(ctx, searchDTO, tx)

		for _, item := range items {
			voList = append(voList, galleryVO.GalleryItemVO{
				Id:              item.Id,
				Name:            item.Name,
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
		page := impl.galleryBizService.GetGalleryPage(ctx, id, tx)
		vo = galleryVO.GalleryPageVO{
			Id:              page.Id,
			Name:            page.Name,
			Description:     page.Description,
			PermissionLevel: page.PermissionLevel,
		}
	})

	return vo
}

func (impl *GalleryFacade) UploadGallery(c *web.Controller, ch chan string) {
	// 上下文
	ctx := impl.GetContext(c)
	// 名称
	name := c.GetString("name", "")
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

	// 图片文件
	headers, err := c.GetFiles("files")
	if err != nil {
		panic(errors.WrapError(err, "get gallery files failed"))
	}

	db.DoTransaction(func(tx orm.TxOrmer) {
		createDTO := galleryDTO.GalleryCreateDTO{
			Name:            name,
			Description:     description,
			Actors:          actors,
			Tags:            tags,
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

		// 创建图集
		id := impl.galleryBizService.CreateGallery(ctx, createDTO, galleryFileDTOList, ch)
		// 更新演员作品
		for _, actorId := range createDTO.Actors {
			artDTO := actorDTO.ActorArtDTO{
				GalleryIds: []int64{id},
			}
			impl.actorBizService.AddArt(ctx, actorId, artDTO, tx)
		}
		// 更新tag作品
		for _, tagName := range createDTO.Tags {
			tag := tagDTO.TagCreateOrUpdateDTO{
				Name: tagName,
				Details: tagDTO.TagDetailsDTO{
					GalleryIds: []int64{id},
				},
			}
			impl.tagBizService.AddArtToTag(ctx, tag, tx)
		}
	})
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
