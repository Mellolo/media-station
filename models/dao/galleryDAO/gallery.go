package galleryDAO

import (
	"github.com/Mellolo/media-station/models/dao/daoCommon"
	"github.com/beego/beego/v2/client/orm"
)

const (
	TableGallery = "gallery"
)

type GalleryRecord struct {
	daoCommon.CommonColumn `orm:"inline"`

	Name            string `orm:"column(name)"`
	Description     string `orm:"column(description)"`
	Uploader        string `orm:"column(uploader)"`
	DirPath         string `orm:"column(dir_path)"`
	PicPaths        string `orm:"column(pic_paths)"`
	PermissionLevel string `orm:"column(permission_level)"`
}

func (g *GalleryRecord) TableName() string {
	return TableGallery
}

func init() {
	orm.RegisterModel(new(GalleryRecord))
}
