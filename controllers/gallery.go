package controllers

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	"github.com/mellolo/common/errors"
	"media-station/controllers/templates"
	"media-station/facade"
	"media-station/util"
)

type GalleryController struct {
	web.Controller
}

// @router search [get]
func (c *GalleryController) SearchGallery() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		voList := facade.NewGalleryFacade().SearchGallery(&c.Controller)
		return templates.NewJsonTemplate200(voList)
	})
}

// @router search/tag [get]
func (c *VideoController) SearchGalleryByTag() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		voList := facade.NewGalleryFacade().SearchGalleryByTag(&c.Controller)
		return templates.NewJsonTemplate200(voList)
	})
}

// @router page/:id [get]
func (c *GalleryController) Page() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		vo := facade.NewGalleryFacade().GetGalleryPage(&c.Controller)
		return templates.NewJsonTemplate200(vo)
	})
}

// @router pic/:id/:page [get]
func (c *GalleryController) Picture() {
	templates.ServePictureTemplate(c.Ctx, func() templates.PictureTemplate {

		vo := facade.NewGalleryFacade().ShowPage(&c.Controller)

		return templates.PictureTemplate{
			Reader: vo.Reader,
			Header: vo.Header,
		}
	})
}

// @router cover/:id/ [get]
func (c *GalleryController) Cover() {
	templates.ServePictureTemplate(c.Ctx, func() templates.PictureTemplate {

		vo := facade.NewGalleryFacade().GetGalleryCover(&c.Controller)

		return templates.PictureTemplate{
			Reader: vo.Reader,
			Header: vo.Header,
		}
	})
}

type GalleryAuthController struct {
	web.Controller
}

// @router upload [post]
func (c *GalleryAuthController) Upload() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		// 上传业务流程
		go func() {
			panicContext := errors.CatchPanic(func() {
				facade.NewGalleryFacade().UploadGallery(&c.Controller)
			})
			if panicContext.Err != nil {
				uniqueId, _ := uuid.NewV7()
				logs.Error(
					fmt.Sprintf("error url(%s)\n%s",
						c.Ctx.Input.URL(),
						util.FormatErrorLog(uniqueId.String(), panicContext.Err.Error(), panicContext.RecoverStack),
					))
			}
		}()

		return templates.NewJsonTemplate200(nil)
	})
}

// @router delete/:id [delete]
func (c *GalleryAuthController) Delete() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		facade.NewGalleryFacade().DeleteGallery(&c.Controller)
		return templates.NewJsonTemplate200(nil)
	})
}

//// @router pic/:id/:page [get]
//func (c *GalleryAuthController) Picture() {
//	templates.ServePictureTemplate(c.Ctx, func() templates.PictureTemplate {
//
//		vo := facade.NewGalleryFacade().ShowPage(&c.Controller)
//
//		return templates.PictureTemplate{
//			Reader: vo.Reader,
//			Header: vo.Header,
//		}
//	})
//}
