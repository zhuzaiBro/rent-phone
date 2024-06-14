package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/user_coupon_rel"
)

type UserCouponRelDao interface {
	Dao[user_coupon_rel.UserCouponRel]
}

type userCouponRelDao struct {
	DaoIMPL[user_coupon_rel.UserCouponRel]
}

func NewUserCouponRelDao() UserCouponRelDao {
	return &userCouponRelDao{
		DaoIMPL: DaoIMPL[user_coupon_rel.UserCouponRel]{
			Db:    db.DB,
			Model: &user_coupon_rel.UserCouponRel{},
		},
	}
}
