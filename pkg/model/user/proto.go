package user

import (
	"github.com/shopspring/decimal"
	coreModel "rentServer/core/model"
	"rentServer/pkg/model/order"
)

type User struct {
	coreModel.BaseModel

	Nickname string `gorm:"column:nickname; type:varchar(255);" json:"nickname"`

	WxMpData string `gorm:"column:wx_mp_data; type:longtext;" json:"-"`

	WxOpenID string `gorm:"column:wx_open_id;type:varchar(255)" json:"-"`

	VipNum string `gorm:"column:vip_num; type: varchar(100);" json:"vip_num"`

	Avatar string `gorm:"column:avatar; type: varchar(255)" json:"avatar"`

	Mobile string `gorm:"column:mobile; type: varchar(255)" json:"mobile"`

	Account decimal.Decimal `gorm:"column:account; " json:"account"`

	ServiceOrderCount int64 `gorm:"-" json:"service_order_count"`

	ProductOrderCount int64 `gorm:"-" json:"product_order_count"`

	CurrentOrder order.Order `gorm:"-" json:"current_order"`
	// 用户默认地址
	AddressID string `gorm:"address_id" json:"address_id"`
}

func (u User) TableName() string {
	return "user"
}
