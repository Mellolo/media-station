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
		response := NewJsonTemplate400(fmt.Sprintf("有问题的请求(URL地址:%s)", c.Ctx.Request.URL.String()))
		response.Data = formErrorData(c.Ctx)
		return response
	})
}

func (c *ErrorController) Error401() {
	ServeErrorJsonTemplate(c.Ctx, func() JsonTemplate {
		response := NewJsonTemplate401(fmt.Sprintf("当前请求未完成鉴权(URL地址:%s)", c.Ctx.Request.URL.String()))
		response.Data = formErrorData(c.Ctx)
		return response
	})
}

func (c *ErrorController) Error404() {
	ServeErrorJsonTemplate(c.Ctx, func() JsonTemplate {
		response := NewJsonTemplate404(fmt.Sprintf("当前请求地址(%s)不存在", c.Ctx.Request.RequestURI))
		response.Data = formErrorData(c.Ctx)
		return response
	})
}

func (c *ErrorController) Error416() {
	ServeErrorJsonTemplate(c.Ctx, func() JsonTemplate {
		response := NewJsonTemplate400(fmt.Sprintf("无法满足指定Range(URL地址:%s)", c.Ctx.Request.URL.String()))
		response.Data = formErrorData(c.Ctx)
		return response
	})
}

func (c *ErrorController) Error500() {
	ServeErrorJsonTemplate(c.Ctx, func() JsonTemplate {
		response := NewJsonTemplate500(fmt.Sprintf("系统检测到当前请求未能正常完成(URL地址：%s)", c.Ctx.Request.URL.String()))
		response.Data = formErrorData(c.Ctx)
		return response
	})
}

func formErrorData(ctx *context.Context) map[string]string {
	result := make(map[string]string)
	if data, ok := ctx.Input.GetData(KeyExceptionData).(ExceptionData); ok {
		if data.Uuid != "" {
			result["uuid"] = data.Uuid
		}
		if data.ErrorMsg != "" {
			result["errorMsg"] = data.ErrorMsg
		}
		if data.Stack != "" {
			result["stack"] = data.Stack
		}
	}
	return result
}
