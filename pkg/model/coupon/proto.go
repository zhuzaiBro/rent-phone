package coupon

import (
	"github.com/shopspring/decimal"
	. "rentServer/core/model"
)

type CouponStatus int

const (
	UNUSE CouponStatus = iota
	USED
)

const (
	PERCENTAGE = "percentage"
	FIXED      = "fixed"
)

type Coupon struct {
	BaseModel
	Name            string          `gorm:"column:name;type:varchar(255)" json:"name"`
	Code            string          `gorm:"column:code;type:VARCHAR(255);comment:'优惠券代码'" json:"code"`
	DiscountType    string          `gorm:"column:discount_type,type:ENUM('percentage', 'fixed');comment:'折扣类型：百分比或固定金额'" json:"discount_type"`
	DiscountValue   decimal.Decimal `gorm:"column:discount_value;comment:'折扣值：表示折扣的金额或百分比'" json:"discount_value"`
	StartDate       int64           `gorm:"column:start_date;comment:'优惠券生效开始日期'" json:"start_date"`
	EndDate         int64           `gorm:"column:end_date;comment:'优惠券失效结束日期'" json:"end_date"`
	UsageLimit      int             `gorm:"column:usage_limit;comment:'优惠券可使用次数限制'" json:"usage_limit"`
	UsedTimes       int             `gorm:"column:used_times;comment:'已使用次数'" json:"used_times"`
	Status          CouponStatus    `gorm:"column:status;comment:'ENUM('active', 'expired', 'disabled')	优惠券状态：激活、过期、禁用'" json:"status"`
	ProductCategory string          `gorm:"column:product_category;type:varchar(255);comment:'适用的商品类目'" json:"product_category"`
	CategoryID      string          `gorm:"column:category_id" json:"-"`
	ProductID       string          `gorm:"column:product_id;comment:'适用的商品唯一标识符（可选）'" json:"product_id"`
	MinPurchase     decimal.Decimal `gorm:"column:min_purchase;comment:'最低消费金额限制（可选）'" json:"min_purchase"`
}

func (_this Coupon) TableName() string {
	return "coupon"
}
