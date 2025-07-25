package db

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/mellolo/common/utils/jsonUtil"
	"media-station/models/dao/tagDAO"
	"media-station/models/do/tagDO"
)

type TagMapper interface {
	SelectByName(name string, tx ...orm.TxOrmer) (*tagDO.TagDO, error)
	InsertOrUpdate(tag *tagDO.TagDO, tx ...orm.TxOrmer) error
	DeleteByName(name string, tx ...orm.TxOrmer) error
}

type TagMapperImpl struct{}

func NewTagMapper() *TagMapperImpl {
	return &TagMapperImpl{}
}

func (impl *TagMapperImpl) SelectByName(name string, tx ...orm.TxOrmer) (*tagDO.TagDO, error) {
	executor := getQueryExecutor(tx...)
	var record tagDAO.TagRecord
	err := executor.QueryTable(tagDAO.TableTag).Filter("name", name).One(&record)
	if err != nil {
		return nil, err
	}

	// todo 有更复杂处理
	var tagArt tagDO.TagArt
	jsonUtil.UnmarshalJsonString(record.Details, &tagArt)
	do := &tagDO.TagDO{
		Id:       record.Id,
		CreateAt: record.CreatedAt.String(),
		Name:     record.Name,
		Creator:  record.Creator,
		Art:      tagArt,
	}
	return do, nil
}

func (impl *TagMapperImpl) InsertOrUpdate(tag *tagDO.TagDO, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)

	// todo 有更复杂处理
	record := tagDAO.TagRecord{
		Name:    tag.Name,
		Creator: tag.Creator,
		Details: jsonUtil.GetJsonString(tag.Art),
	}
	_, err := executor.InsertOrUpdate(&record)
	return err
}

func (impl *TagMapperImpl) DeleteByName(name string, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)
	_, err := executor.QueryTable(tagDAO.TableTag).Filter("name", name).Delete()
	return err
}
