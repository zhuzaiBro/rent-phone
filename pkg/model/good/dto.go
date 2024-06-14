package good

import (
	. "rentServer/core/model"
	"rentServer/pkg/model/user"
)

type GoodCommentDto struct {
	BaseModel
	Level   int       `json:"level"`
	User    user.User `json:"user"`
	Good    Product   `json:"good"`
	Content string    `json:"content"`
}
