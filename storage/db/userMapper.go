package db

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/mellolo/common/utils/jsonUtil"
	"media-station/models/dao/daoCommon"
	"media-station/models/dao/userDAO"
	"media-station/models/do/userDO"
)

type UserMapper interface {
	Insert(user *userDO.UserDO, tx ...orm.TxOrmer) (int64, error)
	SelectById(id int64, tx ...orm.TxOrmer) (*userDO.UserDO, error)
	Update(id int64, user *userDO.UserDO, tx ...orm.TxOrmer) error
	DeleteById(id int64, tx ...orm.TxOrmer) error
}

type UserMapperImpl struct{}

func NewUserMapper() *UserMapperImpl {
	return &UserMapperImpl{}
}

func (impl *UserMapperImpl) Insert(user *userDO.UserDO, tx ...orm.TxOrmer) (int64, error) {
	executor := getQueryExecutor(tx...)

	record := userDAO.UserRecord{
		Username:    user.Username,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		WechatId:    user.WechatId,
		Details:     jsonUtil.GetJsonString(user.Details),
	}

	return executor.Insert(&record)
}

func (impl *UserMapperImpl) SelectById(id int64, tx ...orm.TxOrmer) (*userDO.UserDO, error) {
	executor := getQueryExecutor(tx...)
	record := userDAO.UserRecord{CommonColumn: daoCommon.CommonColumn{Id: id}}
	err := executor.Read(&record)
	if err != nil {
		return nil, err
	}

	var userDetails userDO.UserDetails
	jsonUtil.UnmarshalJsonString(record.Details, &userDetails)
	do := &userDO.UserDO{
		Id:          record.Id,
		Username:    record.Username,
		Password:    record.Password,
		PhoneNumber: record.PhoneNumber,
		WechatId:    record.WechatId,
		Details:     userDetails,
	}
	return do, nil
}

func (impl *UserMapperImpl) Update(id int64, user *userDO.UserDO, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)

	record := userDAO.UserRecord{
		CommonColumn: daoCommon.CommonColumn{
			Id: id,
		},
		Username:    user.Username,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		WechatId:    user.WechatId,
		Details:     jsonUtil.GetJsonString(user.Details),
	}
	_, err := executor.Update(&record)
	return err
}

func (impl *UserMapperImpl) DeleteById(id int64, tx ...orm.TxOrmer) error {
	executor := getQueryExecutor(tx...)
	record := &userDAO.UserRecord{CommonColumn: daoCommon.CommonColumn{Id: id}}
	_, err := executor.Delete(record)
	return err
}
