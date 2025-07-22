package bizUser

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/mellolo/common/cache"
	"github.com/mellolo/common/config"
	"github.com/mellolo/common/errors"
	"github.com/mellolo/common/utils/jwtUtil"
	"media-station/models/do/userDO"
	"media-station/storage/db"
	"time"
)

type UserBizService interface {
	Register(username, password string, tx ...orm.TxOrmer)
	Login(username, password string, tx ...orm.TxOrmer) string
	Logout(token string)
}

func NewUserBizService() *UserBizServiceImpl {
	return &UserBizServiceImpl{
		userMapper: db.NewUserMapper(),
	}
}

type UserBizServiceImpl struct {
	userMapper db.UserMapper
}

func (impl UserBizServiceImpl) Register(username, password string, tx ...orm.TxOrmer) {
	_, err := impl.userMapper.Insert(&userDO.UserDO{
		Username: username,
		Password: password,
	}, tx...)
	if err != nil {
		panic(errors.WrapError(err, "user register failed"))
	}
}

func (impl UserBizServiceImpl) Login(username, password string, tx ...orm.TxOrmer) string {
	userDo, err := impl.userMapper.SelectByUsername(username, tx...)
	if err != nil {
		panic(errors.WrapError(err, "user login failed"))
	}
	if password != userDo.Password {
		panic(errors.NewError("password error"))
	}

	secretKey := config.GetConfig("media", "secretKey", "user")
	token, err := jwtUtil.GenerateToken(userDo.Username, secretKey, 1)
	if err != nil {
		panic(errors.WrapError(err, "generate token failed"))
	}
	return token
}

func (impl UserBizServiceImpl) Logout(token string) {
	client, _ := cache.GetCache()
	data, _ := client.Get("userTokenBacklist")
	if blacklist, ok := data.([]string); ok {
		_ = client.Set("userTokenBacklist", append(blacklist, token), time.Hour)
	} else {
		_ = client.Set("userTokenBacklist", []string{token}, time.Hour)
	}
}
