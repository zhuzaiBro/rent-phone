package v2

import (
	"encoding/json"
	"net/http"
	. "rentServer/core/controller"
	core "rentServer/core/service"
	"rentServer/pkg/model/cart"
	"rentServer/pkg/service"
)

var _ CartController = (*cartController)(nil)

type CartController interface {
	Add(c *Context)
	List(c *Context)
	Delete(c *Context)
}

type cartController struct {
	CartService service.CartService
}

func NewCartController() CartController {
	return &cartController{
		CartService: service.NewCartService(),
	}
}

func (_this cartController) Add(c *Context) {

	kind := c.Query("type")
	switch kind {
	case service.GOOD:
		{
			var addReq cart.AddRequest
			err := c.ShouldBindJSON(&addReq)
			if err != nil {
				return
			}

			cartList, err := _this.CartService.AddGoods2Cart(c.GetUint64(Uid), addReq)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
				return
			}

			c.JSONOK(cartList)
		}
	case service.SERVICE:

	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, "没有给订单类型")
		return
	}

}

func (_this cartController) List(c *Context) {
	var param core.BaseSearchParam

	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	filter := map[string]any{
		"user_id": c.GetUint64(Uid),
	}

	list, err := _this.CartService.FilterList(&param, filter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	for _, cartItem := range list {
		err = json.Unmarshal([]byte(cartItem.Product), &cartItem.P)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSONOK(list)
}

func (_this cartController) Delete(c *Context) {

	var ids []uint64

	err := c.ShouldBindJSON(&ids)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = _this.CartService.Delete(ids)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK()
}
