package bizUser

import (
	"github.com/Mellolo/common/cache"
	"github.com/Mellolo/common/config"
	"github.com/Mellolo/common/errors"
	"github.com/Mellolo/common/utils/jsonUtil"
	"github.com/Mellolo/common/utils/jwtUtil"
	"github.com/beego/beego/v2/client/orm"
	"media-station/models/do/userDO"
	"media-station/models/dto/contextDTO"
	"media-station/models/dto/userDTO"
	"media-station/storage/db"
	"time"
)

type UserBizService interface {
	GetProfile(ctx contextDTO.ContextDTO, username string, tx ...orm.TxOrmer) userDTO.UserProfileDTO
	Register(ctx contextDTO.ContextDTO, registerDTO userDTO.UserRegisterDTO, tx ...orm.TxOrmer)
	Login(ctx contextDTO.ContextDTO, userLoginDTO userDTO.UserLoginDTO, tx ...orm.TxOrmer) string
	Logout(ctx contextDTO.ContextDTO, token string)
}

func NewUserBizService() *UserBizServiceImpl {
	return &UserBizServiceImpl{
		userMapper: db.NewUserMapper(),
	}
}

type UserBizServiceImpl struct {
	userMapper db.UserMapper
}

func (impl UserBizServiceImpl) GetProfile(ctx contextDTO.ContextDTO, username string, tx ...orm.TxOrmer) userDTO.UserProfileDTO {
	user, err := impl.userMapper.SelectByUsername(username, tx...)
	if err != nil {
		panic(errors.WrapError(err, "get user profile failed"))
	}
	return userDTO.UserProfileDTO{
		Username: user.Username,
	}
}

func (impl UserBizServiceImpl) Register(ctx contextDTO.ContextDTO, registerDTO userDTO.UserRegisterDTO, tx ...orm.TxOrmer) {
	_, err := impl.userMapper.Insert(userDO.UserDO{
		Username: registerDTO.Username,
		Password: registerDTO.Password,
	}, tx...)
	if err != nil {
		panic(errors.WrapError(err, "user register failed"))
	}
}

func (impl UserBizServiceImpl) Login(ctx contextDTO.ContextDTO, userLoginDTO userDTO.UserLoginDTO, tx ...orm.TxOrmer) string {
	userDo, err := impl.userMapper.SelectByUsername(userLoginDTO.Username, tx...)
	if err != nil {
		panic(errors.WrapError(err, "user login failed"))
	}
	if userLoginDTO.Password != userDo.Password {
		panic(errors.NewError("password error"))
	}

	secretKey := config.GetConfig("media", "secretKey", "user")

	token, err := jwtUtil.GenerateToken(
		jsonUtil.GetJsonString(userDTO.UserClaimDTO{Username: userDo.Username}),
		secretKey, 1)
	if err != nil {
		panic(errors.WrapError(err, "generate token failed"))
	}
	return token
}

func (impl UserBizServiceImpl) Logout(ctx contextDTO.ContextDTO, token string) {
	client, _ := cache.GetCache()
	data, _ := client.Get("userTokenBacklist")
	if blacklistStr, ok := data.(string); ok {
		var blacklist []string
		jsonUtil.UnmarshalJsonString(blacklistStr, &blacklist)
		_ = client.Set("userTokenBacklist", jsonUtil.GetJsonString(append(blacklist, token)), time.Hour)
	} else {
		_ = client.Set("userTokenBacklist", jsonUtil.GetJsonString([]string{token}), time.Hour)
	}
}
