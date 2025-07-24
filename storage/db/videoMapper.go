package db

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/mellolo/common/utils/jsonUtil"
	"media-station/models/dao/daoCommon"
	"media-station/models/dao/videoDAO"
	"media-station/models/do/videoDO"
)

type VideoMapper interface {
	Insert(video *videoDO.VideoDO, tx ...orm.TxOrmer) (int64, error)
	SelectById(id int64, tx ...orm.TxOrmer) (*videoDO.VideoDO, error)
	SelectAllLimit(limit int, tx ...orm.TxOrmer) ([]*videoDO.VideoDO, error)
	SelectByKeyword(keyword string, tx ...orm.TxOrmer) ([]*videoDO.VideoDO, error)
	Update(id int64, video *videoDO.VideoDO, tx ...orm.TxOrmer) error
	DeleteById(id int64, tx ...orm.TxOrmer) error
}

type VideoMapperImpl struct{}

func NewVideoMapper() *VideoMapperImpl {
	return &VideoMapperImpl{}
}

func (impl *VideoMapperImpl) Insert(video *videoDO.VideoDO, tx ...orm.TxOrmer) (int64, error) {
	executor := getQueryExecutor(tx...)

	record := videoDAO.VideoRecord{
		Name:            video.Name,
		Description:     video.Description,
		Actors:          jsonUtil.GetJsonString(video.Actors),
		Tags:            jsonUtil.GetJsonString(video.Tags),
		Uploader:        video.Uploader,
		VideoUrl:        video.VideoUrl,
		CoverUrl:        video.CoverUrl,
		PermissionLevel: video.PermissionLevel,
	}

	return executor.Insert(&record)
}

func (impl *VideoMapperImpl) SelectById(id int64, tx ...orm.TxOrmer) (*videoDO.VideoDO, error) {
	executor := getQueryExecutor(tx...)
	record := videoDAO.VideoRecord{CommonColumn: daoCommon.CommonColumn{Id: id}}
	err := executor.Read(&record)
	if err != nil {
		return nil, err
	}

	var actors []int64
	var tags []string
	jsonUtil.UnmarshalJsonString(record.Actors, &actors)
	jsonUtil.UnmarshalJsonString(record.Tags, &tags)
	do := &videoDO.VideoDO{
		Id:              record.Id,
		CreateAt:        record.CreatedAt.String(),
		Name:            record.Name,
		Description:     record.Description,
		Actors:          actors,
		Tags:            tags,
		Uploader:        record.Uploader,
		VideoUrl:        record.VideoUrl,
		CoverUrl:        record.CoverUrl,
		PermissionLevel: record.PermissionLevel,
	}
	return do, nil
}

func (impl *VideoMapperImpl) SelectAllLimit(limit int, tx ...orm.TxOrmer) ([]*videoDO.VideoDO, error) {
	executor := getQueryExecutor(tx...)

	var records []videoDAO.VideoRecord
	_, err := executor.QueryTable(videoDAO.TableVideo).Limit(limit).All(&records)
	if err != nil {
		return nil, err
	}

	var doList []*videoDO.VideoDO
	for _, record := range records {
		var actors []int64
		var tags []string
		jsonUtil.UnmarshalJsonString(record.Actors, &actors)
		jsonUtil.UnmarshalJsonString(record.Tags, &tags)
		do := &videoDO.VideoDO{
			Id:              record.Id,
			CreateAt:        record.CreatedAt.String(),
			Name:            record.Name,
			Description:     record.Description,
			Actors:          actors,
			Tags:            tags,
			Uploader:        record.Uploader,
			VideoUrl:        record.VideoUrl,
			CoverUrl:        record.CoverUrl,
			PermissionLevel: record.PermissionLevel,
		}
		doList = append(doList, do)
	}
	return doList, nil
}

func (impl *VideoMapperImpl) SelectByKeyword(keyword string, tx ...orm.TxOrmer) ([]*videoDO.VideoDO, error) {
	executor := getQueryExecutor(tx...)

	var records []videoDAO.VideoRecord
	cond := orm.NewCondition().Or("name__icontains", keyword).Or("description__icontains", keyword)
	_, err := executor.QueryTable(videoDAO.TableVideo).SetCond(cond).All(&records)
	if err != nil {
		return nil, err
	}

	var doList []*videoDO.VideoDO
	for _, record := range records {
		var actors []int64
		var tags []string
		jsonUtil.UnmarshalJsonString(record.Actors, &actors)
		jsonUtil.UnmarshalJsonString(record.Tags, &tags)
		do := &videoDO.VideoDO{
			Id:              record.Id,
			CreateAt:        record.CreatedAt.String(),
			Name:            record.Name,
			Description:     record.Description,
			Actors:          actors,
			Tags:            tags,
			Uploader:        record.Uploader,
			VideoUrl:        record.VideoUrl,
			CoverUrl:        record.CoverUrl,
			PermissionLevel: record.PermissionLevel,
		}
		doList = append(doList, do)
	}
	return doList, nil
}

func (impl *VideoMapperImpl) Update(id int64, video *videoDO.VideoDO, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)

	record := videoDAO.VideoRecord{
		CommonColumn: daoCommon.CommonColumn{
			Id: id,
		},
		Name:            video.Name,
		Description:     video.Description,
		Actors:          jsonUtil.GetJsonString(video.Actors),
		Tags:            jsonUtil.GetJsonString(video.Tags),
		Uploader:        video.Uploader,
		VideoUrl:        video.VideoUrl,
		CoverUrl:        video.CoverUrl,
		PermissionLevel: video.PermissionLevel,
	}
	_, err := executor.Update(&record)
	return err
}

func (impl *VideoMapperImpl) DeleteById(id int64, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)
	record := &videoDAO.VideoRecord{CommonColumn: daoCommon.CommonColumn{Id: id}}
	_, err := executor.Delete(record)
	return err
}
