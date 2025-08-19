package filters

import (
	"github.com/Mellolo/common/cache"
	"github.com/Mellolo/common/config"
	"github.com/Mellolo/common/utils/jsonUtil"
	"github.com/Mellolo/common/utils/jwtUtil"
	"github.com/Mellolo/media-station/controllers/templates"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"k8s.io/apimachinery/pkg/util/sets"
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

	client, _ := cache.GetCache()
	data, _ := client.Get("userTokenBacklist")
	if blacklistStr, ok := data.(string); ok {
		var blacklist []string
		jsonUtil.UnmarshalJsonString(blacklistStr, &blacklist)
		if sets.NewString(blacklist...).Has(tokenStr) {
			ctx.Input.SetData(templates.KeyExceptionData, templates.ExceptionData{
				ErrorMsg: "invalid login",
			})
			web.Exception(401, ctx)
			return
		}
	}

	secretKey := config.GetConfig("media", "secretKey", "user")
	claim, err := jwtUtil.ParseToken(tokenStr, secretKey)
	if err != nil {
		ctx.Input.SetData(templates.KeyExceptionData, templates.ExceptionData{
			ErrorMsg: "invalid login",
		})
		web.Exception(401, ctx)
		return
	}

	ctx.Input.SetData(ContextClaim, claim)
}
