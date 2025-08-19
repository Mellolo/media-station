package controllers

import (
	"fmt"
	"github.com/Mellolo/common/errors"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	"media-station/controllers/templates"
	"media-station/facade"
	"media-station/util"
)

type VideoController struct {
	web.Controller
}

// @router search [get]
func (c *VideoController) SearchVideo() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		voList := facade.NewVideoFacade().SearchVideo(&c.Controller)
		return templates.NewJsonTemplate200(voList)
	})
}

// @router search/tag [get]
func (c *VideoController) SearchVideoByTag() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		voList := facade.NewVideoFacade().SearchVideoByTag(&c.Controller)
		return templates.NewJsonTemplate200(voList)
	})
}

// @router recommend [get]
func (c *VideoController) RecommendVideo() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		voList := facade.NewVideoFacade().RecommendVideo(&c.Controller)
		return templates.NewJsonTemplate200(voList)
	})
}

// @router page/:id [get]
func (c *VideoController) Page() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		vo := facade.NewVideoFacade().GetVideoPage(&c.Controller)
		return templates.NewJsonTemplate200(vo)
	})
}

// @router cover/:id [get]
func (c *VideoController) Cover() {
	templates.ServePictureTemplate(c.Ctx, func() templates.PictureTemplate {
		vo := facade.NewVideoFacade().GetVideoCover(&c.Controller)
		return templates.PictureTemplate{
			Reader: vo.Reader,
			Header: vo.Header,
		}
	})
}

// @router play/:id [get]
func (c *VideoController) Play() {
	templates.ServeVideoTemplate(c.Ctx, func() templates.VideoTemplate {
		// todo Token
		token := c.GetString("token", "")
		logs.Info(token)
		vo := facade.NewVideoFacade().PlayVideo(&c.Controller)
		return templates.VideoTemplate{
			Reader: vo.Reader,
			Header: vo.Header,
		}
	})
}

type VideoAuthController struct {
	web.Controller
}

// @router upload [post]
func (c *VideoAuthController) Upload() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		// 上传业务流程
		go func() {
			panicContext := errors.CatchPanic(func() {
				facade.NewVideoFacade().UploadVideo(&c.Controller)
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

// @router update [post]
func (c *VideoAuthController) Update() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		facade.NewVideoFacade().UpdateVideo(&c.Controller)
		return templates.NewJsonTemplate200(nil)
	})
}

// @router delete/:id [delete]
func (c *VideoAuthController) Delete() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		facade.NewVideoFacade().DeleteVideo(&c.Controller)
		return templates.NewJsonTemplate200(nil)
	})
}
