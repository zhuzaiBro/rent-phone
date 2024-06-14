package admin

import . "rentServer/core/model"

type Admin struct {
	BaseModel

	Nickname string `gorm:"nickname" json:"nickname"`

	UserName string `gorm:"column:username;type:varchar(255)" json:"username"`
	Password string `gorm:"column:password" json:"password"`

	Mobile string `gorm:"column:mobile;type:varchar(255)" json:"mobile"`
	Mail   string `gorm:"column:mail;type:varchar(255)" json:"mail"`
}

func (a Admin) TableName() string {
	return "admin"
}
