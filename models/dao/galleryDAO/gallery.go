package galleryDAO

import (
	"github.com/beego/beego/v2/client/orm"
	"media-station/models/dao/daoCommon"
)

const (
	TableGallery = "gallery"
)

type GalleryRecord struct {
	daoCommon.CommonColumn `orm:"inline"`

	Name            string `orm:"column(name)"`
	Description     string `orm:"column(description)"`
	PageCount       int    `orm:"column(page_count)"`
	Actors          string `orm:"column(actors)"`
	Tags            string `orm:"column(tags)"`
	Uploader        string `orm:"column(uploader)"`
	CoverUrl        string `orm:"column(cover_url)"`
	GalleryUrl      string `orm:"column(gallery_url)"`
	PermissionLevel string `orm:"column(permission_level)"`
}

func (g *GalleryRecord) TableName() string {
	return TableGallery
}

func init() {
	orm.RegisterModel(new(GalleryRecord))
}
