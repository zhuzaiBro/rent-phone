package order

type OrderParam struct {
	Status      OrderStatus `json:"status" form:"status"`
	IsPaid      bool        `json:"is_payed" form:"is_payed"`
	IsCommented bool        `json:"is_commented" form:"is_commented"`
	Type        string      `json:"type" form:"type"`
	UserID      uint64      `json:"user_id"`
}

type OrderCount struct {
	// 所有订单
	Total int64 `json:"total"`
	// 未支付
	UnPayed int64 `json:"un_payed"`
	// 正在进行
	ING int64 `json:"ing"`
	// 待评价
	UnComment int64 `json:"un_comment"`
	// 完成
	Finish int64 `json:"finish"`
}
