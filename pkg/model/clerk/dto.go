package clerk

import (
	. "rentServer/core/model"
	"rentServer/pkg/model/user"
)

type ClerkDto struct {
	BaseModel
}

type ClerkCommentDto struct {
	BaseModel
	User    user.User `json:"user_info"`
	Content string    `json:"content"`
	Level   int       `json:"level"`
	// 评论的对象类型， 有服务员 也有产品
	Type    string `json:"type"`
	ClerkID string ` json:"object_id"`
}
