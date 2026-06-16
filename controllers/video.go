package controllers

import (
	"fmt"

	"github.com/Mellolo/common/cache"
	"github.com/Mellolo/common/config"
	"github.com/Mellolo/common/errors"
	"github.com/Mellolo/common/utils/jsonUtil"
	"github.com/Mellolo/common/utils/jwtUtil"
	"github.com/Mellolo/media-station/controllers/templates"
	"github.com/Mellolo/media-station/facade"
	"github.com/Mellolo/media-station/service/biz/bizVideo"
	"github.com/Mellolo/media-station/util"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/util/sets"
)

type VideoController struct {
	web.Controller
}

// validateTokenFromURL 从 URL 参数中验证 token
func (c *VideoController) validateTokenFromURL() error {
	tokenStr := c.GetString("token", "")
	if tokenStr == "" {
		return errors.NewError("require login")
	}

	// 检查黑名单
	client, _ := cache.GetCache()
	data, _ := client.Get("userTokenBacklist")
	if blacklistStr, ok := data.(string); ok {
		var blacklist []string
		jsonUtil.UnmarshalJsonString(blacklistStr, &blacklist)
		if sets.NewString(blacklist...).Has(tokenStr) {
			return errors.NewError("invalid login")
		}
	}

	// 验证 JWT token
	secretKey := config.GetConfig("media", "secretKey", "user")
	_, err := jwtUtil.ParseToken(tokenStr, secretKey)
	if err != nil {
		return errors.NewError("invalid login")
	}

	return nil
}

// @router search [get]
func (c *VideoController) SearchVideo() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		voList := facade.NewVideoFacade().SearchVideo(&c.Controller)
		return templates.NewJsonTemplate200(voList)
	})
}

// @router search/tag [get]
func (c *VideoController) SearchVideoByTag() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		voList := facade.NewVideoFacade().SearchVideoByTag(&c.Controller)
		return templates.NewJsonTemplate200(voList)
	})
}

// @router recommend [get]
func (c *VideoController) RecommendVideo() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		voList := facade.NewVideoFacade().RecommendVideo(&c.Controller)
		return templates.NewJsonTemplate200(voList)
	})
}

// @router page/:id [get]
func (c *VideoController) Page() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		vo := facade.NewVideoFacade().GetVideoPage(&c.Controller)
		return templates.NewJsonTemplate200(vo)
	})
}

// @router cover/:id [get]
func (c *VideoController) Cover() {
	templates.ServePictureTemplate(c.Ctx, func() templates.PictureTemplate {
		vo := facade.NewVideoFacade().GetVideoCover(&c.Controller)
		return templates.PictureTemplate{
			Reader: vo.Reader,
			Header: vo.Header,
		}
	})
}

// @router play/:id [get]
func (c *VideoController) Play() {
	templates.ServeVideoTemplate(c.Ctx, func() templates.VideoTemplate {
		// todo Token
		token := c.GetString("token", "")
		logs.Info(token)
		vo := facade.NewVideoFacade().PlayVideo(&c.Controller)
		return templates.VideoTemplate{
			Reader: vo.Reader,
			Header: vo.Header,
		}
	})
}

// @router stream/:id [get]
func (c *VideoController) StreamVideo() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		// 验证 token
		if err := c.validateTokenFromURL(); err != nil {
			return templates.NewJsonTemplate401(err.Error())
		}

		result := facade.NewVideoFacade().StreamVideo(&c.Controller)
		return templates.NewJsonTemplate200(result)
	})
}

// @router hls/:session/* [get]
func (c *VideoController) ServeHLSSegment() {
	// 验证 token
	if err := c.validateTokenFromURL(); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	sessionID := c.Ctx.Input.Param(":session")
	filename := c.Ctx.Input.Param(":path")

	filePath, err := facade.NewVideoFacade().GetHLSSegment(sessionID, filename)
	if err != nil {
		// 文件未就绪，返回 404
		c.Ctx.ResponseWriter.WriteHeader(404)
		c.Ctx.Output.JSON(map[string]string{"error": "Segment not ready"}, false, false)
		return
	}

	// 提供静态文件服务
	bizVideo.ServeHLSFile(filePath, c.Ctx.ResponseWriter, c.Ctx.Request)
}

type VideoAuthController struct {
	web.Controller
}

// @router upload [post]
func (c *VideoAuthController) Upload() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		// 上传业务流程
		go func() {
			panicContext := errors.CatchPanic(func() {
				facade.NewVideoFacade().UploadVideo(&c.Controller)
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

		return templates.NewJsonTemplate200(nil)
	})
}

// @router update [post]
func (c *VideoAuthController) Update() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		facade.NewVideoFacade().UpdateVideo(&c.Controller)
		return templates.NewJsonTemplate200(nil)
	})
}

// @router delete/:id [delete]
func (c *VideoAuthController) Delete() {
	templates.ServeJsonTemplate(c.Ctx, func() templates.JsonTemplate {
		facade.NewVideoFacade().DeleteVideo(&c.Controller)
		return templates.NewJsonTemplate200(nil)
	})
}
