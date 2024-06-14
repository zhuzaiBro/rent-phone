package resource_group

import . "rentServer/core/model"

type ResourceGroup struct {
	BaseModel
	Name     string `gorm:"column:name;type:varchar(255)" json:"name"`
	CreateID string `gorm:"column:creator_id;type:varchar(255)" json:"create_id"`
	Sort     int    `gorm:"column:sort" json:"sort"`
}

func (_this ResourceGroup) TableName() string {
	return "resource_group"
}
