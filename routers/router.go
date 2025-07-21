package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"media-station/controllers"
)

func init() {
	namespaceGallery := beego.NewNamespace("/gallery",
		//beego.NSBefore(filters.JWTAuth),
		beego.NSInclude(&controllers.GalleryController{}),
	)

	namespaceVideo := beego.NewNamespace("/video",
		beego.NSInclude(&controllers.VideoController{}),
	)
	beego.AddNamespace(namespaceGallery, namespaceVideo)
}
