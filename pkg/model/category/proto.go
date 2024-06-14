package category

import . "rentServer/core/model"

type Category struct {
	BaseModel
	Name     string      `gorm:"column:name;type:varchar(255);" json:"name"`
	Sort     int         `gorm:"column:sort" json:"sort"`
	StoreID  string      `gorm:"column:store_id" json:"store_id"`
	Cover    string      `gorm:"column:cover" json:"cover"`
	Children []*Category `gorm:"-" json:"children"`
	ParentID string      `gorm:"column:parent_id" json:"parent_id"`
}

func (_this Category) TableName() string {
	return "category"
}
