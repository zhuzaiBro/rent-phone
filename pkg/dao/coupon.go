package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/coupon"
)

func NewCouponDao() CouponDao {
	return &couponDao{
		DaoIMPL: DaoIMPL[coupon.Coupon]{
			Db:    db.DB,
			Model: &coupon.Coupon{},
		},
	}
}

type CouponDao interface {
	Dao[coupon.Coupon]
}

type couponDao struct {
	DaoIMPL[coupon.Coupon]
}
