package bizGallery

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/mellolo/common/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"media-station/enum"
	"media-station/generator"
	"media-station/models/do/galleryDO"
	"media-station/models/do/userDO"
	"media-station/models/dto/contextDTO"
	"media-station/models/dto/fileDTO"
	"media-station/models/dto/galleryDTO"
	"media-station/service/domain/domainPermission"
	"media-station/storage/db"
	"media-station/storage/oss"
	"strconv"
	"strings"
)

const (
	galleryIdGenerateKey = "gallery"
	bucketGallery        = "gallery"
)

type GalleryBizService interface {
	GetGallery(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) galleryDTO.GalleryDTO
	GetGalleryCover(ctx contextDTO.ContextDTO, id int64) galleryDTO.PictureFileDTO
	SearchGallery(ctx contextDTO.ContextDTO, searchDTO galleryDTO.GallerySearchDTO, tx ...orm.TxOrmer) []galleryDTO.GalleryDTO
	SearchGalleryByKeyword(ctx contextDTO.ContextDTO, keyword string, tx ...orm.TxOrmer) []galleryDTO.GalleryDTO
	CreateGallery(ctx contextDTO.ContextDTO, createDTO galleryDTO.GalleryCreateDTO, picDTOList []fileDTO.FileDTO, tx ...orm.TxOrmer) int64
	UpdateGallery(ctx contextDTO.ContextDTO, updateDTO galleryDTO.GalleryUpdateDTO, picDTOList []fileDTO.FileDTO, tx ...orm.TxOrmer) galleryDTO.GalleryDTO
	DeleteGallery(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) galleryDTO.GalleryDTO
	ShowGalleryPage(ctx contextDTO.ContextDTO, id int64, page int) galleryDTO.PictureFileDTO

	RemoveGalleryDir(ctx contextDTO.ContextDTO, dir string, filenames []string)
}

func NewGalleryBizService() *GalleryBizServiceImpl {
	return &GalleryBizServiceImpl{
		userMapper:              db.NewUserMapper(),
		galleryMapper:           db.NewGalleryMapper(),
		idGenerator:             generator.NewIdGenerator(),
		pictureStorage:          oss.NewPictureStorage(),
		permissionDomainService: domainPermission.NewPermissionDomainService(),
	}
}

type GalleryBizServiceImpl struct {
	userMapper     db.UserMapper
	galleryMapper  db.GalleryMapper
	idGenerator    generator.IdGenerator
	pictureStorage oss.PictureStorage

	permissionDomainService domainPermission.PermissionDomainService
}

func (impl *GalleryBizServiceImpl) GetGallery(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) galleryDTO.GalleryDTO {
	gallery, err := impl.galleryMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("get gallery [%d] failed", id)))
	}
	return impl.convertGalleryDO2GalleryDTO(gallery)
}

func (impl *GalleryBizServiceImpl) GetGalleryCover(ctx contextDTO.ContextDTO, id int64) galleryDTO.PictureFileDTO {
	gallery, err := impl.galleryMapper.SelectById(id)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("gallery [%d] doesn't exist", id)))
	}
	if len(gallery.PicPaths) == 0 {
		return galleryDTO.PictureFileDTO{}
	}
	do := impl.pictureStorage.Download(bucketGallery, fmt.Sprintf("%s/%s", gallery.DirPath, gallery.PicPaths[0]))
	return galleryDTO.PictureFileDTO{
		Header: do.Header,
		Reader: do.Reader,
	}
}

func (impl *GalleryBizServiceImpl) SearchGallery(ctx contextDTO.ContextDTO, searchDTO galleryDTO.GallerySearchDTO, tx ...orm.TxOrmer) []galleryDTO.GalleryDTO {
	var galleryDOList []galleryDO.GalleryDO

	for _, id := range searchDTO.Ids {
		do, err := impl.galleryMapper.SelectById(id, tx...)
		if err != nil {
			panic(errors.WrapError(err, fmt.Sprintf("select gallery by id [%d] error", id)))
		}
		if searchDTO.Keyword != "" && !strings.Contains(do.Name, searchDTO.Keyword) && !strings.Contains(do.Description, searchDTO.Keyword) {
			continue
		}
		galleryDOList = append(galleryDOList, do)
	}

	var user userDO.UserDO
	if ctx.UserClaim.Username != "" {
		user, _ = impl.userMapper.SelectByUsername(ctx.UserClaim.Username, tx...)
	}
	// 再筛选，并转化为DTO
	var items []galleryDTO.GalleryDTO
	for _, do := range galleryDOList {
		if !impl.permissionDomainService.IsVisible(user, do.Uploader, do.PermissionLevel) {
			continue
		}

		items = append(items, impl.convertGalleryDO2GalleryDTO(do))
	}

	return items
}

func (impl *GalleryBizServiceImpl) SearchGalleryByKeyword(ctx contextDTO.ContextDTO, keyword string, tx ...orm.TxOrmer) []galleryDTO.GalleryDTO {
	var galleryDOList []galleryDO.GalleryDO
	if keyword != "" {
		doList, err := impl.galleryMapper.SelectByKeyword(keyword)
		if err != nil {
			panic(errors.WrapError(err, fmt.Sprintf("select Gallery by keyword [%s] error", keyword)))
		}
		galleryDOList = append(galleryDOList, doList...)
	} else {
		doList, err := impl.galleryMapper.SelectAllLimit(200, tx...)
		if err != nil {
			panic(errors.WrapError(err, "select all gallery error"))
		}
		galleryDOList = append(galleryDOList, doList...)
	}

	var user userDO.UserDO
	if ctx.UserClaim.Username != "" {
		user, _ = impl.userMapper.SelectByUsername(ctx.UserClaim.Username, tx...)
	}
	// 再筛选，并转化为DTO
	var items []galleryDTO.GalleryDTO
	for _, do := range galleryDOList {
		if !impl.permissionDomainService.IsVisible(user, do.Uploader, do.PermissionLevel) {
			continue
		}

		items = append(items, impl.convertGalleryDO2GalleryDTO(do))
	}

	return items
}

func (impl *GalleryBizServiceImpl) CreateGallery(ctx contextDTO.ContextDTO, createDTO galleryDTO.GalleryCreateDTO, picDTOList []fileDTO.FileDTO, tx ...orm.TxOrmer) int64 {
	if len(picDTOList) == 0 {
		panic(errors.NewError("gallery is empty"))
	}

	gallery := galleryDO.GalleryDO{
		Name:        createDTO.Name,
		Description: createDTO.Description,
		Uploader:    createDTO.Uploader,
	}
	if sets.NewString(enum.PermissionLevels...).Has(createDTO.PermissionLevel) {
		gallery.PermissionLevel = createDTO.PermissionLevel
	} else {
		gallery.PermissionLevel = enum.PermissionPublic
	}
	// 上传画廊资源
	dir := impl.idGenerator.GenerateId(galleryIdGenerateKey)
	var picPaths []string
	for _, picDTO := range picDTOList {
		fileName := fmt.Sprintf("%s.jpg", impl.idGenerator.GenerateId(dir))
		path := fmt.Sprintf("%s/%s", dir, fileName)
		impl.pictureStorage.Upload(bucketGallery, path, picDTO.File, picDTO.Size)
		picPaths = append(picPaths, fileName)
	}
	gallery.DirPath = dir
	gallery.PicPaths = picPaths

	id, err := impl.galleryMapper.Insert(gallery, tx...)
	if err != nil {
		panic(errors.WrapError(err, "create gallery failed"))
	}
	return id
}

func (impl *GalleryBizServiceImpl) UpdateGallery(ctx contextDTO.ContextDTO, updateDTO galleryDTO.GalleryUpdateDTO, picDTOList []fileDTO.FileDTO, tx ...orm.TxOrmer) galleryDTO.GalleryDTO {
	gallery, err := impl.galleryMapper.SelectById(updateDTO.Id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("gallery [%d] doesn't exist", updateDTO.Id)))
	}
	origin := *gallery.DeepCopy()

	// 更新gallery
	gallery.Name = updateDTO.Name
	gallery.Description = updateDTO.Description
	if sets.NewString(enum.PermissionLevels...).Has(updateDTO.PermissionLevel) {
		gallery.PermissionLevel = updateDTO.PermissionLevel
	}

	// 更新画廊资源
	dir := impl.idGenerator.GenerateId(galleryIdGenerateKey)
	var picPaths []string
	for _, pageDTO := range updateDTO.Pages {
		var picDTO fileDTO.FileDTO
		if pageDTO.IsNewUploaded {
			picDTO = picDTOList[pageDTO.Index-1]
		} else {
			picDO := impl.pictureStorage.Download(bucketGallery, fmt.Sprintf("%s/%s", gallery.DirPath, gallery.PicPaths[pageDTO.Index-1]))
			size, sizeErr := strconv.ParseInt(picDO.Header.Get("Content-Length"), 10, 64)
			if sizeErr != nil {
				panic(errors.WrapError(sizeErr, "parse pic size failed"))
			}
			picDTO = fileDTO.FileDTO{
				File: picDO.Reader,
				Size: size,
			}
		}

		fileName := fmt.Sprintf("%s.jpg", impl.idGenerator.GenerateId(dir))
		path := fmt.Sprintf("%s/%s", dir, fileName)
		impl.pictureStorage.Upload(bucketGallery, path, picDTO.File, picDTO.Size)
		picPaths = append(picPaths, fileName)
	}
	gallery.DirPath = dir
	gallery.PicPaths = picPaths

	// 更新写入数据库
	err = impl.galleryMapper.Update(updateDTO.Id, gallery, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("update gallery [%d] failed", updateDTO.Id)))
	}

	return impl.convertGalleryDO2GalleryDTO(origin)
}

func (impl *GalleryBizServiceImpl) DeleteGallery(ctx contextDTO.ContextDTO, id int64, tx ...orm.TxOrmer) galleryDTO.GalleryDTO {
	gallery, err := impl.galleryMapper.SelectById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("gallery [%d] doesn't exist", id)))
	}
	err = impl.galleryMapper.DeleteById(id, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("delete gallery [%d] failed", id)))
	}

	return impl.convertGalleryDO2GalleryDTO(gallery)
}

func (impl *GalleryBizServiceImpl) ShowGalleryPage(ctx contextDTO.ContextDTO, id int64, page int) galleryDTO.PictureFileDTO {
	gallery, err := impl.galleryMapper.SelectById(id)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("gallery [%d] doesn't exist", id)))
	}
	if page > len(gallery.PicPaths) || page < 1 {
		panic(errors.NewError(fmt.Sprintf("gallery [%d] page [%d] doesn't exist", id, page)))
	}
	do := impl.pictureStorage.Download(bucketGallery, fmt.Sprintf("%s/%s", gallery.DirPath, gallery.PicPaths[page-1]))
	return galleryDTO.PictureFileDTO{
		Header: do.Header,
		Reader: do.Reader,
	}
}

func (impl *GalleryBizServiceImpl) RemoveGalleryDir(ctx contextDTO.ContextDTO, dir string, filenames []string) {
	for _, filename := range filenames {
		impl.pictureStorage.Remove(bucketGallery, fmt.Sprintf("%s/%s", dir, filename))
	}
}

func (impl *GalleryBizServiceImpl) convertGalleryDO2GalleryDTO(do galleryDO.GalleryDO) galleryDTO.GalleryDTO {
	return galleryDTO.GalleryDTO{
		Id:              do.Id,
		Name:            do.Name,
		Description:     do.Description,
		Uploader:        do.Uploader,
		DirPath:         do.DirPath,
		PicPaths:        do.PicPaths,
		PermissionLevel: do.PermissionLevel,
	}
}
