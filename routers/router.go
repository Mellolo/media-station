package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"media-station/controllers"
	"media-station/controllers/filters"
)

func init() {
	namespaceGallery := beego.NewNamespace("/gallery",
		beego.NSInclude(&controllers.GalleryController{}),
	)

	namespaceGalleryAuth := beego.NewNamespace("/galleryAuth",
		beego.NSBefore(filters.JWTAuth),
		beego.NSInclude(&controllers.GalleryAuthController{}),
	)

	namespaceVideo := beego.NewNamespace("/user",
		beego.NSInclude(&controllers.VideoAuthController{}),
	)

	namespaceVideoAuth := beego.NewNamespace("/videoAuth",
		beego.NSBefore(filters.JWTAuth),
		beego.NSInclude(&controllers.VideoController{}),
	)

	namespaceUser := beego.NewNamespace("/user",
		beego.NSInclude(&controllers.UserAuthController{}),
	)

	namespaceUserAuth := beego.NewNamespace("/userAuth",
		beego.NSBefore(filters.JWTAuth),
		beego.NSInclude(&controllers.UserAuthController{}),
	)

	beego.AddNamespace(
		namespaceGallery, namespaceVideo, namespaceUser,
		namespaceGalleryAuth, namespaceVideoAuth, namespaceUserAuth)
}
