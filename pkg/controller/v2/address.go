package v2

import (
	"net/http"
	. "rentServer/core/controller"
	. "rentServer/core/service"
	"rentServer/pkg/model/address"
	"rentServer/pkg/service"
)

var _ AddressController = (*addressController)(nil)

type AddressController interface {
	GetMyAddress(c *Context)
	Save(c *Context)
	Delete(c *Context)
}

type addressController struct {
	AddressSvc service.AddressService
}

func NewAddressController() AddressController {
	return &addressController{
		AddressSvc: service.NewAddressService(),
	}
}

func (_this addressController) GetMyAddress(c *Context) {

	var param BaseSearchParam

	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	list, err := _this.AddressSvc.FilterList(&param, map[string]any{
		"user_id": c.GetUint64(Uid),
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK(list)
}

func (_this addressController) Save(c *Context) {

	var addressModel address.Address

	err := c.ShouldBindJSON(&addressModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	addressModel.UserID = int64(c.GetUint64(Uid))
	err = _this.AddressSvc.Save(&addressModel, addressModel.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK(addressModel)
}

func (_this addressController) Delete(c *Context) {

	var ids []uint64

	err := c.ShouldBindJSON(&ids)
	if err != nil {
		return
	}

	err = _this.AddressSvc.Delete(ids)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK()
}
