package facade

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/mellolo/common/errors"
	"github.com/mellolo/common/utils/jsonUtil"
	"media-station/controllers/filters"
	"media-station/models/dto/actorDTO"
	"media-station/models/dto/fileDTO"
	"media-station/models/dto/galleryDTO"
	"media-station/models/dto/tagDTO"
	"media-station/models/dto/userDTO"
	"media-station/models/vo/galleryVO"
	"media-station/service/biz/bizActor"
	"media-station/service/biz/bizGallery"
	"media-station/service/biz/bizTag"
	"media-station/storage/db"
	"strconv"
)

type GalleryFacade struct {
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
	var dto galleryDTO.GallerySearchDTO
	jsonUtil.UnmarshalJsonString(string(c.Ctx.Input.RequestBody), &dto)

	var voList []galleryVO.GalleryItemVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		items := impl.galleryBizService.SearchGallery(dto, tx)

		for _, item := range items {
			voList = append(voList, galleryVO.GalleryItemVO{
				Id:              item.Id,
				Name:            item.Name,
				CoverUrl:        item.CoverUrl,
				PermissionLevel: item.PermissionLevel,
			})
		}
	})

	return voList
}

func (impl *GalleryFacade) GetGalleryPage(c *web.Controller) galleryVO.GalleryPageVO {
	// id
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		panic(errors.WrapError(err, "get videoId failed"))
	}

	var vo galleryVO.GalleryPageVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		page := impl.galleryBizService.GetGalleryPage(id, tx)
		vo = galleryVO.GalleryPageVO{
			Id:              page.Id,
			Name:            page.Name,
			Description:     page.Description,
			CoverUrl:        page.CoverUrl,
			PermissionLevel: page.PermissionLevel,
		}
	})

	return vo
}

func (impl *GalleryFacade) UploadGallery(c *web.Controller, ch chan string) {
	// 名称
	name := c.GetString("name", "")
	// 描述
	description := c.GetString("description", "")
	// 演员
	var actors []int64
	jsonUtil.UnmarshalJsonString(c.GetString("actors", "[]"), &actors)
	// tag
	var tags []string
	jsonUtil.UnmarshalJsonString(c.GetString("tags", "[]"), &tags)
	// 用户
	uploader := ""
	if claim, ok := c.Ctx.Input.GetData(filters.ContextClaim).(string); ok {
		var userClaim userDTO.UserClaimDTO
		jsonUtil.UnmarshalJsonString(claim, &userClaim)
		uploader = userClaim.Username
	}
	// 权限
	permissionLevel := c.GetString("permissionLevel", "")

	createDTO := galleryDTO.GalleryCreateDTO{
		Name:            name,
		Description:     description,
		Actors:          actors,
		Tags:            tags,
		Uploader:        uploader,
		PermissionLevel: permissionLevel,
	}

	// 图片文件
	headers, err := c.GetFiles("files")
	if err != nil {
		panic(errors.WrapError(err, "get gallery files failed"))
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

	db.DoTransaction(func(tx orm.TxOrmer) {
		// 创建图集
		id := impl.galleryBizService.CreateGallery(createDTO, galleryFileDTOList, ch)
		// 更新演员作品
		for _, actorId := range createDTO.Actors {
			updateDTO := actorDTO.ActorUpdateDTO{
				Id: actorId,
				Art: actorDTO.ActorArtDTO{
					GalleryIds: []int64{id},
				},
			}
			impl.actorBizService.UpdateActor(actorId, updateDTO, fileDTO.FileDTO{}, tx)
		}
		// 更新tag作品
		for _, tagName := range createDTO.Tags {
			tag := tagDTO.TagCreateOrUpdateDTO{
				Name: tagName,
				Details: tagDTO.TagDetailsDTO{
					GalleryIds: []int64{id},
				},
			}
			impl.tagBizService.CreateOrUpdateTag(tag, tx)
		}
	})
}

func (impl *GalleryFacade) ShowPage(c *web.Controller) galleryVO.GalleryFileVO {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("param [id] %s is invalid", idStr)))
	}
	pageStr := c.Ctx.Input.Param(":page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("param [page] %s is invalid", pageStr)))
	}

	dto := impl.galleryBizService.ShowGalleryPage(id, page)

	return galleryVO.GalleryFileVO{
		Header: dto.Header,
		Reader: dto.Reader,
	}
}
