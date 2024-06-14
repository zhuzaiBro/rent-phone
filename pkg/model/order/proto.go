package order

import (
	"github.com/shopspring/decimal"
	. "rentServer/core/model"
	"rentServer/pkg/model/address"
	"rentServer/pkg/model/cart"
	"rentServer/pkg/model/coupon"
	"rentServer/pkg/model/store"
)

type OrderStatus int

const (
	// 待付款
	PENDING OrderStatus = iota
	// ING 进行中 活着是待收货
	ING
	// WAIT_COMMENT 待评价
	WAIT_COMMENT
	// 完成
	FINISHED
	// 取消
	CANCELED
)

type OrderPayType int

const (
	DATE OrderPayType = iota
	// WXPAY 微信支付
	WXPAY
	// BALANCE 余额
	BALANCE
)

type Order struct {
	// 创建时间和更新时间都在这里
	BaseModel
	// 订单状态
	Status OrderStatus `json:"status" gorm:"column:status"`

	// 关联的店员ID
	ClerkID string `json:"clerk_id" gorm:"column:clerk_id;type:varchar(255)"`

	// 订单类型 有服务 也有商品
	Type string `json:"type" gorm:"column:type;type:varchar(255)"`

	// 是否支付
	IsPayed bool `json:"is_payed" gorm:"column:is_payed"`
	// 是否评论
	IsCommented bool `json:"is_commented" gorm:"column:is_commented"`

	// 订单中的产品列表
	Products    string              `gorm:"column:products" json:"-"`
	ProductList []*cart.CartProduct `json:"product_list" gorm:"-"`

	// 订单总价
	Total decimal.Decimal `json:"total" gorm:"column:total;"`

	// 订单地址
	Address address.Address `json:"address" gorm:"-"  binding:"required"`
	Addr    string          `json:"-" gorm:"column:address"`

	// 流水号
	TradeNo string `json:"trade_no" gorm:"column:trade_no;type:varchar(255)"  `

	StoreID string      `json:"store_id" gorm:"column:store_id;type:varchar(255)"`
	Store   store.Store `json:"store" gorm:"-"`
	// 付款方式
	PayType OrderPayType `json:"pay_type" gorm:"column:pay_type"  `

	// 优惠券ID
	CouponRelID string `json:"coupon_rel_id" gorm:"coupon_rel_id"`

	// 优惠券
	Coupon coupon.Coupon `json:"coupon" gorm:"-"`

	UserID   string `json:"user_id" gorm:"column:user_id;type:varchar(255)" `
	UserName string `json:"user_name" gorm:"column:user_name;type:varchar(255);comment:'用户名称'"`

	// 买家备注
	Mark string `json:"mark" gorm:"column:mark; comment:'买家备注'"`

	StoreMark string `json:"store_mark" gorm:"column:store_mark; comment:'卖家备注'"`
}

func (_this Order) TableName() string {
	return "order"
}
