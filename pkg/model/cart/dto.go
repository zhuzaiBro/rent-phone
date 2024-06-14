package cart

import "github.com/shopspring/decimal"

type AddRequest []*AddItem

type AddItem struct {
	UniqueID string `json:"unique_id"`
	Num      int    `json:"num"`
}

type CartProduct struct {
	Cover string          `json:"cover"`
	Title string          `json:"title"`
	Price decimal.Decimal `json:"price"`
	// 选购产品数量
	Num int `json:"num"`
	// 预留字段，多商户模式开启
	MerID string `json:"mer_id"`

	// 前端通过这两个字段跳转 包括购买，后端通过这两个查询处理sku
	ProductType string `json:"product_type"`
	ProductID   string `json:"product_id"`

	Sku string `json:"sku"`

	ArriveTime string `json:"arrive_time"`
}
