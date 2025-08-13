package bizTag

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/mellolo/common/errors"
	pkgErrors "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"media-station/models/do/tagDO"
	"media-station/models/dto/contextDTO"
	"media-station/models/dto/tagDTO"
	"media-station/storage/db"
)

type TagBizService interface {
	AddArt(ctx contextDTO.ContextDTO, dto tagDTO.TagCreateOrUpdateDTO, tx ...orm.TxOrmer)
	DeleteArt(ctx contextDTO.ContextDTO, dto tagDTO.TagDeleteArtDTO, tx ...orm.TxOrmer)
	RemoveArt(ctx contextDTO.ContextDTO, dto tagDTO.TagRemoveArtDTO, tx ...orm.TxOrmer)
	DeleteTag(ctx contextDTO.ContextDTO, tagName string, tx ...orm.TxOrmer)
}

type TagBizServiceImpl struct {
	tagMapper db.TagMapper
}

func NewTagBizService() *TagBizServiceImpl {
	return &TagBizServiceImpl{
		tagMapper: db.NewTagMapper(),
	}
}

func (impl *TagBizServiceImpl) AddArt(ctx contextDTO.ContextDTO, dto tagDTO.TagCreateOrUpdateDTO, tx ...orm.TxOrmer) {
	if dto.Name == "" {
		panic(errors.NewError("tag name cannot be empty"))
	}
	tag, err := impl.tagMapper.SelectByName(dto.Name, tx...)
	if err != nil && !pkgErrors.Is(err, orm.ErrNoRows) {
		panic(errors.WrapError(err, "select tag failed"))
	}
	if pkgErrors.Is(err, orm.ErrNoRows) {
		tag = &tagDO.TagDO{
			Name:    dto.Name,
			Creator: dto.Creator,
		}
	}

	videoIds := sets.NewInt64(tag.Art.VideoIds...)
	videoIds.Insert(dto.Details.VideoIds...)
	tag.Art.VideoIds = videoIds.List()

	galleryIds := sets.NewInt64(tag.Art.GalleryIds...)
	galleryIds.Insert(dto.Details.GalleryIds...)
	tag.Art.GalleryIds = galleryIds.List()

	err = impl.tagMapper.InsertOrUpdate(tag, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("create or update tag [%s] failed", dto.Name)))
	}
}

func (impl *TagBizServiceImpl) DeleteArt(ctx contextDTO.ContextDTO, dto tagDTO.TagDeleteArtDTO, tx ...orm.TxOrmer) {
	tag, err := impl.tagMapper.SelectByName(dto.Name, tx...)
	if err != nil && !pkgErrors.Is(err, orm.ErrNoRows) {
		panic(errors.WrapError(err, "select tag failed"))
	}
	if pkgErrors.Is(err, orm.ErrNoRows) {
		return
	}

	tag.Art.VideoIds = sets.NewInt64(tag.Art.VideoIds...).Delete(dto.Details.VideoIds...).List()

	tag.Art.GalleryIds = sets.NewInt64(tag.Art.GalleryIds...).Delete(dto.Details.GalleryIds...).List()

	err = impl.tagMapper.InsertOrUpdate(tag, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("create or update tag [%s] failed", dto.Name)))
	}
}

func (impl *TagBizServiceImpl) RemoveArt(ctx contextDTO.ContextDTO, dto tagDTO.TagRemoveArtDTO, tx ...orm.TxOrmer) {
	tag, err := impl.tagMapper.SelectByName(dto.Name, tx...)
	if err != nil && !pkgErrors.Is(err, orm.ErrNoRows) {
		panic(errors.WrapError(err, "select tag failed"))
	}
	tag.Art.VideoIds = sets.NewInt64(tag.Art.VideoIds...).Delete(dto.Details.VideoIds...).List()
	tag.Art.GalleryIds = sets.NewInt64(tag.Art.GalleryIds...).Delete(dto.Details.GalleryIds...).List()
	err = impl.tagMapper.InsertOrUpdate(tag, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("remove artwork for tag [%s] failed", dto.Name)))
	}
}

func (impl *TagBizServiceImpl) DeleteTag(ctx contextDTO.ContextDTO, tagName string, tx ...orm.TxOrmer) {
	err := impl.tagMapper.DeleteByName(tagName, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("delete tag [%s] failed", tagName)))
	}
}
