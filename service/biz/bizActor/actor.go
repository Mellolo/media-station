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
	GetActorPage(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) actorDTO.ActorPageDTO
	GetActorCover(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) actorDTO.ActorCoverFileDTO
	CreateActor(ctx contextDTO.ContextDTO, createDTO actorDTO.ActorCreateDTO, coverDTO fileDTO.FileDTO, tx ...orm.TxOrmer) int64
	UpdateActor(ctx contextDTO.ContextDTO, id int64, updateDTO actorDTO.ActorUpdateDTO, coverDTO fileDTO.FileDTO, tx ...orm.TxOrmer) string
	RemoveLastCover(ctx contextDTO.ContextDTO, lastCoverUrl string)
	RemoveArt(ctx contextDTO.ContextDTO, dto actorDTO.ActorRemoveArtDTO, tx ...orm.TxOrmer)
	DeleteActor(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) string
	SearchActor(ctx contextDTO.ContextDTO, dto actorDTO.ActorSearchDTO, tx ...orm.TxOrmer) []actorDTO.ActorItemDTO
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

func (impl *ActorBizServiceImpl) GetActorPage(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) actorDTO.ActorPageDTO {
	actor, err := impl.actorMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("get actor [%d] failed", id)))
	}

	return actorDTO.ActorPageDTO{
		Id:          actor.Id,
		Name:        actor.Name,
		Description: actor.Description,
		Creator:     actor.Creator,
		CoverUrl:    actor.CoverUrl,
		Art: actorDTO.ActorArtDTO{
			VideoIds:   actor.Art.VideoIds,
			GalleryIds: actor.Art.GalleryIds,
		},
	}
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

func (impl *ActorBizServiceImpl) UpdateActor(ctx contextDTO.ContextDTO, id int64, updateDTO actorDTO.ActorUpdateDTO, coverDTO fileDTO.FileDTO, tx ...orm.TxOrmer) string {
	actor, err := impl.actorMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("actor [%d] doesn't exist", id)))
	}

	if updateDTO.Name != "" {
		actor.Name = updateDTO.Name
	}
	if updateDTO.Description != "" {
		actor.Description = updateDTO.Description
	}

	videoIds := sets.NewInt64(actor.Art.VideoIds...)
	videoIds.Insert(updateDTO.Art.VideoIds...)
	actor.Art.VideoIds = videoIds.List()

	galleryIds := sets.NewInt64(actor.Art.GalleryIds...)
	galleryIds.Insert(updateDTO.Art.GalleryIds...)
	actor.Art.GalleryIds = galleryIds.List()

	// 上传封面
	lastCoverUrl := ""
	if coverDTO.File != nil {
		filename := impl.idGenerator.GenerateId(actorCoverIdGenerateKey)
		path := fmt.Sprintf("%s.jpg", filename)
		impl.pictureStorage.Upload(bucketActor, path, coverDTO.File, coverDTO.Size)
		lastCoverUrl = actor.CoverUrl
		actor.CoverUrl = path
	}

	err = impl.actorMapper.Update(id, actor)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("update actor [%d] failed", id)))
	}
	return lastCoverUrl
}

func (impl *ActorBizServiceImpl) RemoveLastCover(ctx contextDTO.ContextDTO, lastCoverUrl string) {
	impl.pictureStorage.Remove(bucketActor, lastCoverUrl)
}

func (impl *ActorBizServiceImpl) RemoveArt(ctx contextDTO.ContextDTO, dto actorDTO.ActorRemoveArtDTO, tx ...orm.TxOrmer) {
	actor, err := impl.actorMapper.SelectById(dto.Id, tx...)
	if err != nil && !pkgErrors.Is(err, orm.ErrNoRows) {
		panic(errors.WrapError(err, "select actor failed"))
	}
	actor.Art.VideoIds = sets.NewInt64(actor.Art.VideoIds...).Delete(dto.Art.VideoIds...).List()
	actor.Art.GalleryIds = sets.NewInt64(actor.Art.GalleryIds...).Delete(dto.Art.GalleryIds...).List()
	err = impl.actorMapper.Update(dto.Id, actor, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("remove artwork for actor [%s] failed", dto.Id)))
	}
}

func (impl *ActorBizServiceImpl) DeleteActor(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) string {
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

func (impl *ActorBizServiceImpl) SearchActor(ctx contextDTO.ContextDTO, searchDTO actorDTO.ActorSearchDTO, tx ...orm.TxOrmer) []actorDTO.ActorItemDTO {
	// 读取数据库
	var actorDOList []*actorDO.ActorDO
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

	var actorItems []actorDTO.ActorItemDTO
	for _, actor := range actorDOList {
		actorItems = append(actorItems, actorDTO.ActorItemDTO{
			Id:          actor.Id,
			Name:        actor.Name,
			Description: actor.Description,
			CoverUrl:    actor.CoverUrl,
		})
	}
	return actorItems
}
