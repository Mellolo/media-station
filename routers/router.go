package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"media-station/controllers"
	"media-station/controllers/filters"
)

func init() {
	namespaceApi := beego.NewNamespace("/api",
		beego.NSNamespace("/gallery", beego.NSInclude(&controllers.GalleryController{})),
		beego.NSNamespace("/video", beego.NSInclude(&controllers.VideoController{})),
		beego.NSNamespace("/user", beego.NSInclude(&controllers.UserController{})),
		beego.NSNamespace("/auth",
			beego.NSBefore(filters.JWTAuth),
			beego.NSNamespace("/gallery", beego.NSInclude(&controllers.GalleryAuthController{})),
			beego.NSNamespace("/video", beego.NSInclude(&controllers.VideoAuthController{})),
			beego.NSNamespace("/user", beego.NSInclude(&controllers.UserAuthController{})),
		),
	)

	beego.AddNamespace(namespaceApi)
}
