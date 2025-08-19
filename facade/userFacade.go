package facade

import (
	"github.com/Mellolo/common/cache"
	"github.com/Mellolo/common/config"
	"github.com/Mellolo/common/utils/jsonUtil"
	"github.com/Mellolo/common/utils/jwtUtil"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"k8s.io/apimachinery/pkg/util/sets"
	"media-station/models/dto/contextDTO"
	"media-station/models/dto/userDTO"
	"media-station/models/vo/userVO"
	"media-station/service/biz/bizUser"
	"media-station/storage/db"
)

type UserFacade struct {
	AbstractFacade
	userBizService bizUser.UserBizService
}

func NewUserFacade() *UserFacade {
	return &UserFacade{
		userBizService: bizUser.NewUserBizService(),
	}
}

func (impl *UserFacade) Register(c *web.Controller) {
	// 上下文
	ctx := impl.GetContext(c)

	var dto userDTO.UserRegisterDTO
	jsonUtil.UnmarshalJsonString(string(c.Ctx.Input.RequestBody), &dto)

	db.DoTransaction(func(tx orm.TxOrmer) {
		impl.userBizService.Register(ctx, dto, tx)
	})
}

func (impl *UserFacade) Login(c *web.Controller) string {
	// 上下文
	ctx := impl.GetContext(c)

	var dto userDTO.UserLoginDTO
	jsonUtil.UnmarshalJsonString(string(c.Ctx.Input.RequestBody), &dto)

	var token string
	db.DoTransaction(func(tx orm.TxOrmer) {
		token = impl.userBizService.Login(ctx, dto, tx)
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
	if blacklistStr, ok := data.(string); ok {
		var blacklist []string
		jsonUtil.UnmarshalJsonString(blacklistStr, &blacklist)
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
		impl.userBizService.Logout(contextDTO.ContextDTO{}, token)
	})
}

func (impl *UserFacade) GetProfile(c *web.Controller) userVO.UserProfileVO {
	// 上下文
	ctx := impl.GetContext(c)

	var profile userDTO.UserProfileDTO
	db.DoTransaction(func(tx orm.TxOrmer) {
		profile = impl.userBizService.GetProfile(ctx, ctx.UserClaim.Username, tx)
	})

	return userVO.UserProfileVO{
		Username: profile.Username,
	}
}
