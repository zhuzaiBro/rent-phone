package address

import . "rentServer/core/model"

// Address 用户收货地址
type Address struct {
	BaseModel
	City      string `gorm:"column:city;type:varchar(255)" json:"city"`            // 收货人所在市
	Detail    string `gorm:"column:detail;type:varchar(255)" json:"detail"`        // 收货人详细地址
	District  string `gorm:"column:district;type:varchar(255)" json:"district"`    // 收货人所在区
	Latitude  string `gorm:"column:latitude;type:varchar(1000)" json:"latitude"`   // 纬度
	Longitude string `gorm:"column:longitude;type:varchar(1000)" json:"longitude"` // 经度
	Phone     string `gorm:"column:phone;type:varchar(255)" json:"phone"`          // 收货人电话
	PostCode  string `gorm:"column:post_code" json:"post_code"`                    // 邮编
	Province  string `gorm:"column:province" json:"province"`                      // 收货人所在省
	RealName  string `gorm:"column:real_name" json:"real_name"`                    // 收货人姓名
	UserID    int64  `gorm:"column:user_id" json:"-"`                              // 用户id
}

func (_this Address) TableName() string {
	return "user_address"
}
