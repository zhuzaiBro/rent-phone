package wechat_template

import . "rentServer/core/model"

const (
	TEMPLATE  = "template"
	SUBSCRIBE = "订阅消息"
)

// WechatTemplate ，微信模板
type WechatTemplate struct {
	BaseModel
	Content string `json:"content" gorm:"column:content;type:varchar(1000);comment:'回复内容'"`
	Name    string `json:"name" gorm:"column:name;type:varchar(255);comment:'模板名'"`
	Status  bool   `json:"status" gorm:"column:status;comment:'状态'"`
	TempId  string `json:"temp_id,omitempty" gorm:"column:temp_id;type:varchar(255);comment:'模板ID'"`
	TempKey string `json:"temp_key" gorm:"column:temp_key;type:varchar(255);comment:'模板编号'"`
	Type    string `json:"type,omitempty" gorm:"column:type;type:varchar(255);comment:'类型：template:模板消息 subscribe:订阅消息'"`
}

func (_this WechatTemplate) TableName() string {
	return "wechat_template"
}
