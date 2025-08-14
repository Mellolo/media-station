package db

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/mellolo/common/utils/jsonUtil"
	"media-station/models/dao/daoCommon"
	"media-station/models/dao/galleryDAO"
	"media-station/models/do/galleryDO"
)

type GalleryMapper interface {
	Insert(gallery galleryDO.GalleryDO, tx ...orm.TxOrmer) (int64, error)
	SelectAllLimit(limit int, tx ...orm.TxOrmer) ([]galleryDO.GalleryDO, error)
	SelectByKeyword(keyword string, tx ...orm.TxOrmer) ([]galleryDO.GalleryDO, error)
	SelectById(id int64, tx ...orm.TxOrmer) (galleryDO.GalleryDO, error)
	Update(id int64, gallery galleryDO.GalleryDO, tx ...orm.TxOrmer) error
	DeleteById(id int64, tx ...orm.TxOrmer) error
}

type GalleryMapperImpl struct{}

func NewGalleryMapper() *GalleryMapperImpl {
	return &GalleryMapperImpl{}
}

func (impl *GalleryMapperImpl) Insert(gallery galleryDO.GalleryDO, tx ...orm.TxOrmer) (int64, error) {
	executor := getQueryExecutor(tx...)

	record := galleryDAO.GalleryRecord{
		Name:            gallery.Name,
		Description:     gallery.Description,
		Uploader:        gallery.Uploader,
		DirPath:         gallery.DirPath,
		PicPaths:        jsonUtil.GetJsonString(gallery.PicPaths),
		PermissionLevel: gallery.PermissionLevel,
	}

	return executor.Insert(&record)
}

func (impl *GalleryMapperImpl) SelectAllLimit(limit int, tx ...orm.TxOrmer) ([]galleryDO.GalleryDO, error) {
	executor := getQueryExecutor(tx...)

	var records []galleryDAO.GalleryRecord
	_, err := executor.QueryTable(galleryDAO.TableGallery).Limit(limit).All(&records)
	if err != nil {
		return nil, err
	}

	var doList []galleryDO.GalleryDO
	for _, record := range records {
		var picPaths []string
		jsonUtil.UnmarshalJsonString(record.PicPaths, &picPaths)
		do := galleryDO.GalleryDO{
			Id:              record.Id,
			CreateAt:        record.CreatedAt.String(),
			Name:            record.Name,
			Description:     record.Description,
			Uploader:        record.Uploader,
			DirPath:         record.DirPath,
			PicPaths:        picPaths,
			PermissionLevel: record.PermissionLevel,
		}
		doList = append(doList, do)
	}
	return doList, nil
}

func (impl *GalleryMapperImpl) SelectByKeyword(keyword string, tx ...orm.TxOrmer) ([]galleryDO.GalleryDO, error) {
	executor := getQueryExecutor(tx...)

	var records []galleryDAO.GalleryRecord
	cond := orm.NewCondition().Or("name__icontains", keyword).Or("description__icontains", keyword)
	_, err := executor.QueryTable(galleryDAO.TableGallery).SetCond(cond).All(&records)
	if err != nil {
		return nil, err
	}

	var doList []galleryDO.GalleryDO
	for _, record := range records {
		var picPaths []string
		jsonUtil.UnmarshalJsonString(record.PicPaths, &picPaths)
		do := galleryDO.GalleryDO{
			Id:              record.Id,
			CreateAt:        record.CreatedAt.String(),
			Name:            record.Name,
			Description:     record.Description,
			Uploader:        record.Uploader,
			DirPath:         record.DirPath,
			PicPaths:        picPaths,
			PermissionLevel: record.PermissionLevel,
		}
		doList = append(doList, do)
	}
	return doList, nil
}

func (impl *GalleryMapperImpl) SelectById(id int64, tx ...orm.TxOrmer) (galleryDO.GalleryDO, error) {
	executor := getQueryExecutor(tx...)
	record := galleryDAO.GalleryRecord{CommonColumn: daoCommon.CommonColumn{Id: id}}
	err := executor.Read(&record)
	if err != nil {
		return galleryDO.GalleryDO{}, err
	}

	var picPaths []string
	jsonUtil.UnmarshalJsonString(record.PicPaths, &picPaths)
	do := galleryDO.GalleryDO{
		Id:              record.Id,
		CreateAt:        record.CreatedAt.String(),
		Name:            record.Name,
		Description:     record.Description,
		Uploader:        record.Uploader,
		DirPath:         record.DirPath,
		PicPaths:        picPaths,
		PermissionLevel: record.PermissionLevel,
	}
	return do, nil
}

func (impl *GalleryMapperImpl) Update(id int64, gallery galleryDO.GalleryDO, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)

	record := galleryDAO.GalleryRecord{
		CommonColumn: daoCommon.CommonColumn{
			Id: id,
		},
		Name:            gallery.Name,
		Description:     gallery.Description,
		Uploader:        gallery.Uploader,
		DirPath:         gallery.DirPath,
		PicPaths:        jsonUtil.GetJsonString(gallery.PicPaths),
		PermissionLevel: gallery.PermissionLevel,
	}
	_, err := executor.Update(&record)
	return err
}

func (impl *GalleryMapperImpl) DeleteById(id int64, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)
	record := &galleryDAO.GalleryRecord{CommonColumn: daoCommon.CommonColumn{Id: id}}
	_, err := executor.Delete(record)
	return err
}
