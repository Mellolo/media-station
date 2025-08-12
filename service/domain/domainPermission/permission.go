package domainPermission

import (
	"media-station/enum"
	"media-station/models/do/userDO"
)

type PermissionDomainService interface {
	IsAccessAllowed(user userDO.UserDO, uploader, permissionLevel string) bool
}

func NewPermissionDomainService() *PermissionDomainServiceImpl {
	return &PermissionDomainServiceImpl{}
}

type PermissionDomainServiceImpl struct {
}

func (impl PermissionDomainServiceImpl) IsAccessAllowed(user userDO.UserDO, uploader, permissionLevel string) bool {
	if user.Username != "" && user.Username == uploader && permissionLevel != enum.PermissionForbidden {
		return true
	}

	switch permissionLevel {
	case enum.PermissionForbidden:
		return false
	case enum.PermissionPrivate:
		if user.Username == uploader {
			return true
		} else {
			return false
		}
	case enum.PermissionLogin:
		if user.Username != "" {
			return true
		} else {
			return false
		}
	case enum.PermissionVIP:
		if user.Username != "" {
			return true
		} else {
			return false
		}
	default:
		return true
	}
}
