package attr_value

import . "rentServer/core/model"

type AttrValue struct {
	BaseModel
	// 业务代码
	BarCode string `gorm:"column:bar_code;type:varchar(255)" json:"bar_code"`
	// 成本
	Cost string `gorm:"column:cost;" json:"cost"`
	// 图像
	Image string `gorm:"column:image" json:"image"`
	// 积分
	Integral int64 `gorm:"column:integral" json:"integral"`
	// 原价
	OtPrice int64 `gorm:"column:ot_price" json:"ot_price"`
	// 团购价
	PinkPrice string `gorm:"column:pink_price" json:"pink_price"`
	// 团购库存
	PinkStock int64  `gorm:"column:pink_stock" json:"pink_stock"`
	Price     string `gorm:"column:price" json:"price"`
	ProductID string `gorm:"column:product_id;type:varchar(255)" json:"product_id"`
	// 销量
	Sales int64 `gorm:"column:sales;comment:'销量'" json:"sales"`
	// 优惠价格
	SeckillPrice string `gorm:"column:seckill_price;comment:'优惠价格'" json:"seckill_price"`
	// 优惠库存
	SeckillStock int64  `gorm:"column:seckill_stock;comment:'优惠价格的库存';" json:"seckill_stock"`
	Sku          string `gorm:"column:sku" json:"sku"`
	Stock        int64  `gorm:"column:stock" json:"stock"`
	Unique       string `gorm:"column:unique" json:"unique"`
	Volume       int64  `gorm:"column:volume" json:"volume"`
	Weight       int64  `gorm:"column:weight" json:"weight"`
}

func (_this AttrValue) TableName() string {
	return "attr_value"
}
