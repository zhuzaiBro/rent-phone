package v1

import (
	"net/http"
	. "rentServer/core/controller"
	. "rentServer/core/service"
	. "rentServer/pkg/model/customized_project"
	"rentServer/pkg/service"
)

var _ CustomProjController = (*customProjectController)(nil)

type CustomProjController interface {
	Save(c *Context)
	Get(c *Context)
	Delete(c *Context)
}

type customProjectController struct {
	CustomProjService service.CustomProjService
}

func NewCustomProjController() CustomProjController {
	return &customProjectController{
		CustomProjService: service.NewCustomProjService(),
	}
}

func (_this customProjectController) Save(c *Context) {

	var customProjModel CProject

	err := c.ShouldBindJSON(&customProjModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = _this.CustomProjService.Save(&customProjModel, customProjModel.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSONOK(customProjModel)
}

func (_this customProjectController) Get(c *Context) {

	var param BaseSearchParam

	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	list, err := _this.CustomProjService.FilterList(&param, map[string]any{})
	if err != nil {
		return
	}

	c.JSONOK(list)
}

func (_this customProjectController) Delete(c *Context) {

	var ids []uint64

	err := c.ShouldBindJSON(&ids)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = _this.CustomProjService.Delete(ids)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK()
}
