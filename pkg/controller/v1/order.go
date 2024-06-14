package v1

import (
	"encoding/json"
	"net/http"
	. "rentServer/core/controller"
	core "rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/logger"
	"rentServer/pkg/model/order"
	"rentServer/pkg/service"
)

type OrderController interface {
	GetOrders(c *Context)
	Detail(c *Context)
	Update(c *Context)
}

type orderController struct {
	OrderSvc   service.OrderService
	AddressDao dao.AddressDao
}

func NewOrderController() OrderController {
	return &orderController{
		OrderSvc: service.NewOrderService(),
	}
}

func (_this orderController) GetOrders(c *Context) {
	var param core.BaseSearchParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	filter := map[string]any{
		"status": c.Query("status"),
	}
	kind := c.Query("type")
	if kind != "" {
		filter["type"] = kind
	}

	list, err := _this.OrderSvc.FilterList(&param, filter)
	if err != nil {
		return
	}

	for _, orderItem := range list {
		// 商品需要处理地址
		if orderItem.Type == service.GOOD {
			err = json.Unmarshal([]byte(orderItem.Addr), &orderItem.Address)
			if err != nil {
				logger.Error(err.Error())
				//c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
				//return
			}
		}

		err = json.Unmarshal([]byte(orderItem.Products), &orderItem.ProductList)
		if err != nil {
			logger.Error(err.Error())
			//c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			//return
		}

	}
	c.JSONOK(list)
}

func (_this orderController) Detail(c *Context) {

	filter := map[string]any{
		"id": c.Query("id"),
	}
	var orderModel order.Order
	err := _this.OrderSvc.FilterFindOne(&orderModel, filter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = json.Unmarshal([]byte(orderModel.Products), &orderModel.ProductList)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK(orderModel)
}

func (_this orderController) Update(c *Context) {
	var orderModel order.Order

	err := c.ShouldBindJSON(&orderModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = _this.OrderSvc.Update(&orderModel, orderModel.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK(orderModel)
}
