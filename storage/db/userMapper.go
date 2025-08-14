package db

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/mellolo/common/utils/jsonUtil"
	"media-station/models/dao/userDAO"
	"media-station/models/do/userDO"
)

type UserMapper interface {
	Insert(user userDO.UserDO, tx ...orm.TxOrmer) (int64, error)
	SelectByUsername(username string, tx ...orm.TxOrmer) (userDO.UserDO, error)
	Update(user userDO.UserDO, tx ...orm.TxOrmer) error
	DeleteByUsername(username string, tx ...orm.TxOrmer) error
}

type UserMapperImpl struct{}

func NewUserMapper() *UserMapperImpl {
	return &UserMapperImpl{}
}

func (impl *UserMapperImpl) Insert(user userDO.UserDO, tx ...orm.TxOrmer) (int64, error) {
	executor := getQueryExecutor(tx...)

	// todo 有更复杂处理
	record := userDAO.UserRecord{
		Username:    user.Username,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		WechatId:    user.WechatId,
		Details:     jsonUtil.GetJsonString(user.Art),
	}

	return executor.Insert(&record)
}

func (impl *UserMapperImpl) SelectByUsername(username string, tx ...orm.TxOrmer) (userDO.UserDO, error) {
	executor := getQueryExecutor(tx...)
	var record userDAO.UserRecord
	err := executor.QueryTable(userDAO.TableUser).Filter("username", username).One(&record)
	if err != nil {
		return userDO.UserDO{}, err
	}

	// todo 有更复杂处理
	var userArt userDO.UserArt
	jsonUtil.UnmarshalJsonString(record.Details, &userArt)
	do := userDO.UserDO{
		Username:    record.Username,
		Password:    record.Password,
		PhoneNumber: record.PhoneNumber,
		WechatId:    record.WechatId,
		Art:         userArt,
	}
	return do, nil
}

func (impl *UserMapperImpl) Update(user userDO.UserDO, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)

	// todo 有更复杂处理
	_, err := executor.QueryTable(userDAO.TableUser).Filter("username", user.Username).Update(orm.Params{
		"password":     user.Password,
		"phone_number": user.PhoneNumber,
		"wechat_id":    user.WechatId,
		"details":      jsonUtil.GetJsonString(user.Art),
	})
	return err
}

func (impl *UserMapperImpl) DeleteByUsername(username string, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)
	_, err := executor.QueryTable(userDAO.TableUser).Filter("username", username).Delete()
	return err
}
