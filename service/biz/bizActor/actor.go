package bizActor

import (
	"fmt"
	"github.com/Mellolo/common/errors"
	"github.com/beego/beego/v2/client/orm"
	"media-station/generator"
	"media-station/models/do/actorDO"
	"media-station/models/dto/actorDTO"
	"media-station/models/dto/contextDTO"
	"media-station/models/dto/fileDTO"
	"media-station/storage/db"
	"media-station/storage/oss"
)

const (
	actorCoverIdGenerateKey = "actor"
	bucketActor             = "actor"
)

type ActorBizService interface {
	GetActor(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) actorDTO.ActorDTO
	GetActorCover(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) actorDTO.ActorCoverFileDTO
	CreateActor(ctx contextDTO.ContextDTO, createDTO actorDTO.ActorCreateDTO, coverDTO fileDTO.FileDTO, tx ...orm.TxOrmer) int64
	UpdateActor(ctx contextDTO.ContextDTO, id int64, updateDTO actorDTO.ActorUpdateDTO, coverDTO fileDTO.FileDTO, tx ...orm.TxOrmer) actorDTO.ActorDTO
	RemoveLastCover(ctx contextDTO.ContextDTO, lastCoverUrl string)
	DeleteActor(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) actorDTO.ActorDTO
	SearchActor(ctx contextDTO.ContextDTO, dto actorDTO.ActorSearchDTO, tx ...orm.TxOrmer) []actorDTO.ActorDTO
}

type ActorBizServiceImpl struct {
	actorMapper    db.ActorMapper
	idGenerator    generator.IdGenerator
	pictureStorage oss.PictureStorage
}

func NewActorBizService() *ActorBizServiceImpl {
	return &ActorBizServiceImpl{
		actorMapper:    db.NewActorMapper(),
		idGenerator:    generator.NewIdGenerator(),
		pictureStorage: oss.NewPictureStorage(),
	}
}

func (impl *ActorBizServiceImpl) GetActor(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) actorDTO.ActorDTO {
	actor, err := impl.actorMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("get actor [%d] failed", id)))
	}

	return impl.convertActorDO2ActorDTO(actor)
}

func (impl *ActorBizServiceImpl) GetActorCover(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) actorDTO.ActorCoverFileDTO {
	actor, err := impl.actorMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("get actor [%d] failed", id)))
	}

	pic := impl.pictureStorage.Download(bucketActor, actor.CoverUrl)
	return actorDTO.ActorCoverFileDTO{
		Reader: pic.Reader,
		Header: pic.Header,
	}
}

func (impl *ActorBizServiceImpl) CreateActor(ctx contextDTO.ContextDTO, createDTO actorDTO.ActorCreateDTO, coverDTO fileDTO.FileDTO, tx ...orm.TxOrmer) int64 {
	actor := actorDO.ActorDO{
		Name:        createDTO.Name,
		Description: createDTO.Description,
		Creator:     createDTO.Creator,
	}
	// 上传封面
	if coverDTO.File != nil {
		filename := impl.idGenerator.GenerateId(actorCoverIdGenerateKey)
		path := fmt.Sprintf("%s.jpg", filename)
		impl.pictureStorage.Upload(bucketActor, path, coverDTO.File, coverDTO.Size)
		actor.CoverUrl = path
	}

	id, err := impl.actorMapper.Insert(actor, tx...)
	if err != nil {
		panic(errors.WrapError(err, "create actor failed"))
	}
	return id
}

func (impl *ActorBizServiceImpl) UpdateActor(ctx contextDTO.ContextDTO, id int64, updateDTO actorDTO.ActorUpdateDTO, coverDTO fileDTO.FileDTO, tx ...orm.TxOrmer) actorDTO.ActorDTO {
	actor, err := impl.actorMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("actor [%d] doesn't exist", id)))
	}
	origin := *actor.DeepCopy()

	if updateDTO.Name != "" {
		actor.Name = updateDTO.Name
	}
	if updateDTO.Description != "" {
		actor.Description = updateDTO.Description
	}

	// 上传封面
	if coverDTO.File != nil {
		filename := impl.idGenerator.GenerateId(actorCoverIdGenerateKey)
		path := fmt.Sprintf("%s.jpg", filename)
		impl.pictureStorage.Upload(bucketActor, path, coverDTO.File, coverDTO.Size)
		actor.CoverUrl = path
	}

	err = impl.actorMapper.Update(id, actor)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("update actor [%d] failed", id)))
	}
	return impl.convertActorDO2ActorDTO(origin)
}

func (impl *ActorBizServiceImpl) RemoveLastCover(ctx contextDTO.ContextDTO, lastCoverUrl string) {
	impl.pictureStorage.Remove(bucketActor, lastCoverUrl)
}

func (impl *ActorBizServiceImpl) DeleteActor(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) actorDTO.ActorDTO {
	actor, err := impl.actorMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("actor [%d] doesn't exist", id)))
	}

	err = impl.actorMapper.DeleteById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("delete actor [%d] failed", id)))
	}

	return impl.convertActorDO2ActorDTO(actor)
}

func (impl *ActorBizServiceImpl) SearchActor(ctx contextDTO.ContextDTO, searchDTO actorDTO.ActorSearchDTO, tx ...orm.TxOrmer) []actorDTO.ActorDTO {
	// 读取数据库
	var actorDOList []actorDO.ActorDO
	if searchDTO.Keyword == "" {
		doList, err := impl.actorMapper.SelectAllLimit(200, tx...)
		if err != nil {
			panic(errors.WrapError(err, "select all Actor error"))
		}
		actorDOList = append(actorDOList, doList...)
	} else {
		doList, err := impl.actorMapper.SelectByKeyword(searchDTO.Keyword)
		if err != nil {
			panic(errors.WrapError(err, fmt.Sprintf("select Actor by keyword [%s] error", searchDTO.Keyword)))
		}
		actorDOList = append(actorDOList, doList...)
	}

	var actorItems []actorDTO.ActorDTO
	for _, actor := range actorDOList {
		actorItems = append(actorItems, impl.convertActorDO2ActorDTO(actor))
	}
	return actorItems
}

func (impl *ActorBizServiceImpl) convertActorDO2ActorDTO(do actorDO.ActorDO) actorDTO.ActorDTO {
	return actorDTO.ActorDTO{
		Id:          do.Id,
		Name:        do.Name,
		Description: do.Description,
		Creator:     do.Creator,
		CoverUrl:    do.CoverUrl,
	}
}
