package controllers

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mellolo/common/errors"
	"media-station/controllers/templates"
	"media-station/facade"
	"media-station/util"
	"net/http"
)

type GalleryController struct {
	web.Controller
}

// @router search [get]
func (c *GalleryController) SearchGallery() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		voList := facade.NewGalleryFacade().SearchGallery(&c.Controller)
		return templates.NewJsonTemplate200(voList)
	})
}

type GalleryAuthController struct {
	web.Controller
}

// @router upload [post]
func (c *GalleryAuthController) Upload() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源的连接，可加限制
			},
		}
		conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
		if err != nil {
			panic(errors.WrapError(err, "websocket error"))
		}

		ch := make(chan string, 100)

		// 上传业务流程
		go func() {
			panicContext := errors.CatchPanic(func() {
				facade.NewGalleryFacade().UploadGallery(&c.Controller, ch)
			})
			if panicContext.Err != nil {
				uniqueId, _ := uuid.NewV7()
				logs.Error(
					fmt.Sprintf("error url(%s)\n%s",
						c.Ctx.Input.URL(),
						util.FormatErrorLog(uniqueId.String(), panicContext.Err.Error(), panicContext.RecoverStack),
					))
			}
		}()

		// 进度条
		go func() {
			for str := range ch {
				err = conn.WriteMessage(websocket.TextMessage, []byte(str))
				if err != nil {
					return
				}
			}
		}()

		return templates.NewJsonTemplate200(nil)
	})
}

// @router :id/:page [get]
func (c *GalleryAuthController) Picture() {
	templates.ServePictureTemplate(c.Ctx, func() templates.PictureTemplate {

		vo := facade.NewGalleryFacade().ShowPage(&c.Controller)

		return templates.PictureTemplate{
			Reader: vo.Reader,
			Header: vo.Header,
		}
	})
}
