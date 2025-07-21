package templates

import (
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	// 错误请求处理模板
	web.ErrorController(&ErrorController{})
}

type ErrorController struct {
	web.Controller
}

func (c *ErrorController) Error400() {
	ServeErrorJsonTemplate(c.Ctx, func() JsonTemplate {
		errorMsg := getErrorMsg(c.Ctx, fmt.Sprintf("有问题的请求(URL地址:%s)", c.Ctx.Request.URL.String()))
		return NewJsonTemplate400(errorMsg)
	})
}

func (c *ErrorController) Error401() {
	ServeErrorJsonTemplate(c.Ctx, func() JsonTemplate {
		errorMsg := getErrorMsg(c.Ctx, fmt.Sprintf("当前请求未完成鉴权(URL地址:%s)", c.Ctx.Request.URL.String()))
		return NewJsonTemplate401(errorMsg)
	})
}

func (c *ErrorController) Error404() {
	ServeErrorJsonTemplate(c.Ctx, func() JsonTemplate {
		errorMsg := getErrorMsg(c.Ctx, fmt.Sprintf("当前请求地址(%s)不存在", c.Ctx.Request.RequestURI))
		return NewJsonTemplate404(errorMsg)
	})
}

func (c *ErrorController) Error416() {
	ServeErrorJsonTemplate(c.Ctx, func() JsonTemplate {
		errorMsg := getErrorMsg(c.Ctx, fmt.Sprintf("无法满足指定Range(URL地址:%s)", c.Ctx.Request.URL.String()))
		return NewJsonTemplate400(errorMsg)
	})
}

func (c *ErrorController) Error500() {
	ServeErrorJsonTemplate(c.Ctx, func() JsonTemplate {
		errorMsg := getErrorMsg(c.Ctx, fmt.Sprintf("系统检测到当前请求未能正常完成(URL地址：%s)", c.Ctx.Request.URL.String()))
		return NewJsonTemplate500(errorMsg)
	})
}

func getErrorMsg(ctx *context.Context, msg string) string {
	if data, ok := ctx.Input.GetData(KeyExceptionData).(ExceptionData); ok {
		if data.Uuid != "" {
			msg = fmt.Sprintf("%s\n[uuid] %s", msg, data.Uuid)
		}
		if data.ErrorMsg != "" {
			msg = fmt.Sprintf("%s\n[message] %s", msg, data.ErrorMsg)
		}
		if data.Stack != "" {
			msg = fmt.Sprintf("%s\n[stack]\n%s", msg, data.Stack)
		}
	}
	return msg
}
