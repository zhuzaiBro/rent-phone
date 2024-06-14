package v1

import (
	"net/http"
	. "rentServer/core/controller"
	. "rentServer/core/service"
	"rentServer/pkg/model/store"
	"rentServer/pkg/service"
	v1 "rentServer/pkg/service/v1"
	"strconv"
)

var _ StoreController = (*storeController)(nil)

type StoreController interface {
	StoreInfo(c *Context)
	Save(c *Context)
}

type storeController struct {
	StoreService      v1.StoreService
	rentService       service.ReservationService
	CustomProjService service.CustomProjService
}

func NewStoreController() StoreController {
	return &storeController{
		StoreService:      v1.NewStoreService(),
		rentService:       service.NewrentService(),
		CustomProjService: service.NewCustomProjService(),
	}
}

func (_this storeController) Save(c *Context) {

	var storeModel store.Store
	err := c.ShouldBindJSON(&storeModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = _this.StoreService.Save(&storeModel, storeModel.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK(storeModel)
}

func (_this storeController) StoreInfo(c *Context) {

	var storeModel store.Store

	// 先默认是1号店
	err := _this.StoreService.FindByID(&storeModel, 1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	timestamp := c.Query("timestamp")

	if timestamp != "" {
		_timestamp, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return
		}

		storeModel.ServiceTime = _this.rentService.GetServiceTime(_timestamp)
	}

	param := BaseSearchParam{
		Page: "1",
		Size: "100000",
	}

	storeModel.Projects, err = _this.CustomProjService.FilterList(&param, map[string]any{})
	if err != nil {
		return
	}

	c.JSONOK(storeModel)
}
