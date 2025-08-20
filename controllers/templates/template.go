package templates

import (
	"fmt"
	"github.com/Mellolo/common/errors"
	"github.com/Mellolo/media-station/util"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/google/uuid"
	"io"
	"net/http"
)

const KeyExceptionData = "exceptionData"

type ExceptionData struct {
	Uuid     string // 日志里的uuid
	ErrorMsg string // 错误信息
	Stack    string // 错误堆栈
}

// ServeJsonTemplate 用于业务controller的返回模板
func ServeJsonTemplate(ctx *context.Context, f func() JsonTemplate) {
	var response JsonTemplate
	panicContext := errors.CatchPanic(func() {
		response = f()
	})
	if panicContext.Err != nil {
		uniqueId, _ := uuid.NewV7()
		logs.Error(
			fmt.Sprintf("error url(%s)\n%s",
				ctx.Input.URL(),
				util.FormatErrorLog(uniqueId.String(), panicContext.Err.Error(), panicContext.RecoverStack),
			))
		data := ExceptionData{
			ErrorMsg: panicContext.Err.Error(),
		}
		if web.BConfig.RunMode == web.DEV && web.BConfig.EnableErrorsRender {
			data.Uuid = uniqueId.String()
			data.Stack = panicContext.RecoverStack
		}
		ctx.Input.SetData(KeyExceptionData, data)
		web.Exception(500, ctx)
	}
	err := ctx.JSONResp(response)
	if err != nil {
		logs.Critical(fmt.Sprintf("ServeJson failed: %v", err))
		ctx.Input.SetData(KeyExceptionData, ExceptionData{
			ErrorMsg: fmt.Sprintf("ServeJson failed: %v", err),
		})
		web.Exception(500, ctx)
	}
}

// ServeErrorJsonTemplate 用于error controller的返回模板
func ServeErrorJsonTemplate(ctx *context.Context, f func() JsonTemplate) {
	var response JsonTemplate
	panicContext := errors.CatchPanic(func() {
		response = f()
	})
	if panicContext.Err != nil {
		logs.Critical(fmt.Sprintf("error controller failed: %v\n[stack]%s\n", panicContext.Err, panicContext.RecoverStack))
		panic(errors.WrapError(panicContext.Err, "error controller failed"))
	}
	err := ctx.JSONResp(response)
	if err != nil {
		logs.Critical(fmt.Sprintf("ServeJson failed: %v", err))
		panic(errors.WrapError(err, "ServeJson failed"))
	}
}

func ServePictureTemplate(ctx *context.Context, f func() PictureTemplate) {
	var response PictureTemplate
	defer func() {
		if response.Reader != nil {
			_ = response.Reader.Close()
		}
	}()

	panicContext := errors.CatchPanic(func() {
		response = f()
	})
	if panicContext.Err != nil {
		uniqueId, _ := uuid.NewV7()
		logs.Error(
			fmt.Sprintf("error url(%s)\n%s",
				ctx.Input.URL(),
				util.FormatErrorLog(uniqueId.String(), panicContext.Err.Error(), panicContext.RecoverStack),
			))
		data := ExceptionData{
			ErrorMsg: panicContext.Err.Error(),
		}
		if web.BConfig.RunMode == web.DEV && web.BConfig.EnableErrorsRender {
			data.Uuid = uniqueId.String()
			data.Stack = panicContext.RecoverStack
		}
		ctx.Input.SetData(KeyExceptionData, data)
		web.Exception(500, ctx)
		return
	}

	// 响应头设置
	ctx.Output.Header("Content-Type", response.Header.Get("Content-Type"))
	ctx.ResponseWriter.Header().Set("Content-Length", response.Header.Get("Content-Length"))

	_, err := io.Copy(ctx.ResponseWriter, response.Reader)
	if err != nil {
		uniqueId, _ := uuid.NewV7()
		logs.Error(
			fmt.Sprintf("error url(%s)\n%s",
				ctx.Input.URL(),
				util.FormatErrorLog(uniqueId.String(), err.Error()),
			))
		data := ExceptionData{
			ErrorMsg: err.Error(),
		}
		if web.BConfig.RunMode == web.DEV && web.BConfig.EnableErrorsRender {
			data.Uuid = uniqueId.String()
		}
		ctx.Input.SetData(KeyExceptionData, data)
		web.Exception(500, ctx)
		return
	}
}

func ServeVideoTemplate(ctx *context.Context, f func() VideoTemplate) {
	var response VideoTemplate
	defer func() {
		if response.Reader != nil {
			_ = response.Reader.Close()
		}
	}()

	panicContext := errors.CatchPanic(func() {
		response = f()
	})
	if panicContext.Err != nil {
		uniqueId, _ := uuid.NewV7()
		logs.Error(
			fmt.Sprintf("error url(%s)\n%s",
				ctx.Input.URL(),
				util.FormatErrorLog(uniqueId.String(), panicContext.Err.Error(), panicContext.RecoverStack),
			))
		data := ExceptionData{
			ErrorMsg: panicContext.Err.Error(),
		}
		if web.BConfig.RunMode == web.DEV && web.BConfig.EnableErrorsRender {
			data.Uuid = uniqueId.String()
			data.Stack = panicContext.RecoverStack
		}
		ctx.Input.SetData(KeyExceptionData, data)
		web.Exception(500, ctx)
		return
	}

	// 响应头设置
	ctx.Output.Header("Content-Type", response.Header.Get("Content-Type"))
	ctx.Output.Header("Accept-Ranges", response.Header.Get("Accept-Ranges"))
	ctx.ResponseWriter.Header().Set("Content-Length", response.Header.Get("Content-Length"))
	if contentRange := response.Header.Get("Content-Range"); contentRange != "" {
		ctx.Output.Header("Content-Range", contentRange)
		ctx.ResponseWriter.WriteHeader(http.StatusPartialContent)
	} else {
		ctx.ResponseWriter.WriteHeader(http.StatusOK)
	}

	_, err := io.Copy(ctx.ResponseWriter, response.Reader)
	if err != nil {
		uniqueId, _ := uuid.NewV7()
		logs.Error(
			fmt.Sprintf("error url(%s)\n%s",
				ctx.Input.URL(),
				util.FormatErrorLog(uniqueId.String(), err.Error()),
			))
		data := ExceptionData{
			ErrorMsg: err.Error(),
		}
		if web.BConfig.RunMode == web.DEV && web.BConfig.EnableErrorsRender {
			data.Uuid = uniqueId.String()
		}
		ctx.Input.SetData(KeyExceptionData, data)
		web.Exception(500, ctx)
		return
	}
}
