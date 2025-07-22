package facade

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/mellolo/common/utils/jsonUtil"
	"media-station/service/biz/bizUser"
	"media-station/storage/db"
)

type UserFacade struct {
	userBizService bizUser.UserBizService
}

func NewUserFacade() *UserFacade {
	return &UserFacade{
		userBizService: bizUser.NewUserBizService(),
	}
}

func (impl *UserFacade) Register(c *web.Controller) {
	body := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	jsonUtil.UnmarshalJsonString(string(c.Ctx.Input.RequestBody), body)
	db.DoTransaction(func(tx orm.TxOrmer) {
		impl.userBizService.Register(body.Username, body.Password, tx)
	})
}

func (impl *UserFacade) Login(c *web.Controller) string {
	body := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	jsonUtil.UnmarshalJsonString(string(c.Ctx.Input.RequestBody), body)
	var token string
	db.DoTransaction(func(tx orm.TxOrmer) {
		token = impl.userBizService.Login(body.Username, body.Password, tx)
	})
	return token
}

func (impl *UserFacade) Logout(token string) {
	db.DoTransaction(func(tx orm.TxOrmer) {
		impl.userBizService.Logout(token)
	})
}
