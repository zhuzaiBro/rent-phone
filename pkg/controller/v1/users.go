package v1

import (
	. "rentServer/core/controller"
	"rentServer/pkg/model/admin"
	admin_svc "rentServer/pkg/service/v1/admin"
)

var _ AdminController = (*adminController)(nil)

type AdminController interface {
	Login(c *Context)
	GetUserInfo(c *Context)
}

type adminController struct {
	AdminService admin_svc.AdminService
}

func NewAdminController() AdminController {
	return &adminController{
		AdminService: admin_svc.NewAdminService(),
	}
}

func (_this adminController) Login(c *Context) {
	var adminLogin admin.AdminLogin

	err := c.ShouldBindJSON(&adminLogin)
	if err != nil {
		c.JSONOK(ErrBadRequest)
		return
	}

	token, err := _this.AdminService.Login(adminLogin)
	if err != nil {
		c.JSONOK(ErrInternalServer)
		return
	}

	c.JSONOK(map[string]string{"token": token})
}

func (_this adminController) GetUserInfo(c *Context) {
	var adminModel admin.Admin
	filter := map[string]any{
		"is_del": "0",
		"id":     c.GetUint64(AdminID),
	}

	err := _this.AdminService.FilterFindOne(&adminModel, filter)
	if err != nil {
		c.JSONOK(ErrInternalServer)
		return
	}

	c.JSONOK(adminModel)
}
