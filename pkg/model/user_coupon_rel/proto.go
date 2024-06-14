package user_coupon_rel

import (
	"github.com/shopspring/decimal"
	. "rentServer/core/model"
	"rentServer/pkg/model/coupon"
)

type UserCouponRel struct {
	BaseModel
	Category       string          `json:"category,omitempty" gorm:"column:category"`
	CouponID       string          `json:"coupon_id" gorm:"column:coupon_id"`
	Discount       decimal.Decimal `gorm:"column:discount;comment:'优惠券折扣金额'" json:"discount"`
	ExpirationDate int64           `gorm:"column:expiration_date;comment:'过期时间'" json:"expiration_date"`
	IsActive       bool            `gorm:"column:active" json:"active"`
	CouponCode     string          `gorm:"column:coupon_code;type:varchar(255);comment:'优惠券代码'" json:"coupon_code"`
	// 产品名称
	Product string        `gorm:"column:product;type:varchar(255)" json:"product,omitempty"`
	UserID  string        `gorm:"column:user_id;type:varchar(255)" json:"user_id"`
	Coupon  coupon.Coupon `gorm:"-" json:"coupon"`
}

func (_this UserCouponRel) TableName() string {
	return "user_coupon_rel"
}
