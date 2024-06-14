package v2

import (
	. "rentServer/core/controller"
	"rentServer/pkg/logger"
	"rentServer/pkg/model/wx_api"
	"rentServer/pkg/service"
)

type PayController interface {
	WxPayCallback(c *Context)
}

type payController struct {
	WxPayService service.WxPayService
	OrderService service.OrderService
}

func NewPayController() PayController {
	return &payController{
		WxPayService: service.NewWxPayService(),
		OrderService: service.NewOrderService(),
	}
}

func (_this payController) WxPayCallback(c *Context) {
	var (
		wxCallback wx_api.WxPayCallback
		notify     wx_api.WechatPayNotify
	)

	err := c.ShouldBind(&notify)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = _this.WxPayService.DecodeNotify(&notify, &wxCallback)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info(wxCallback)

	err = _this.OrderService.PayOrder(&wxCallback)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
