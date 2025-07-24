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
func (c *UserController) Login() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		token := facade.NewUserFacade().Login(&c.Controller)
		return templates.NewJsonTemplate200(token)
	})
}

// @router login/status [get]
func (c *UserController) LoginStatus() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		profile, loggedIn := facade.NewUserFacade().LoginStatus(&c.Controller)
		return templates.NewJsonTemplate200(map[string]interface{}{
			"loggedIn": loggedIn,
			"profile":  profile,
		})
	})
}

type UserAuthController struct {
	web.Controller
}

// @router profile [get]
func (c *UserAuthController) Profile() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		profile := facade.NewUserFacade().GetProfile(&c.Controller)
		return templates.NewJsonTemplate200(profile)
	})
}

// @router logout [get]
func (c *UserAuthController) Logout() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		tokenStr := c.Ctx.Request.Header.Get("Authorization")
		facade.NewUserFacade().Logout(tokenStr)
		return templates.NewJsonTemplate200(nil)
	})
}
