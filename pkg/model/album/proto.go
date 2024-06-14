package album

import . "rentServer/core/model"

type AlbumFrom int

const (
	Customer = AlbumFrom(iota)
	ENV
	SERVICE
)

type Album struct {
	BaseModel
	StoreID string    `gorm:"column:store_id;type:varchar(255);comment:'店铺ID'" json:"store_id"`
	From    AlbumFrom `gorm:"column:from;comment:'相册来源 enum:1：客户上传 2：环境设施 3：服务特色'" json:"from"`
	Url     string    `gorm:"column:url;comment:'资源访问地址'" json:"url"`
}

func (_this Album) TableName() string {
	return "album"
}
