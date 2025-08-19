package db

import (
	"github.com/Mellolo/media-station/models/dao/tagDAO"
	"github.com/Mellolo/media-station/models/do/tagDO"
	"github.com/beego/beego/v2/client/orm"
)

type TagMapper interface {
	InsertOrUpdateTagsOfArt(artType string, artId int64, items []tagDO.TagDO, tx ...orm.TxOrmer) error
	SelectArtByTag(artType string, tag string, tx ...orm.TxOrmer) ([]tagDO.TagDO, error)
	SelectTagByArt(artType string, artId int64, tx ...orm.TxOrmer) ([]tagDO.TagDO, error)
	DeleteArt(artType string, artId int64, tx ...orm.TxOrmer) error
	DeleteTag(tag string, tx ...orm.TxOrmer) error
}

func NewTagMapper() *TagMapperImpl {
	return &TagMapperImpl{}
}

type TagMapperImpl struct{}

func (impl *TagMapperImpl) InsertOrUpdateTagsOfArt(artType string, artId int64, items []tagDO.TagDO, tx ...orm.TxOrmer) error {
	err := impl.DeleteArt(artType, artId, tx...)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	executor := getQueryExecutor(tx...)

	var records []tagDAO.TagRecord
	for _, item := range items {
		records = append(records, tagDAO.TagRecord{
			ArtType: artType,
			ArtId:   artId,
			Tag:     item.Tag,
		})
	}

	_, err = executor.InsertMulti(100, records)
	return err
}

func (impl *TagMapperImpl) SelectArtByTag(artType string, tag string, tx ...orm.TxOrmer) ([]tagDO.TagDO, error) {
	executor := getQueryExecutor(tx...)

	var records []tagDAO.TagRecord
	_, err := executor.QueryTable(tagDAO.TableTag).
		Filter("art_type", artType).
		Filter("tag", tag).
		All(&records, "art_type", "art_id", "tag")
	if err != nil {
		return nil, err
	}

	var items []tagDO.TagDO
	for _, record := range records {
		do := tagDO.TagDO{
			ArtType: record.ArtType,
			ArtId:   record.ArtId,
			Tag:     record.Tag,
		}
		items = append(items, do)
	}

	return items, nil
}

func (impl *TagMapperImpl) SelectTagByArt(artType string, artId int64, tx ...orm.TxOrmer) ([]tagDO.TagDO, error) {
	executor := getQueryExecutor(tx...)

	var records []tagDAO.TagRecord
	_, err := executor.QueryTable(tagDAO.TableTag).
		Filter("art_type", artType).
		Filter("art_id", artId).
		All(&records, "art_type", "art_id", "tag")
	if err != nil {
		return nil, err
	}

	var items []tagDO.TagDO
	for _, record := range records {
		do := tagDO.TagDO{
			ArtType: record.ArtType,
			ArtId:   record.ArtId,
			Tag:     record.Tag,
		}
		items = append(items, do)
	}

	return items, nil
}

func (impl *TagMapperImpl) DeleteArt(artType string, artId int64, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)

	_, err := executor.QueryTable(tagDAO.TableTag).
		Filter("art_type", artType).
		Filter("art_id", artId).
		Delete()
	return err
}

func (impl *TagMapperImpl) DeleteTag(tag string, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)

	_, err := executor.QueryTable(tagDAO.TableTag).
		Filter("tag", tag).
		Delete()
	return err
}
