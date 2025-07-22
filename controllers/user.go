package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"media-station/controllers/templates"
	"media-station/facade"
)

type UserController struct {
	web.Controller
}

// @router login [post]
func (c *VideoController) Login() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		token := facade.NewUserFacade().Login(&c.Controller)
		return templates.NewJsonTemplate200(token)
	})
}

// @router logout [get]
func (c *VideoController) Logout() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		tokenStr := c.Ctx.Request.Header.Get("Authorization")
		facade.NewUserFacade().Logout(tokenStr)
		return templates.NewJsonTemplate200(nil)
	})
}
