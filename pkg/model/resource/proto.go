package resource

import . "rentServer/core/model"

type Resource struct {
	BaseModel
	Name    string `gorm:"column:name;type:varchar(255);comment:'文件名称'" json:"name"`
	Url     string `gorm:"column:url;type:varchar(500);comment:'链接'" json:"url"`
	GroupID string `gorm:"column:group_id;type:varchar(255);comment:'文件分组ID'" json:"group_id"`
	Type    string `gorm:"column:type;type:varchar(255);comment:'文件类型'" json:"type"`
}

func (_this Resource) TableName() string {
	return "resource"
}
