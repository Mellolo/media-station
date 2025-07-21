package bizActor

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/mellolo/common/errors"
	pkgErrors "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"media-station/generator"
	"media-station/models/do/actorDO"
	"media-station/models/dto/actorDTO"
	"media-station/models/dto/fileDTO"
	"media-station/storage/db"
	"media-station/storage/oss"
)

const (
	actorCoverIdGenerateKey = "actor"
	bucketActor             = "actor"
)

type ActorBizService interface {
	CreateActor(createDTO actorDTO.ActorCreateDTO, coverDTO fileDTO.FileDTO, tx ...orm.TxOrmer) int64
	UpdateActor(id int64, updateDTO actorDTO.ActorUpdateDTO, coverDTO fileDTO.FileDTO, tx ...orm.TxOrmer) string
	RemoveArt(dto actorDTO.ActorRemoveArtDTO, tx ...orm.TxOrmer)
	DeleteActor(id int64, tx ...orm.TxOrmer) string
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

func (impl *ActorBizServiceImpl) CreateActor(createDTO actorDTO.ActorCreateDTO, coverDTO fileDTO.FileDTO, tx ...orm.TxOrmer) int64 {
	actor := &actorDO.ActorDO{
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

func (impl *ActorBizServiceImpl) UpdateActor(id int64, updateDTO actorDTO.ActorUpdateDTO, coverDTO fileDTO.FileDTO, tx ...orm.TxOrmer) string {
	actor, err := impl.actorMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("actor [%d] doesn't exist", id)))
	}
	lastCoverUrl := actor.CoverUrl

	if updateDTO.Name != "" {
		actor.Name = updateDTO.Name
	}
	if updateDTO.Description != "" {
		actor.Description = updateDTO.Description
	}

	videoIds := sets.NewInt64(updateDTO.Details.VideoIds...)
	videoIds.Insert(updateDTO.Details.VideoIds...)
	actor.Details.VideoIds = videoIds.List()

	galleryIds := sets.NewInt64(updateDTO.Details.GalleryIds...)
	galleryIds.Insert(updateDTO.Details.GalleryIds...)
	actor.Details.GalleryIds = galleryIds.List()

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
	return lastCoverUrl
}

func (impl *ActorBizServiceImpl) RemoveArt(dto actorDTO.ActorRemoveArtDTO, tx ...orm.TxOrmer) {
	actor, err := impl.actorMapper.SelectById(dto.Id, tx...)
	if err != nil && !pkgErrors.Is(err, orm.ErrNoRows) {
		panic(errors.WrapError(err, "select actor failed"))
	}
	actor.Details.VideoIds = sets.NewInt64(actor.Details.VideoIds...).Delete(dto.Details.VideoIds...).List()
	actor.Details.GalleryIds = sets.NewInt64(actor.Details.GalleryIds...).Delete(dto.Details.GalleryIds...).List()
	err = impl.actorMapper.Update(dto.Id, actor, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("remove artwork for actor [%s] failed", dto.Id)))
	}
}

func (impl *ActorBizServiceImpl) DeleteActor(id int64, tx ...orm.TxOrmer) string {
	actor, err := impl.actorMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("actor [%d] doesn't exist", id)))
	}
	err = impl.actorMapper.DeleteById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("delete actor [%d] failed", id)))
	}
	return actor.CoverUrl
}
