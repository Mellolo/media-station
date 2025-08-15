package bizTag

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/mellolo/common/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"media-station/models/do/tagDO"
	"media-station/models/dto/contextDTO"
	"media-station/models/dto/tagDTO"
	"media-station/storage/db"
)

type TagBizService interface {
	SelectArtByTag(ctx contextDTO.ContextDTO, artType string, tags []string, tx ...orm.TxOrmer) []int64
	SelectTagByArt(ctx contextDTO.ContextDTO, artType string, artId int64, tx ...orm.TxOrmer) []string
	InsertOrUpdateTagsOfArt(ctx contextDTO.ContextDTO, dto tagDTO.ArtTagDTO, tx ...orm.TxOrmer)
	DeleteArt(ctx contextDTO.ContextDTO, artType string, artId int64, tx ...orm.TxOrmer)
	DeleteTag(ctx contextDTO.ContextDTO, tag string, tx ...orm.TxOrmer)
}

func NewTagBizService() *TagBizServiceImpl {
	return &TagBizServiceImpl{
		tagMapper: db.NewTagMapper(),
	}
}

type TagBizServiceImpl struct {
	tagMapper db.TagMapper
}

func (impl TagBizServiceImpl) SelectArtByTag(ctx contextDTO.ContextDTO, artType string, tags []string, tx ...orm.TxOrmer) []int64 {
	if len(tags) == 0 {
		return nil
	}

	firstLoop := true
	artIdSet := sets.NewInt64()
	for _, tag := range tags {
		taggedList, err := impl.tagMapper.SelectArtByTag(artType, tag, tx...)
		if err != nil {
			panic(errors.WrapError(err, fmt.Sprintf("get %s of tag [%s] failed", artType, tag)))
		}
		if firstLoop {
			for _, tagged := range taggedList {
				artIdSet.Insert(tagged.ArtId)
			}
		} else {
			thisArtIdSet := sets.NewInt64()
			for _, tagged := range taggedList {
				thisArtIdSet.Insert(tagged.ArtId)
			}
			artIdSet = artIdSet.Intersection(thisArtIdSet)
		}
		firstLoop = false
	}

	return artIdSet.List()
}

func (impl TagBizServiceImpl) SelectTagByArt(ctx contextDTO.ContextDTO, artType string, artId int64, tx ...orm.TxOrmer) []string {
	taggedList, err := impl.tagMapper.SelectTagByArt(artType, artId, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("get tags of %s[%d] failed", artType, artId)))
	}
	var tags []string
	for _, tagged := range taggedList {
		tags = append(tags, tagged.Tag)
	}

	return tags
}

func (impl TagBizServiceImpl) InsertOrUpdateTagsOfArt(ctx contextDTO.ContextDTO, dto tagDTO.ArtTagDTO, tx ...orm.TxOrmer) {
	var doList []tagDO.TagDO
	for _, tag := range dto.Tags {
		doList = append(doList, tagDO.TagDO{
			ArtType: dto.ArtType,
			ArtId:   dto.ArtId,
			Tag:     tag,
		})
	}

	err := impl.tagMapper.InsertOrUpdateTagsOfArt(dto.ArtType, dto.ArtId, doList, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("insert tags of art [%s][%d] failed", dto.ArtType, dto.ArtId)))
	}
}

func (impl TagBizServiceImpl) DeleteArt(ctx contextDTO.ContextDTO, artType string, artId int64, tx ...orm.TxOrmer) {
	err := impl.tagMapper.DeleteArt(artType, artId, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("delete art [%s][%d] tags failed", artType, artId)))
	}
}

func (impl TagBizServiceImpl) DeleteTag(ctx contextDTO.ContextDTO, tag string, tx ...orm.TxOrmer) {
	err := impl.tagMapper.DeleteTag(tag, tx...)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("delete tag [%s] arts failed", tag)))
	}
}
