package clerk

import (
	. "rentServer/core/model"
	"rentServer/pkg/model/user"
)

type Clerk struct {
	BaseModel
	Name       string  `json:"name" gorm:"column:name;type:varchar(255)"`
	EnName     string  `json:"en_name" gorm:"column:en_name;type:varchar(255)"`
	Commission float64 `json:"commission" gorm:"column:commission"`
	Avatar     string  `json:"avatar" json:"avatar"`
	Level      int     `json:"level" gorm:"column:level"`

	ServiceTime any `json:"service_time" gorm:"-"`
}

func (_this Clerk) TableName() string {
	return "clerk"
}

type Comment struct {
	BaseModel
	UserID  string `gorm:"column:user_id;type:varchar(255)" json:"user_id"`
	Content string `gorm:"column:content; type:varchar(1000)" json:"content"`
	Level   int    `gorm:"column:level" json:"level"`
	// 评论的对象ID
	Type    string `gorm:"column:type;type:varchar(255)" json:"type"`
	ClerkID string `gorm:"column:clerk_id;type:varchar(255)" json:"clerk_id"`
}

func (_this Comment) TableName() string {
	return "clerk_comment"
}
func (_this Comment) BuildDto(dto *ClerkCommentDto) {
	dto.BaseModel = _this.BaseModel
	dto.User = user.User{}
	dto.Content = _this.Content
	dto.Level = _this.Level
	dto.Type = _this.Type
	dto.ClerkID = _this.ClerkID

}
