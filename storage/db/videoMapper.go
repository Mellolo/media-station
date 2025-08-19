package db

import (
	"github.com/Mellolo/media-station/models/dao/daoCommon"
	"github.com/Mellolo/media-station/models/dao/videoDAO"
	"github.com/Mellolo/media-station/models/do/videoDO"
	"github.com/beego/beego/v2/client/orm"
)

type VideoMapper interface {
	Insert(video videoDO.VideoDO, tx ...orm.TxOrmer) (int64, error)
	SelectById(id int64, tx ...orm.TxOrmer) (videoDO.VideoDO, error)
	SelectAllLimit(limit int, tx ...orm.TxOrmer) ([]videoDO.VideoDO, error)
	SelectByKeyword(keyword string, tx ...orm.TxOrmer) ([]videoDO.VideoDO, error)
	Update(id int64, video videoDO.VideoDO, tx ...orm.TxOrmer) error
	DeleteById(id int64, tx ...orm.TxOrmer) error
}

type VideoMapperImpl struct{}

func NewVideoMapper() *VideoMapperImpl {
	return &VideoMapperImpl{}
}

func (impl *VideoMapperImpl) Insert(video videoDO.VideoDO, tx ...orm.TxOrmer) (int64, error) {
	executor := getQueryExecutor(tx...)

	record := videoDAO.VideoRecord{
		Name:            video.Name,
		Description:     video.Description,
		Uploader:        video.Uploader,
		VideoUrl:        video.VideoUrl,
		CoverUrl:        video.CoverUrl,
		Duration:        video.Duration,
		PermissionLevel: video.PermissionLevel,
	}

	return executor.Insert(&record)
}

func (impl *VideoMapperImpl) SelectById(id int64, tx ...orm.TxOrmer) (videoDO.VideoDO, error) {
	executor := getQueryExecutor(tx...)
	record := videoDAO.VideoRecord{CommonColumn: daoCommon.CommonColumn{Id: id}}
	err := executor.Read(&record)
	if err != nil {
		return videoDO.VideoDO{}, err
	}

	do := videoDO.VideoDO{
		Id:              record.Id,
		CreateAt:        record.CreatedAt.String(),
		Name:            record.Name,
		Description:     record.Description,
		Uploader:        record.Uploader,
		VideoUrl:        record.VideoUrl,
		CoverUrl:        record.CoverUrl,
		Duration:        record.Duration,
		PermissionLevel: record.PermissionLevel,
	}
	return do, nil
}

func (impl *VideoMapperImpl) SelectAllLimit(limit int, tx ...orm.TxOrmer) ([]videoDO.VideoDO, error) {
	executor := getQueryExecutor(tx...)

	var records []videoDAO.VideoRecord
	_, err := executor.QueryTable(videoDAO.TableVideo).Limit(limit).All(&records)
	if err != nil {
		return nil, err
	}

	var doList []videoDO.VideoDO
	for _, record := range records {
		do := videoDO.VideoDO{
			Id:              record.Id,
			CreateAt:        record.CreatedAt.String(),
			Name:            record.Name,
			Description:     record.Description,
			Uploader:        record.Uploader,
			VideoUrl:        record.VideoUrl,
			CoverUrl:        record.CoverUrl,
			Duration:        record.Duration,
			PermissionLevel: record.PermissionLevel,
		}
		doList = append(doList, do)
	}
	return doList, nil
}

func (impl *VideoMapperImpl) SelectByKeyword(keyword string, tx ...orm.TxOrmer) ([]videoDO.VideoDO, error) {
	executor := getQueryExecutor(tx...)

	var records []videoDAO.VideoRecord
	cond := orm.NewCondition().Or("name__icontains", keyword).Or("description__icontains", keyword)
	_, err := executor.QueryTable(videoDAO.TableVideo).SetCond(cond).All(&records)
	if err != nil {
		return nil, err
	}

	var doList []videoDO.VideoDO
	for _, record := range records {
		do := videoDO.VideoDO{
			Id:              record.Id,
			CreateAt:        record.CreatedAt.String(),
			Name:            record.Name,
			Description:     record.Description,
			Uploader:        record.Uploader,
			VideoUrl:        record.VideoUrl,
			CoverUrl:        record.CoverUrl,
			Duration:        record.Duration,
			PermissionLevel: record.PermissionLevel,
		}
		doList = append(doList, do)
	}
	return doList, nil
}

func (impl *VideoMapperImpl) Update(id int64, video videoDO.VideoDO, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)

	record := videoDAO.VideoRecord{
		CommonColumn: daoCommon.CommonColumn{
			Id: id,
		},
		Name:            video.Name,
		Description:     video.Description,
		Uploader:        video.Uploader,
		VideoUrl:        video.VideoUrl,
		CoverUrl:        video.CoverUrl,
		Duration:        video.Duration,
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
