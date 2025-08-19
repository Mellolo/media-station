package controllers

import (
	"github.com/Mellolo/media-station/controllers/templates"
	"github.com/Mellolo/media-station/facade"
	"github.com/beego/beego/v2/server/web"
)

type ActorController struct {
	web.Controller
}

// @router search [get]
func (c *ActorController) SearchActor() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		voList := facade.NewActorFacade().SearchActor(&c.Controller)
		return templates.NewJsonTemplate200(voList)
	})
}

// @router page/:id [get]
func (c *ActorController) Page() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		vo := facade.NewActorFacade().GetActorPage(&c.Controller)
		return templates.NewJsonTemplate200(vo)
	})
}

// @router cover/:id [get]
func (c *ActorController) Cover() {
	templates.ServePictureTemplate(c.Ctx, func() templates.PictureTemplate {
		vo := facade.NewActorFacade().GetActorCover(&c.Controller)
		return templates.PictureTemplate{
			Reader: vo.Reader,
			Header: vo.Header,
		}
	})
}

type ActorAuthController struct {
	web.Controller
}

// @router update [post]
func (c *ActorAuthController) Update() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		facade.NewActorFacade().UpdateActor(&c.Controller)
		return templates.NewJsonTemplate200(nil)
	})
}

// @router create [post]
func (c *ActorAuthController) Create() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		id := facade.NewActorFacade().CreateActor(&c.Controller)
		return templates.NewJsonTemplate200(id)
	})
}

// @router delete/:id [delete]
func (c *ActorAuthController) Delete() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		facade.NewActorFacade().DeleteActor(&c.Controller)
		return templates.NewJsonTemplate200(nil)
	})
}
