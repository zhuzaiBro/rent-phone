package v1

import (
	"net/http"
	. "rentServer/core/controller"
	. "rentServer/core/service"
	"rentServer/pkg/model/coupon"
	"rentServer/pkg/service"
	"strconv"
)

type CouponController interface {
	SaveCoupon(c *Context)
	DeleteCoupon(c *Context)
	GetCouponRecord(c *Context)
}

type couponController struct {
	service.CouponService
}

func NewCouponController() CouponController {
	return &couponController{}
}

func (_this couponController) SaveCoupon(c *Context) {

	var couponModel coupon.Coupon

	err := c.ShouldBindJSON(&couponModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = _this.CouponService.Save(&couponModel, couponModel.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK(couponModel)

}

func (_this couponController) DeleteCoupon(c *Context) {
	var ids []uint64

	err := c.ShouldBindJSON(&ids)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = _this.CouponService.Delete(ids)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK()

}

// GetCouponRecord 管理端查询优惠券领取发放情况
func (_this couponController) GetCouponRecord(c *Context) {

	statusStr := c.Query("status")

	status, err := strconv.Atoi(statusStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	var param BaseSearchParam

	err = c.ShouldBindQuery(&param)
	if err != nil {
		return
	}

	list, err := _this.CouponService.FilterList(&param, map[string]any{
		"status": status,
	})
	if err != nil {
		return
	}

	c.JSONOK(list)
}
