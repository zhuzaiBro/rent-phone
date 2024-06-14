package coupon_goods

import . "rentServer/core/model"

type CouponGoods struct {
	BaseModel
	ProductID string `gorm:"column:product_id;type:varchar(255)" json:"product_id"`
	CouponID  string `gorm:"column:coupon_id;type:varchar(255)" json:"coupon_id"`
}

func (_this CouponGoods) TableName() string {
	return "coupon_goods"
}
