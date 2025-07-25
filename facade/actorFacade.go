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
	"media-station/models/dto/userDTO"
	"media-station/service/biz/bizActor"
	"media-station/storage/db"
	"strconv"
)

type ActorFacade struct {
	actorBizService bizActor.ActorBizService
}

func NewActorFacade() *ActorFacade {
	return &ActorFacade{
		actorBizService: bizActor.NewActorBizService(),
	}
}

func (impl *ActorFacade) CreateActor(c *web.Controller) int64 {
	// 名称
	name := c.GetString("name", "")
	if name == "" {
		panic(errors.NewError("actor name is empty"))
	}
	// 描述
	description := c.GetString("description", "")

	// 创建者
	creator := ""
	if claim, ok := c.Ctx.Input.GetData(filters.ContextClaim).(string); ok {
		var userClaim userDTO.UserClaimDTO
		jsonUtil.UnmarshalJsonString(claim, &userClaim)
		creator = userClaim.Username
	}
	if creator == "" {
		panic(errors.NewError("creator is empty"))
	}

	createDTO := actorDTO.ActorCreateDTO{
		Name:        name,
		Description: description,
		Creator:     creator,
	}

	// 获取封面文件
	coverDTO := fileDTO.FileDTO{}
	reader, header, err := c.GetFile("cover")
	if err == nil {
		coverDTO = fileDTO.FileDTO{
			File: reader,
			Size: header.Size,
		}
	}

	var id int64
	db.DoTransaction(func(tx orm.TxOrmer) {
		id = impl.actorBizService.CreateActor(createDTO, coverDTO, tx)
	})

	return id
}

func (impl *ActorFacade) UpdateActor(c *web.Controller) string {
	// 名称
	id, err := c.GetInt64("id")
	if err != nil {
		panic(errors.WrapError(err, "id is invalid"))
	}
	// 名称
	name := c.GetString("name", "")
	if name == "" {
		panic(errors.NewError("actor name is empty"))
	}
	// 描述
	description := c.GetString("description", "")

	// 获取封面文件
	coverDTO := fileDTO.FileDTO{}
	reader, header, err := c.GetFile("cover")
	if err == nil {
		coverDTO = fileDTO.FileDTO{
			File: reader,
			Size: header.Size,
		}
	}

	dto := actorDTO.ActorUpdateDTO{
		Id:          id,
		Name:        name,
		Description: description,
	}

	var lastCoverUrl string
	db.DoTransaction(func(tx orm.TxOrmer) {
		lastCoverUrl = impl.actorBizService.UpdateActor(dto.Id, dto, coverDTO, tx)
	})

	return lastCoverUrl
}

func (impl *ActorFacade) DeleteActor(c *web.Controller) string {
	// 获取演员ID
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("param [id] %s is invalid", idStr)))
	}

	var lastCoverUrl string
	db.DoTransaction(func(tx orm.TxOrmer) {
		lastCoverUrl = impl.actorBizService.DeleteActor(id, tx)
	})

	return lastCoverUrl
}
