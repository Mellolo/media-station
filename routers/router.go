package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"media-station/controllers"
)

func init() {
	namespaceGallery := beego.NewNamespace("/gallery",
		//beego.NSBefore(filters.JWTAuth),
		beego.NSInclude(&controllers.GalleryAuthController{}),
	)

	namespaceVideo := beego.NewNamespace("/video",
		beego.NSInclude(&controllers.VideoAuthController{}),
	)

	namespaceUser := beego.NewNamespace("/user",
		beego.NSInclude(&controllers.UserAuthController{}),
	)

	beego.AddNamespace(namespaceGallery, namespaceVideo, namespaceUser)
}
