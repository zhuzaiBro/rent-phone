package admin

import (
	"errors"
	"gorm.io/gorm"
	. "rentServer/core/service"
	admin_dao "rentServer/pkg/dao/v1/admin"
	"rentServer/pkg/logger"
	"rentServer/pkg/middle_ware"
	"rentServer/pkg/model/admin"
)

var _ AdminService = (*adminServiceIMPL)(nil)

type AdminService interface {
	Service[admin.Admin]
	Login(login admin.AdminLogin) (token string, err error)
}

type adminServiceIMPL struct {
	ServiceIMPL[admin.Admin]
	AdminDao admin_dao.AdminDao
}

func NewAdminService() AdminService {
	adminDao := admin_dao.NewAdminDao()
	return &adminServiceIMPL{
		ServiceIMPL: ServiceIMPL[admin.Admin]{
			BaseDao: adminDao,
		},
		AdminDao: adminDao,
	}
}

func (_this adminServiceIMPL) Login(param admin.AdminLogin) (token string, err error) {

	var adminModel admin.Admin
	filter := map[string]any{
		"username": param.UserName,
	}
	err = _this.FilterFindOne(&adminModel, filter)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error(err.Error())
		return "", err
	}

	if adminModel.Password != param.Password {
		return "", errors.New("账号或密码错误")
	}

	token, err = middle_ware.GenerateAdminToken(adminModel.ID)
	if err != nil {
		logger.Error("GenerateAdminToken", err.Error())
		return "", err
	}

	return
}
