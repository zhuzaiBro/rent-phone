package customized_project

import . "rentServer/core/model"

type CProject struct {
	BaseModel

	IsMain bool `json:"is_main" gorm:"column:is_main;"`

	Name     string `gorm:"column:name;type:varchar(200)" json:"name"`
	Desc     string `gorm:"column:desc; type:varchar(255)" json:"desc"`
	Strength int    `gorm:"column:strength" json:"strength"`
	MinTime  int    `gorm:"column:min_time;comment:'最短服务时间30min/60min不等'" json:"min_time"`
	Info     string `gorm:"column:info;type:varchar(1000);" json:"info"`
}

func (_this CProject) TableName() string {
	return "custom_proj"
}
