package v2

import (
	"net/http"
	. "rentServer/core/controller"
	. "rentServer/core/service"
	"rentServer/pkg/service"
)

type AlbumController interface {
	List(c *Context)
}

type albumController struct {
	service.AlbumService
}

func NewAlbumController() AlbumController {
	return &albumController{
		AlbumService: service.NewAlbumService(),
	}
}

func (_this albumController) List(c *Context) {

	var param BaseSearchParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	list, err := _this.AlbumService.List(&param)
	if err != nil {
		return
	}

	c.JSONOK(list)
}
