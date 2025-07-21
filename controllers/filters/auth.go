package filters

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/mellolo/common/utils/jwtUtil"
	"media-station/controllers/templates"
)

const (
	ContextClaim = "claim"
)

func JWTAuth(ctx *context.Context) {
	tokenStr := ctx.Request.Header.Get("Authorization")
	if tokenStr == "" {
		ctx.Input.SetData(templates.KeyExceptionData, templates.ExceptionData{
			ErrorMsg: "require login",
		})
		web.Exception(401, ctx)
		return
	}

	claim, err := jwtUtil.ParseToken(tokenStr, "")
	if err != nil {
		ctx.Input.SetData(templates.KeyExceptionData, templates.ExceptionData{
			ErrorMsg: "invalid login",
		})
		web.Exception(401, ctx)
		return
	}

	ctx.Input.SetData(ContextClaim, claim)
}
