package v1

import (
	. "rentServer/core/controller"
	core "rentServer/core/service"
	"rentServer/pkg/model/clerk"
	"rentServer/pkg/service"
)

var _ ClerkController = (*clerkController)(nil)

type ClerkController interface {
	GetClerkList(c *Context)
	Save(c *Context)
	Delete(c *Context)
}

type clerkController struct {
	ClerkService service.ClerkService
}

func NewClerkController() ClerkController {
	return &clerkController{
		ClerkService: service.NewClerkService(),
	}
}

func (_this clerkController) GetClerkList(c *Context) {

	var param core.BaseSearchParam

	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSONOK(ErrBadRequest)
		return
	}

	filter := map[string]any{}

	list, err := _this.ClerkService.FilterList(&param, filter)
	if err != nil {
		c.JSONOK(ErrInternalServer)
		return
	}

	c.JSONOK(list)
}

func (_this clerkController) Save(c *Context) {

	var clerkModel clerk.Clerk

	err := c.ShouldBindJSON(&clerkModel)
	if err != nil {
		c.JSONOK(err.Error())
		return
	}

	err = _this.ClerkService.Save(&clerkModel, clerkModel.ID)
	if err != nil {
		c.JSONOK(ErrInternalServer)
		return
	}

	c.JSONOK(clerkModel)
}

func (_this clerkController) Delete(c *Context) {

	var ids []uint64
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		c.JSONOK(ErrBadRequest)
		return
	}

	err = _this.ClerkService.Delete(ids)
	if err != nil {
		c.JSONOK(ErrInternalServer)
		return
	}

	c.JSONOK()

}
