package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

	beego.GlobalControllerRouter["media-station/controllers:ActorAuthController"] = append(beego.GlobalControllerRouter["media-station/controllers:ActorAuthController"],
		beego.ControllerComments{
			Method:           "Create",
			Router:           `create`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:ActorAuthController"] = append(beego.GlobalControllerRouter["media-station/controllers:ActorAuthController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `delete/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:ActorAuthController"] = append(beego.GlobalControllerRouter["media-station/controllers:ActorAuthController"],
		beego.ControllerComments{
			Method:           "Update",
			Router:           `update`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:ActorController"] = append(beego.GlobalControllerRouter["media-station/controllers:ActorController"],
		beego.ControllerComments{
			Method:           "Cover",
			Router:           `cover/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:ActorController"] = append(beego.GlobalControllerRouter["media-station/controllers:ActorController"],
		beego.ControllerComments{
			Method:           "Page",
			Router:           `page/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:ActorController"] = append(beego.GlobalControllerRouter["media-station/controllers:ActorController"],
		beego.ControllerComments{
			Method:           "SearchActor",
			Router:           `search`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:GalleryAuthController"] = append(beego.GlobalControllerRouter["media-station/controllers:GalleryAuthController"],
		beego.ControllerComments{
			Method:           "Picture",
			Router:           `:id/:page`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:GalleryAuthController"] = append(beego.GlobalControllerRouter["media-station/controllers:GalleryAuthController"],
		beego.ControllerComments{
			Method:           "Upload",
			Router:           `upload`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:GalleryController"] = append(beego.GlobalControllerRouter["media-station/controllers:GalleryController"],
		beego.ControllerComments{
			Method:           "Page",
			Router:           `page/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:GalleryController"] = append(beego.GlobalControllerRouter["media-station/controllers:GalleryController"],
		beego.ControllerComments{
			Method:           "SearchGallery",
			Router:           `search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:UserAuthController"] = append(beego.GlobalControllerRouter["media-station/controllers:UserAuthController"],
		beego.ControllerComments{
			Method:           "Logout",
			Router:           `logout`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:UserAuthController"] = append(beego.GlobalControllerRouter["media-station/controllers:UserAuthController"],
		beego.ControllerComments{
			Method:           "Profile",
			Router:           `profile`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:UserController"] = append(beego.GlobalControllerRouter["media-station/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           `login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:UserController"] = append(beego.GlobalControllerRouter["media-station/controllers:UserController"],
		beego.ControllerComments{
			Method:           "LoginStatus",
			Router:           `login/status`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:VideoAuthController"] = append(beego.GlobalControllerRouter["media-station/controllers:VideoAuthController"],
		beego.ControllerComments{
			Method:           "Play",
			Router:           `play/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:VideoAuthController"] = append(beego.GlobalControllerRouter["media-station/controllers:VideoAuthController"],
		beego.ControllerComments{
			Method:           "Upload",
			Router:           `upload`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:VideoController"] = append(beego.GlobalControllerRouter["media-station/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "Page",
			Router:           `page/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:VideoController"] = append(beego.GlobalControllerRouter["media-station/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "Play",
			Router:           `play/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["media-station/controllers:VideoController"] = append(beego.GlobalControllerRouter["media-station/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "SearchVideo",
			Router:           `search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
