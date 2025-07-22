package facade

import (
	"github.com/beego/beego/v2/client/orm"
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

func (impl *UserFacade) Register(username, password string) {
	db.DoTransaction(func(tx orm.TxOrmer) {
		impl.userBizService.Register(username, password, tx)
	})
}

func (impl *UserFacade) Login(username, password string) string {
	var token string
	db.DoTransaction(func(tx orm.TxOrmer) {
		token = impl.userBizService.Login(username, password, tx)
	})
	return token
}

func (impl *UserFacade) Logout(token string) {
	db.DoTransaction(func(tx orm.TxOrmer) {
		impl.userBizService.Logout(token)
	})
}
