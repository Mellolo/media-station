package bizPerform

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/mellolo/common/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"media-station/models/do/performDO"
	"media-station/models/dto/contextDTO"
	"media-station/models/dto/performDTO"
	"media-station/storage/db"
)

type PerformBizService interface {
	SelectArtByActor(ctx contextDTO.ContextDTO, artType string, actorIds []int64, tx ...orm.TxOrmer) []int64
	SelectActorByArt(ctx contextDTO.ContextDTO, artType string, artId int64, tx ...orm.TxOrmer) []int64
	InsertOrUpdateActorsOfArt(ctx contextDTO.ContextDTO, dto performDTO.ArtPerformDTO, tx ...orm.TxOrmer)
	DeleteArt(ctx contextDTO.ContextDTO, artType string, artId int64, tx ...orm.TxOrmer)
	DeleteActor(ctx contextDTO.ContextDTO, actorId int64, tx ...orm.TxOrmer)
}

func NewPerformBizService() *PerformBizServiceImpl {
	return &PerformBizServiceImpl{
		performMapper: db.NewPerformMapper(),
	}
}

type PerformBizServiceImpl struct {
	performMapper db.PerformMapper
}

func (impl PerformBizServiceImpl) SelectArtByActor(ctx contextDTO.ContextDTO, artType string, actorIds []int64, tx ...orm.TxOrmer) []int64 {
	if len(actorIds) == 0 {
		return nil
	}

	firstLoop := true
	artIdSet := sets.NewInt64()
	for _, actorId := range actorIds {
		performs, err := impl.performMapper.SelectArtByActor(artType, actorId, tx...)
		if err != nil {
			panic(errors.WrapError(err, fmt.Sprintf("get %s of actor [%d] failed", artType, actorId)))
		}
		if firstLoop {
			for _, perform := range performs {
				artIdSet.Insert(perform.ArtId)
			}
		} else {
			thisArtIdSet := sets.NewInt64()
			for _, perform := range performs {
				thisArtIdSet.Insert(perform.ArtId)
			}
			artIdSet = artIdSet.Intersection(thisArtIdSet)
		}
		firstLoop = false
	}

	return artIdSet.List()
}

func (impl PerformBizServiceImpl) SelectActorByArt(ctx contextDTO.ContextDTO, artType string, artId int64, tx ...orm.TxOrmer) []int64 {
	performs, err := impl.performMapper.SelectActorByArt(artType, artId, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("get actors of %s[%d] failed", artType, artId)))
	}
	var actorId []int64
	for _, perform := range performs {
		actorId = append(actorId, perform.ActorId)
	}

	return actorId
}

func (impl PerformBizServiceImpl) InsertOrUpdateActorsOfArt(ctx contextDTO.ContextDTO, dto performDTO.ArtPerformDTO, tx ...orm.TxOrmer) {
	var doList []performDO.PerformDO
	for _, actorId := range dto.ActorIds {
		doList = append(doList, performDO.PerformDO{
			ArtType: dto.ArtType,
			ArtId:   dto.ArtId,
			ActorId: actorId,
		})
	}

	err := impl.performMapper.InsertOrUpdateActorsOfArt(dto.ArtType, dto.ArtId, doList, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("insert actors perform art [%s][%d] failed", dto.ArtType, dto.ArtId)))
	}
}

func (impl PerformBizServiceImpl) DeleteArt(ctx contextDTO.ContextDTO, artType string, artId int64, tx ...orm.TxOrmer) {
	err := impl.performMapper.DeleteArt(artType, artId, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("delete art perform [%s][%d] failed", artType, artId)))
	}
}

func (impl PerformBizServiceImpl) DeleteActor(ctx contextDTO.ContextDTO, actorId int64, tx ...orm.TxOrmer) {
	err := impl.performMapper.DeleteActor(actorId, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("delete actor perform [%d] failed", actorId)))
	}
}
