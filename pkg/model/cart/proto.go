package cart

import (
	"github.com/shopspring/decimal"
	. "rentServer/core/model"
)

type Cart struct {
	BaseModel
	Product string          `gorm:"column:product; type:longtext" json:"-"`
	P       CartProduct     `gorm:"-" json:"product"`
	Total   decimal.Decimal `gorm:"column:total;" json:"total"`
	// 用户收藏的 东西 是商品 good 还是 服务service
	Type string `gorm:"column:type;type:varchar(255)" json:"type"`
	//
	ProductID string `gorm:"column:product_id" json:"product_id"`
	UserID    string `gorm:"column:user_id;type:varchar(255)" json:"-"`
}

func (_this Cart) TableName() string {
	return "cart"
}
