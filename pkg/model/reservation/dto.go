package reservation

import (
	"rentServer/pkg/model/clerk"
)

type Reservation struct {
	People int `gorm:"column:people;comment:'预约人数'" json:"people"`

	// 以小时为单位
	TotalTime float64 `gorm:"column:total_time;comment:'项目服务时间'" json:"total_time"`

	NeedShower bool `gorm:"column:need_shower;comment:'是否淋浴'" json:"need_shower"`

	Clerk clerk.Clerk `gorm:"-" json:"clerk"`

	ArriveTime string `gorm:"-" json:"arrive_time"`
}
