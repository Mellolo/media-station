package bizGallery

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/mellolo/common/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"media-station/enum"
	"media-station/generator"
	"media-station/models/do/galleryDO"
	"media-station/models/dto/fileDTO"
	"media-station/models/dto/galleryDTO"
	"media-station/storage/db"
	"media-station/storage/oss"
)

const (
	galleryIdGenerateKey = "gallery"
	bucketGallery        = "gallery"
)

type GalleryBizService interface {
	GetGalleryPage(id int64, tx ...orm.TxOrmer) galleryDTO.GalleryPageDTO
	SearchGallery(searchDTO galleryDTO.GallerySearchDTO, tx ...orm.TxOrmer) []galleryDTO.GalleryItemDTO
	CreateGallery(createDTO galleryDTO.GalleryCreateDTO, picDTOList []fileDTO.FileDTO, ch chan string, tx ...orm.TxOrmer) int64
	UpdateGallery(id int64, updateDTO galleryDTO.GalleryUpdateDTO, tx ...orm.TxOrmer)
	DeleteGallery(id int64, tx ...orm.TxOrmer) (string, int)
	ShowGalleryPage(id int64, page int) galleryDTO.PictureFileDTO

	RemoveGalleryDir(dir string, count int)
}

func NewGalleryBizService() *GalleryBizServiceImpl {
	return &GalleryBizServiceImpl{
		galleryMapper:  db.NewGalleryMapper(),
		idGenerator:    generator.NewIdGenerator(),
		pictureStorage: oss.NewPictureStorage(),
	}
}

type GalleryBizServiceImpl struct {
	galleryMapper  db.GalleryMapper
	idGenerator    generator.IdGenerator
	pictureStorage oss.PictureStorage
}

func (impl *GalleryBizServiceImpl) GetGalleryPage(id int64, tx ...orm.TxOrmer) galleryDTO.GalleryPageDTO {
	gallery, err := impl.galleryMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("get gallery [%d] failed", id)))
	}
	return galleryDTO.GalleryPageDTO{
		Id:              gallery.Id,
		Name:            gallery.Name,
		Description:     gallery.Description,
		Actors:          gallery.Actors,
		Tags:            gallery.Tags,
		Uploader:        gallery.Uploader,
		CoverUrl:        gallery.CoverUrl,
		GalleryUrl:      gallery.GalleryUrl,
		PermissionLevel: gallery.PermissionLevel,
	}
}

func (impl *GalleryBizServiceImpl) SearchGallery(searchDTO galleryDTO.GallerySearchDTO, tx ...orm.TxOrmer) []galleryDTO.GalleryItemDTO {
	// 读取数据库
	var galleryDOList []*galleryDO.GalleryDO
	if searchDTO.Keyword == "" {
		doList, err := impl.galleryMapper.SelectAllLimit(200, tx...)
		if err != nil {
			panic(errors.WrapError(err, "select all Gallery error"))
		}
		galleryDOList = append(galleryDOList, doList...)
	} else {
		doList, err := impl.galleryMapper.SelectByKeyword(searchDTO.Keyword)
		if err != nil {
			panic(errors.WrapError(err, fmt.Sprintf("select Gallery by keyword [%s] error", searchDTO.Keyword)))
		}
		galleryDOList = append(galleryDOList, doList...)
	}

	// 再筛选，并转化为DTO
	var items []galleryDTO.GalleryItemDTO
	for _, do := range galleryDOList {
		if do.PermissionLevel == enum.PermissionForbidden || do.PermissionLevel == enum.PermissionPrivate {
			continue
		}
		if !sets.NewInt64(do.Actors...).HasAll(searchDTO.Actors...) {
			continue
		}
		if !sets.NewString(do.Tags...).HasAll(searchDTO.Tags...) {
			continue
		}

		items = append(items, galleryDTO.GalleryItemDTO{
			Id:              do.Id,
			Name:            do.Name,
			CoverUrl:        do.CoverUrl,
			PermissionLevel: do.PermissionLevel,
		})
	}

	return items
}

func (impl *GalleryBizServiceImpl) CreateGallery(createDTO galleryDTO.GalleryCreateDTO, picDTOList []fileDTO.FileDTO, ch chan string, tx ...orm.TxOrmer) int64 {
	defer func() {
		if ch != nil {
			close(ch)
		}
	}()

	if len(picDTOList) == 0 {
		panic(errors.NewError("gallery is empty"))
	}

	gallery := &galleryDO.GalleryDO{
		Name:        createDTO.Name,
		Description: createDTO.Description,
		PageCount:   len(picDTOList),
		Actors:      createDTO.Actors,
		Tags:        createDTO.Tags,
		Uploader:    createDTO.Uploader,
	}
	if sets.NewString(enum.PermissionLevels...).Has(createDTO.PermissionLevel) {
		gallery.PermissionLevel = createDTO.PermissionLevel
	} else {
		gallery.PermissionLevel = enum.PermissionPublic
	}
	// 上传画廊资源
	dir := impl.idGenerator.GenerateId(galleryIdGenerateKey)
	for i, picDTO := range picDTOList {
		path := fmt.Sprintf("%s/%d.jpg", dir, i+1)
		impl.pictureStorage.Upload(bucketGallery, path, picDTO.File, picDTO.Size)
		if ch != nil {
			ch <- fmt.Sprintf("%.2f", float64(i+1)/float64(len(picDTOList)))
		}
	}
	if ch != nil {
		ch <- "complete"
	}
	gallery.GalleryUrl = dir
	gallery.CoverUrl = fmt.Sprintf("%s/1.jpg", dir)

	id, err := impl.galleryMapper.Insert(gallery, tx...)
	if err != nil {
		panic(errors.WrapError(err, "create gallery failed"))
	}
	return id
}

func (impl *GalleryBizServiceImpl) UpdateGallery(id int64, updateDTO galleryDTO.GalleryUpdateDTO, tx ...orm.TxOrmer) {
	gallery, err := impl.galleryMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("gallery [%d] doesn't exist", id)))
	}

	// 更新gallery
	if updateDTO.Name != "" {
		gallery.Name = updateDTO.Name
	}
	if updateDTO.Description != "" {
		gallery.Description = updateDTO.Description
	}
	if len(updateDTO.Actors) > 0 {
		gallery.Actors = updateDTO.Actors
	}
	if len(updateDTO.Tags) > 0 {
		gallery.Tags = updateDTO.Tags
	}
	if sets.NewString(enum.PermissionLevels...).Has(updateDTO.PermissionLevel) {
		gallery.PermissionLevel = updateDTO.PermissionLevel
	}

	// 更新写入数据库
	err = impl.galleryMapper.Update(id, gallery, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("update gallery [%d] failed", id)))
	}
}

func (impl *GalleryBizServiceImpl) DeleteGallery(id int64, tx ...orm.TxOrmer) (string, int) {
	gallery, err := impl.galleryMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("gallery [%d] doesn't exist", id)))
	}
	err = impl.galleryMapper.DeleteById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("delete gallery [%d] failed", id)))
	}

	return gallery.GalleryUrl, gallery.PageCount
}

func (impl *GalleryBizServiceImpl) ShowGalleryPage(id int64, page int) galleryDTO.PictureFileDTO {
	gallery, err := impl.galleryMapper.SelectById(id)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("gallery [%d] doesn't exist", id)))
	}
	if page > gallery.PageCount || page < 1 {
		panic(errors.NewError(fmt.Sprintf("gallery [%d] page [%d] doesn't exist", id, page)))
	}
	do := impl.pictureStorage.Download(bucketGallery, fmt.Sprintf("%s/%d.jpg", gallery.GalleryUrl, page))
	return galleryDTO.PictureFileDTO{
		Header: do.Header,
		Reader: do.Reader,
	}
}

func (impl *GalleryBizServiceImpl) RemoveGalleryDir(dir string, count int) {
	for i := 1; i <= count; i++ {
		impl.pictureStorage.Remove(bucketGallery, fmt.Sprintf("%s/%d.jpg", dir, i))
	}
}
