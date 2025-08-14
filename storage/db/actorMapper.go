package db

import (
	"github.com/beego/beego/v2/client/orm"
	"media-station/models/dao/actorDAO"
	"media-station/models/dao/daoCommon"
	"media-station/models/do/actorDO"
)

type ActorMapper interface {
	Insert(actor actorDO.ActorDO, tx ...orm.TxOrmer) (int64, error)
	SelectById(id int64, tx ...orm.TxOrmer) (actorDO.ActorDO, error)
	SelectAllLimit(limit int, tx ...orm.TxOrmer) ([]actorDO.ActorDO, error)
	SelectByKeyword(keyword string, tx ...orm.TxOrmer) ([]actorDO.ActorDO, error)
	Update(id int64, actor actorDO.ActorDO, tx ...orm.TxOrmer) error
	DeleteById(id int64, tx ...orm.TxOrmer) error
}

type ActorMapperImpl struct{}

func NewActorMapper() *ActorMapperImpl {
	return &ActorMapperImpl{}
}

func (impl *ActorMapperImpl) Insert(actor actorDO.ActorDO, tx ...orm.TxOrmer) (int64, error) {
	executor := getQueryExecutor(tx...)

	// todo 可能会扩展更复杂的处理
	record := actorDAO.ActorRecord{
		Name:        actor.Name,
		Description: actor.Description,
		Creator:     actor.Creator,
		CoverUrl:    actor.CoverUrl,
		Details:     "{}",
	}

	id, err := executor.Insert(&record)
	return id, err
}

func (impl *ActorMapperImpl) SelectById(id int64, tx ...orm.TxOrmer) (actorDO.ActorDO, error) {
	executor := getQueryExecutor(tx...)
	record := actorDAO.ActorRecord{CommonColumn: daoCommon.CommonColumn{Id: id}}
	err := executor.Read(&record)
	if err != nil {
		return actorDO.ActorDO{}, err
	}

	do := actorDO.ActorDO{
		Id:          record.Id,
		CreateAt:    record.CreatedAt.String(),
		Name:        record.Name,
		Description: record.Description,
		Creator:     record.Creator,
		CoverUrl:    record.CoverUrl,
	}
	return do, nil
}

func (impl *ActorMapperImpl) SelectAllLimit(limit int, tx ...orm.TxOrmer) ([]actorDO.ActorDO, error) {
	executor := getQueryExecutor(tx...)

	var records []actorDAO.ActorRecord

	_, err := executor.QueryTable(actorDAO.TableActor).Limit(limit).All(&records)
	if err != nil {
		return nil, err
	}

	var actors []actorDO.ActorDO
	for _, record := range records {
		do := actorDO.ActorDO{
			Id:          record.Id,
			CreateAt:    record.CreatedAt.String(),
			Name:        record.Name,
			Description: record.Description,
			Creator:     record.Creator,
			CoverUrl:    record.CoverUrl,
		}
		actors = append(actors, do)
	}

	return actors, nil
}

func (impl *ActorMapperImpl) SelectByKeyword(keyword string, tx ...orm.TxOrmer) ([]actorDO.ActorDO, error) {
	executor := getQueryExecutor(tx...)

	var records []actorDAO.ActorRecord
	_, err := executor.QueryTable(actorDAO.TableActor).
		Filter("name__icontains", keyword).
		All(&records)
	if err != nil {
		return nil, err
	}

	var actors []actorDO.ActorDO
	for _, record := range records {
		do := actorDO.ActorDO{
			Id:          record.Id,
			CreateAt:    record.CreatedAt.String(),
			Name:        record.Name,
			Description: record.Description,
			Creator:     record.Creator,
			CoverUrl:    record.CoverUrl,
		}
		actors = append(actors, do)
	}

	return actors, nil
}

func (impl *ActorMapperImpl) Update(id int64, actor actorDO.ActorDO, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)

	record := actorDAO.ActorRecord{
		CommonColumn: daoCommon.CommonColumn{
			Id: id,
		},
		Name:        actor.Name,
		Description: actor.Description,
		Creator:     actor.Creator,
		CoverUrl:    actor.CoverUrl,
		Details:     "{}",
	}
	_, err := executor.Update(&record)
	return err
}

func (impl *ActorMapperImpl) DeleteById(id int64, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)
	record := &actorDAO.ActorRecord{CommonColumn: daoCommon.CommonColumn{Id: id}}
	_, err := executor.Delete(record)
	return err
}
