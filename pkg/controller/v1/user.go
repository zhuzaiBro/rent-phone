package v1

import (
	"net/http"
	. "rentServer/core/controller"
	. "rentServer/core/service"
	"rentServer/pkg/model/user"
	"rentServer/pkg/service/v1"
)

type UserController interface {
	List(c *Context)
	Update(c *Context)
}

type userController struct {
	UserSvc v1.UserService
}

func NewUserController() UserController {
	return &userController{
		UserSvc: v1.NewUserService(),
	}
}

func (_this userController) List(c *Context) {

	var (
		param BaseSearchParam
	)

	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	list, err := _this.UserSvc.List(&param, map[string]any{})
	if err != nil {
		return
	}

	c.JSONOK(list)
}

func (_this userController) Update(c *Context) {
	var userModel user.User

	err := c.ShouldBindJSON(&userModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = _this.UserSvc.Update(&userModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK(userModel)

}
