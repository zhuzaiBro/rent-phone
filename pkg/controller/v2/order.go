package v2

import (
	"encoding/json"

	"net/http"
	. "rentServer/core/controller"
	service2 "rentServer/core/service"
	"rentServer/core/utils"
	"rentServer/pkg/logger"
	"rentServer/pkg/model/order"
	"rentServer/pkg/model/reservation"
	"rentServer/pkg/model/user"
	"rentServer/pkg/service"
	v1 "rentServer/pkg/service/v1"
	"strconv"
)

var _ OrderController = (*orderController)(nil)

type OrderController interface {
	GenerateOrder(c *Context)
	WxPayOrder(c *Context)
	GetMyOrder(c *Context)
	Detail(c *Context)
	ComputeOrder(c *Context)
	GetOrderCount(c *Context)
}

type orderController struct {
	OrderService service.OrderService
	UserService  service.UserService
	AddressSvc   service.AddressService
	StoreSvc     v1.StoreService
}

func NewOrderController() OrderController {
	return &orderController{
		OrderService: service.NewOrderService(),
		UserService:  service.NewUserService(),
		AddressSvc:   service.NewAddressService(),
		StoreSvc:     v1.NewStoreService(),
	}
}

// GetMyOrder 获取我的订单
func (_this orderController) GetMyOrder(c *Context) {

	var (
		_param order.OrderParam
		param  service2.BaseSearchParam
		list   []*order.Order
	)

	err := c.ShouldBindQuery(&_param)
	if err != nil {
		return
	}

	err = c.ShouldBindQuery(&param)
	if err != nil {
		return
	}

	_param.UserID = c.GetUint64(Uid)

	logger.Info(c.Query("all"))

	if c.Query("all") == "true" {
		list, err = _this.OrderService.FilterList(&param, map[string]any{
			"user_id": c.GetUint64(Uid),
			"type":    _param.Type,
		})
	} else {
		list, err = _this.OrderService.FilterList(&param, utils.FlattenStruct(_param))
	}
	if err != nil {
		logger.Error("查询出错了")
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	for _, orderItem := range list {
		err := json.Unmarshal([]byte(orderItem.Products), &orderItem.ProductList)
		if err != nil {
			logger.Error("list里面有什么？", orderItem.TradeNo, orderItem.Products)
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSONOK(list)
}

// GenerateOrder 生成订单
func (_this orderController) GenerateOrder(c *Context) {
	orderType := c.Query("type")

	var orderModel order.Order

	orderModel.UserID = strconv.FormatUint(c.GetUint64(Uid), 10)
	if orderType == service.GOOD {
		var cartIds []uint64
		err := c.ShouldBindJSON(&cartIds)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
		err = _this.OrderService.GenerateGoodOrder(c.GetUint64(Uid), cartIds, &orderModel)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

	} else if orderType == service.SERVICE {
		var rentModel reservation.Reservation
		err := c.ShouldBindJSON(&rentModel)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
		err = _this.OrderService.GenerateServiceOrder(&orderModel, &rentModel)
		if err != nil {
			c.JSONOK(err.Error())
			return
		}
	} else if orderType == service.DEPOSIT {
		err := _this.OrderService.GenerateDepositOrder(&orderModel, c.Query("deposit_num"), c.GetUint64(Uid))
		if err != nil {
			c.JSONOK(err.Error())
			return
		}
	}

	c.JSONOK(orderModel)
}

// WxPayOrder 微信支付订单
func (_this orderController) WxPayOrder(c *Context) {

	tradeNo := c.Query("trade_no")
	payType := c.Query("pay_type")
	couponRelID := c.Query("coupon_rel_id")

	if len(tradeNo) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, "请传入订单流水号ID！")
		return
	}

	_payType, err := strconv.Atoi(payType)

	var (
		userModel user.User
	)
	filter := map[string]any{
		"id": c.GetUint64(Uid),
	}
	err = _this.UserService.FilterFindOne(&userModel, filter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if int(order.DATE) == _payType {
		// 直接修改到已支付状态
		var orderModel order.Order
		err = _this.OrderService.HandleDateOrder(tradeNo, &orderModel, c.GetUint64(Uid))
		if err != nil {
			return
		}
		c.JSONOK()
		return
	}

	payOrder, err := _this.OrderService.GenerateWxMpPayOrder(tradeNo, userModel.WxOpenID, c.GetUint64(Uid), couponRelID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONOK(payOrder)
}

// Detail 查询订单信息
func (_this orderController) Detail(c *Context) {

	var orderModel order.Order
	err := _this.OrderService.FilterFindOne(&orderModel, map[string]any{
		"id": c.Query("id"),
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if len(orderModel.Addr) > 0 {
		err = json.Unmarshal([]byte(orderModel.Addr), &orderModel.Address)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	err = json.Unmarshal([]byte(orderModel.Products), &orderModel.ProductList)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = _this.StoreSvc.FindByID(&orderModel.Store, 1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK(orderModel)
}

// ComputeOrder 修改未提交订单信息
func (_this orderController) ComputeOrder(c *Context) {
	var orderModel order.Order

	err := c.ShouldBindJSON(&orderModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = _this.OrderService.ComputeOrder(&orderModel, c.GetUint64(Uid))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK(orderModel)
}

func (_this orderController) GetOrderCount(c *Context) {
	var resp order.OrderCount

	filter := map[string]any{
		"user_id": c.GetUint64(Uid),
		"is_del":  0,
		"type":    c.Query("type"),
	}
	total, err := _this.OrderService.GetMyOrderCount(filter)
	resp.Total = total
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	filter["status"] = order.PENDING
	resp.UnPayed, err = _this.OrderService.GetMyOrderCount(filter)

	filter["status"] = order.ING
	resp.ING, err = _this.OrderService.GetMyOrderCount(filter)

	filter["status"] = order.WAIT_COMMENT
	resp.UnComment, err = _this.OrderService.GetMyOrderCount(filter)

	filter["status"] = order.FINISHED
	resp.Finish, err = _this.OrderService.GetMyOrderCount(filter)

	c.JSONOK(resp)
}
