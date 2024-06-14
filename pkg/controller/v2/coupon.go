package v2

import (
	"net/http"
	. "rentServer/core/controller"
	service2 "rentServer/core/service"
	"rentServer/pkg/model/coupon"
	"rentServer/pkg/service"
	"strconv"
	"time"
)

type CouponController interface {
	FetchCoupon(c *Context)
	UserGetCouponList(c *Context)
}

type couponController struct {
	service.CouponService
}

func NewCouponController() CouponController {
	return &couponController{
		CouponService: service.NewCouponService(),
	}
}

func (_this couponController) FetchCoupon(c *Context) {

	couponID := c.Query("coupon_id")
	//
	err := _this.CouponService.Dispatch(c.GetUint64(Uid), couponID)
	if err != nil {
		return
	}

	c.JSONOK()
}

func (_this couponController) UserGetCouponList(c *Context) {

	statusStr := c.Query("status")
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		return
	}
	var param service2.BaseSearchParam

	err = c.ShouldBindQuery(&param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	list, err := _this.CouponService.UserList(c.GetUint64(Uid), coupon.CouponStatus(status), &param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	for _, rel := range list {
		err = _this.CouponService.FindByID(&rel.Coupon, rel.CouponID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSONOK(list)
}

func (_this couponController) FetchReceiveList(c *Context) {
	var param service2.BaseSearchParam

	param.Page = "1"
	param.Size = "100"

	now := time.Now().UnixMilli()
	list, err := _this.CouponService.FilterList(&param, map[string]any{
		"end_date >= ?": now,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK(list)
}
