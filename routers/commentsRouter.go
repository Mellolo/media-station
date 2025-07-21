package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["media-station/controllers:GalleryController"] = append(beego.GlobalControllerRouter["media-station/controllers:GalleryController"],
        beego.ControllerComments{
            Method: "Picture",
            Router: `:id/:page`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["media-station/controllers:GalleryController"] = append(beego.GlobalControllerRouter["media-station/controllers:GalleryController"],
        beego.ControllerComments{
            Method: "Upload",
            Router: `upload`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["media-station/controllers:VideoController"] = append(beego.GlobalControllerRouter["media-station/controllers:VideoController"],
        beego.ControllerComments{
            Method: "Play",
            Router: `:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["media-station/controllers:VideoController"] = append(beego.GlobalControllerRouter["media-station/controllers:VideoController"],
        beego.ControllerComments{
            Method: "SearchGallery",
            Router: `search`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["media-station/controllers:VideoController"] = append(beego.GlobalControllerRouter["media-station/controllers:VideoController"],
        beego.ControllerComments{
            Method: "SearchVideo",
            Router: `search`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["media-station/controllers:VideoController"] = append(beego.GlobalControllerRouter["media-station/controllers:VideoController"],
        beego.ControllerComments{
            Method: "Upload",
            Router: `upload`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
