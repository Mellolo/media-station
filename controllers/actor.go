package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"media-station/controllers/templates"
	"media-station/facade"
)

type ActorController struct {
	web.Controller
}

// @router page/:id [get]
func (c *ActorController) Page() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		vo := facade.NewActorFacade().GetActorPage(&c.Controller)
		return templates.NewJsonTemplate200(vo)
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
