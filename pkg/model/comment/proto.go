package comment

import (
	. "rentServer/core/model"
	"rentServer/pkg/model/cart"
	"rentServer/pkg/model/user"
)

type Score int

const (
	BAD Score = iota
	TWO
	THREE
	FOUR
	GREAT
)

type Comment struct {
	BaseModel
	StoreID      string           `gorm:"column:store_id;type:varchar(255)" json:"store_id"`
	OrderID      string           `gorm:"column:order_id;type:varchar(255)" json:"order_id"`
	ProductID    string           `gorm:"column:product_id;type:varchar(255)" json:"product_id"`
	ClerkID      string           `gorm:"column:clerk_id;type:varchar(255);comment:'店员ID'" json:"clerk_id"`
	Quality      Score            `gorm:"column:quality" json:"quality"`
	Service      Score            `gorm:"column:service" json:"service"`
	PicList      string           `gorm:"column:pic_list" json:"-"`
	Content      string           `gorm:"column:content" json:"content"`
	Hide         bool             `gorm:"column:hide" json:"hide"`
	P            []string         `gorm:"-" json:"pic_list"`
	ReplyContent string           `gorm:"column:reply_content" json:"reply_content"`
	User         user.User        `gorm:"-" json:"user"`
	Product      cart.CartProduct `gorm:"-" json:"product"`
	UserID       string           `gorm:"column:user_id" json:"user_id"`
}

func (_this Comment) TableName() string {
	return "user_comments"
}
