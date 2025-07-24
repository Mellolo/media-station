package bizUser

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/mellolo/common/cache"
	"github.com/mellolo/common/config"
	"github.com/mellolo/common/errors"
	"github.com/mellolo/common/utils/jsonUtil"
	"github.com/mellolo/common/utils/jwtUtil"
	"media-station/models/do/userDO"
	"media-station/models/dto/userDTO"
	"media-station/storage/db"
	"time"
)

type UserBizService interface {
	GetProfile(username string, tx ...orm.TxOrmer) userDTO.UserProfileDTO
	Register(registerDTO userDTO.UserRegisterDTO, tx ...orm.TxOrmer)
	Login(userLoginDTO userDTO.UserLoginDTO, tx ...orm.TxOrmer) string
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

func (impl UserBizServiceImpl) GetProfile(username string, tx ...orm.TxOrmer) userDTO.UserProfileDTO {
	user, err := impl.userMapper.SelectByUsername(username, tx...)
	if err != nil {
		panic(errors.WrapError(err, "get user profile failed"))
	}
	return userDTO.UserProfileDTO{
		Username: user.Username,
		Details: userDTO.UserDetailsDTO{
			VideoIds:   user.Details.VideoIds,
			GalleryIds: user.Details.GalleryIds,
		},
	}
}

func (impl UserBizServiceImpl) Register(registerDTO userDTO.UserRegisterDTO, tx ...orm.TxOrmer) {
	_, err := impl.userMapper.Insert(&userDO.UserDO{
		Username: registerDTO.Username,
		Password: registerDTO.Password,
	}, tx...)
	if err != nil {
		panic(errors.WrapError(err, "user register failed"))
	}
}

func (impl UserBizServiceImpl) Login(userLoginDTO userDTO.UserLoginDTO, tx ...orm.TxOrmer) string {
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

func (impl UserBizServiceImpl) Logout(token string) {
	client, _ := cache.GetCache()
	data, _ := client.Get("userTokenBacklist")
	if blacklist, ok := data.([]string); ok {
		_ = client.Set("userTokenBacklist", append(blacklist, token), time.Hour)
	} else {
		_ = client.Set("userTokenBacklist", []string{token}, time.Hour)
	}
}
