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
	"media-station/models/dto/tagDTO"
	"media-station/models/dto/userDTO"
	"media-station/models/dto/videoDTO"
	"media-station/models/vo/videoVO"
	"media-station/service/biz/bizActor"
	"media-station/service/biz/bizTag"
	"media-station/service/biz/bizVideo"
	"media-station/storage/db"
	"strconv"
)

type VideoFacade struct {
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
	var dto videoDTO.VideoSearchDTO
	jsonUtil.UnmarshalJsonString(string(c.Ctx.Input.RequestBody), &dto)

	var voList []videoVO.VideoItemVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		items := impl.videoBizService.SearchVideo(dto, tx)

		for _, item := range items {
			voList = append(voList, videoVO.VideoItemVO{
				Id:              item.Id,
				Name:            item.Name,
				PermissionLevel: item.PermissionLevel,
			})
		}
	})

	return voList
}

func (impl *VideoFacade) GetVideoPage(c *web.Controller) videoVO.VideoPageVO {
	// id
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		panic(errors.WrapError(err, "get videoId failed"))
	}

	var vo videoVO.VideoPageVO
	db.DoTransaction(func(tx orm.TxOrmer) {
		page := impl.videoBizService.GetVideoPage(id, tx)
		vo = videoVO.VideoPageVO{
			Id:              page.Id,
			Name:            page.Name,
			Description:     page.Description,
			PermissionLevel: page.PermissionLevel,
		}
	})

	return vo
}

func (impl *VideoFacade) UploadVideo(c *web.Controller, ch chan string) {
	// 名称
	name := c.GetString("name", "")
	// 描述
	description := c.GetString("description", "")
	// 演员
	var actors []int64
	jsonUtil.UnmarshalJsonString(c.GetString("actors", "[]"), &actors)
	// tags
	var tags []string
	jsonUtil.UnmarshalJsonString(c.GetString("tags", "[]"), &tags)
	// 上传者
	uploader := ""
	if claim, ok := c.Ctx.Input.GetData(filters.ContextClaim).(string); ok {
		var userClaim userDTO.UserClaimDTO
		jsonUtil.UnmarshalJsonString(claim, &userClaim)
		uploader = userClaim.Username
	}
	// 权限
	permissionLevel := c.GetString("permissionLevel", "")

	createDTO := videoDTO.VideoCreateDTO{
		Name:            name,
		Description:     description,
		Actors:          actors,
		Tags:            tags,
		Uploader:        uploader,
		PermissionLevel: permissionLevel,
	}

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
		// 创建视频
		id := impl.videoBizService.CreateVideo(createDTO, videoFileDTO, ch)
		// 更新actor作品
		for _, actorId := range createDTO.Actors {
			updateDTO := actorDTO.ActorUpdateDTO{
				Id: actorId,
				Art: actorDTO.ActorArtDTO{
					VideoIds: []int64{id},
				},
			}
			impl.actorBizService.UpdateActor(actorId, updateDTO, fileDTO.FileDTO{}, tx)
		}
		// 更新tag作品
		for _, tagName := range createDTO.Tags {
			tag := tagDTO.TagCreateOrUpdateDTO{
				Name: tagName,
				Details: tagDTO.TagDetailsDTO{
					VideoIds: []int64{id},
				},
			}
			impl.tagBizService.CreateOrUpdateTag(tag, tx)
		}
	})
}

func (impl *VideoFacade) PlayVideo(c *web.Controller) videoVO.VideoFileVO {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("param [id] %s is invalid", idStr)))
	}

	dto := impl.videoBizService.PlayVideo(id, c.Ctx.Request.Header["Range"]...)
	return videoVO.VideoFileVO{
		Header: dto.Header,
		Reader: dto.Reader,
	}
}
