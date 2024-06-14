package good

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	. "rentServer/core/model"
	"rentServer/pkg/model/attr"
	"rentServer/pkg/model/attr_value"
	"rentServer/pkg/model/comment"
	"rentServer/pkg/model/user"
	"strconv"
)

type Product struct {
	BaseModel
	Cover       string                  `gorm:"column:cover; type:varchar(255)" json:"cover"`
	Title       string                  `gorm:"column:title;type:varchar(255)" json:"title"`
	Price       decimal.Decimal         `gorm:"column:price" json:"price"`
	Description string                  `gorm:"column:description" json:"description"`
	Banners     string                  `gorm:"column:banners;type:longtext" json:"-"`
	B           []string                `gorm:"-" json:"banners"`
	Stock       int                     `gorm:"column:stock" json:"stock"`
	Unit        string                  `gorm:"column:unit;type:varchar(255)" json:"unit"`
	Attrs       []*attr.Attr            `gorm:"-" json:"attrs"`
	Values      []*attr_value.AttrValue `gorm:"-" json:"values"`
	Comments    []*comment.Comment      `gorm:"-" json:"comments"`
	Detail      string                  `gorm:"column:detail;type:longtext" json:"detail"`
	CategoryID  string                  `gorm:"column:category_id;type:varchar(255)" json:"category_id"`
}

func (g Product) TableName() string {
	return "product"
}

type GoodComment struct {
	BaseModel

	Level   int    `gorm:"column:level" json:"level"`
	GoodID  string `gorm:"column:good_id;type:varchar(255)" json:"good_id"`
	UserID  string `gorm:"column:user_id;type:varchar(255)" json:"user_id"`
	Content string `gorm:"column:content" json:"content"`
}

func (_this GoodComment) TableName() string {
	return "good_comment"
}

func (_this Product) BeforeCreate(tx *gorm.DB) (err error) {

	banner, err := json.Marshal(&_this.B)
	if err != nil {
		return err
	}

	_this.Banners = string(banner)

	return nil
}

func (_this Product) AfterFind(tx *gorm.DB) (err error) {
	return json.Unmarshal([]byte(_this.Banners), &_this.B)
}

func (_this GoodComment) BuildDto(dto *GoodCommentDto) {

	dto.BaseModel = _this.BaseModel
	dto.Level = _this.Level
	dto.User = user.User{}
	dto.Good = Product{}
	dto.Content = _this.Content

}

func (_this GoodComment) BuildFromDto(dto *GoodCommentDto) GoodComment {
	return GoodComment{
		BaseModel: dto.BaseModel,
		Level:     dto.Level,
		GoodID:    strconv.FormatUint(dto.Good.ID, 10),
		UserID:    strconv.FormatUint(dto.User.ID, 10),
		Content:   dto.Content,
	}
}
