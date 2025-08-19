package facade

import (
	"fmt"
	"github.com/Mellolo/common/errors"
	"github.com/Mellolo/common/utils/jsonUtil"
	"github.com/Mellolo/media-station/controllers/filters"
	"github.com/Mellolo/media-station/models/dto/contextDTO"
	"github.com/Mellolo/media-station/models/dto/userDTO"
	"github.com/beego/beego/v2/server/web"
	"io"
	"strconv"
)

type AbstractFacade struct {
}

func (facade *AbstractFacade) GetFile(c *web.Controller, key string) (io.ReadCloser, int64) {
	reader, header, err := c.GetFile("cover")
	if err != nil {
		return nil, 0
	}
	return reader, header.Size
}

func (facade *AbstractFacade) GetFileNotInvalid(c *web.Controller, key string) (io.ReadCloser, int64) {
	reader, header, err := c.GetFile("cover")
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("get param [%s] as int64 error", key)))
	}
	return reader, header.Size
}

func (facade *AbstractFacade) GetInt64NotInvalid(c *web.Controller, key string) int64 {
	result, err := c.GetInt64(key)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("get param [%s] as int64 error", key)))
	}
	return result
}

func (facade *AbstractFacade) GetStringNotEmpty(c *web.Controller, key string) string {
	result := c.GetString(key, "")
	if result == "" {
		panic(errors.NewError(fmt.Sprintf("get param [%s] as string is empty", key)))
	}
	return result
}

func (facade *AbstractFacade) GetStringAsInt64List(c *web.Controller, key string) []int64 {
	var results []int64
	jsonUtil.UnmarshalJsonString(c.GetString(key, "[]"), &results)
	return results
}

func (facade *AbstractFacade) GetStringAsStringList(c *web.Controller, key string) []string {
	var results []string
	jsonUtil.UnmarshalJsonString(c.GetString(key, "[]"), &results)
	return results
}

func (facade *AbstractFacade) GetStringAsStruct(c *web.Controller, key string, obj interface{}) {
	jsonUtil.UnmarshalJsonString(c.GetString(key, "null"), obj)
}

func (facade *AbstractFacade) GetRestfulParamInt(c *web.Controller, key string) int {
	str := c.Ctx.Input.Param(":page")
	result, err := strconv.Atoi(str)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("get restful param [%s] as int failed", str)))
	}
	return result
}

func (facade *AbstractFacade) GetRestfulParamInt64(c *web.Controller, key string) int64 {
	str := c.Ctx.Input.Param(key)
	result, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("get restful param [%s] as int64 failed", key)))
	}
	return result
}

func (facade *AbstractFacade) GetContext(c *web.Controller) contextDTO.ContextDTO {
	var ctx contextDTO.ContextDTO

	if claim, ok := c.Ctx.Input.GetData(filters.ContextClaim).(string); ok {
		var userClaim userDTO.UserClaimDTO
		jsonUtil.UnmarshalJsonString(claim, &userClaim)
		ctx.UserClaim = userClaim
	}

	return ctx
}
