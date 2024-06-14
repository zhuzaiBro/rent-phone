package v1

import (
	"net/http"
	. "rentServer/core/controller"
	. "rentServer/core/service"
	"rentServer/pkg/model/album"
	"rentServer/pkg/service"
)

type AlbumController interface {
	List(ctx *Context)
	Save(ctx *Context)
	Delete(ctx *Context)
}

type albumController struct {
	service.AlbumService
}

func NewAlbumController() AlbumController {
	return &albumController{}
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

func (_this albumController) Save(c *Context) {

	var albumModel album.Album

	err := c.ShouldBindJSON(&albumModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = _this.AlbumService.Save(&albumModel, albumModel.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK(albumModel)
}

func (_this albumController) Delete(c *Context) {
	var ids []uint64

	err := c.ShouldBindJSON(&ids)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = _this.AlbumService.Delete(ids)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK()
}
