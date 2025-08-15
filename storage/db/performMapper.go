package db

import (
	"github.com/beego/beego/v2/client/orm"
	"media-station/models/dao/performDAO"
	"media-station/models/do/performDO"
)

type PerformMapper interface {
	InsertOrUpdateActorsOfArt(artType string, artId int64, items []performDO.PerformDO, tx ...orm.TxOrmer) error
	SelectArtByActor(artType string, actorId int64, tx ...orm.TxOrmer) ([]performDO.PerformDO, error)
	SelectActorByArt(artType string, artId int64, tx ...orm.TxOrmer) ([]performDO.PerformDO, error)
	DeleteArt(artType string, artId int64, tx ...orm.TxOrmer) error
	DeleteActor(actorId int64, tx ...orm.TxOrmer) error
}

func NewPerformMapper() *PerformMapperImpl {
	return &PerformMapperImpl{}
}

type PerformMapperImpl struct{}

func (impl PerformMapperImpl) InsertOrUpdateActorsOfArt(artType string, artId int64, items []performDO.PerformDO, tx ...orm.TxOrmer) error {
	err := impl.DeleteArt(artType, artId, tx...)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	executor := getQueryExecutor(tx...)

	var records []performDAO.PerformRecord
	for _, item := range items {
		records = append(records, performDAO.PerformRecord{
			ArtType: artType,
			ArtId:   artId,
			ActorId: item.ActorId,
		})
	}

	_, err = executor.InsertMulti(100, records)
	return err
}

func (impl PerformMapperImpl) SelectArtByActor(artType string, actorId int64, tx ...orm.TxOrmer) ([]performDO.PerformDO, error) {
	executor := getQueryExecutor(tx...)

	var records []performDAO.PerformRecord
	_, err := executor.QueryTable(performDAO.TablePerform).
		Filter("art_type", artType).
		Filter("actor_id", actorId).
		All(&records, "art_type", "art_id", "actor_id")
	if err != nil {
		return nil, err
	}

	var items []performDO.PerformDO
	for _, record := range records {
		do := performDO.PerformDO{
			ArtType: record.ArtType,
			ArtId:   record.ArtId,
			ActorId: record.ActorId,
		}
		items = append(items, do)
	}

	return items, nil
}

func (impl PerformMapperImpl) SelectActorByArt(artType string, artId int64, tx ...orm.TxOrmer) ([]performDO.PerformDO, error) {
	executor := getQueryExecutor(tx...)

	var records []performDAO.PerformRecord
	_, err := executor.QueryTable(performDAO.TablePerform).
		Filter("art_type", artType).
		Filter("art_id", artId).
		All(&records, "art_type", "actor_id", "art_id")
	if err != nil {
		return nil, err
	}

	var items []performDO.PerformDO
	for _, record := range records {
		do := performDO.PerformDO{
			ArtType: record.ArtType,
			ArtId:   record.ArtId,
			ActorId: record.ActorId,
		}
		items = append(items, do)
	}

	return items, nil
}

func (impl PerformMapperImpl) DeleteArt(artType string, artId int64, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)

	_, err := executor.QueryTable(performDAO.TablePerform).
		Filter("art_type", artType).
		Filter("art_id", artId).
		Delete()
	return err
}

func (impl PerformMapperImpl) DeleteActor(actorId int64, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)

	_, err := executor.QueryTable(performDAO.TablePerform).
		Filter("actor_id", actorId).
		Delete()
	return err
}
