package facade

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/mellolo/common/cache"
	"github.com/mellolo/common/config"
	"github.com/mellolo/common/utils/jsonUtil"
	"github.com/mellolo/common/utils/jwtUtil"
	"k8s.io/apimachinery/pkg/util/sets"
	"media-station/controllers/filters"
	"media-station/models/dto/userDTO"
	"media-station/models/vo/userVO"
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
	var dto userDTO.UserRegisterDTO
	jsonUtil.UnmarshalJsonString(string(c.Ctx.Input.RequestBody), dto)
	db.DoTransaction(func(tx orm.TxOrmer) {
		impl.userBizService.Register(dto, tx)
	})
}

func (impl *UserFacade) Login(c *web.Controller) string {
	var dto userDTO.UserLoginDTO
	jsonUtil.UnmarshalJsonString(string(c.Ctx.Input.RequestBody), dto)
	var token string
	db.DoTransaction(func(tx orm.TxOrmer) {
		token = impl.userBizService.Login(dto, tx)
	})
	return token
}

func (impl *UserFacade) LoginStatus(c *web.Controller) (userVO.UserStatusProfileVO, bool) {
	tokenStr := c.Ctx.Request.Header.Get("Authorization")
	if tokenStr == "" {
		return userVO.UserStatusProfileVO{}, false
	}

	client, _ := cache.GetCache()
	data, _ := client.Get("userTokenBacklist")
	if blacklist, ok := data.([]string); ok {
		if sets.NewString(blacklist...).Has(tokenStr) {
			return userVO.UserStatusProfileVO{}, false
		}
	}

	secretKey := config.GetConfig("media", "secretKey", "user")
	claim, err := jwtUtil.ParseToken(tokenStr, secretKey)
	if err != nil {
		return userVO.UserStatusProfileVO{}, false
	}

	var userClaim userDTO.UserClaimDTO
	jsonUtil.UnmarshalJsonString(claim, &userClaim)
	username := userClaim.Username

	return userVO.UserStatusProfileVO{
		Username: username,
	}, true
}

func (impl *UserFacade) Logout(token string) {
	db.DoTransaction(func(tx orm.TxOrmer) {
		impl.userBizService.Logout(token)
	})
}

func (impl *UserFacade) GetProfile(c *web.Controller) userVO.UserProfileVO {
	username := ""
	if claim, ok := c.Ctx.Input.GetData(filters.ContextClaim).(string); ok {
		var userClaim userDTO.UserClaimDTO
		jsonUtil.UnmarshalJsonString(claim, &userClaim)
		username = userClaim.Username
	}
	var profile userDTO.UserProfileDTO
	db.DoTransaction(func(tx orm.TxOrmer) {
		profile = impl.userBizService.GetProfile(username, tx)
	})

	return userVO.UserProfileVO{
		Username: profile.Username,
	}
}
