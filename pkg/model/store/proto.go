package store

import (
	. "rentServer/core/model"
	"rentServer/pkg/model/customized_project"
)

type Store struct {
	BaseModel
	CityName          string                         `gorm:"column:city_name" json:"city_name,omitempty"`       // 商户城市名称
	ContactName       string                         `gorm:"column:contact_name" json:"contact_name,omitempty"` // 联系人姓名
	ContactPhone      string                         `gorm:"contact_phone" json:"contact_phone,omitempty"`      // 联系人电话
	Name              string                         `gorm:"column:name;type:varchar(255)" json:"name"`
	Info              string                         `gorm:"column:info" json:"info"`
	Lat               float32                        `gorm:"column:lat" json:"lat"`
	Lng               float32                        `gorm:"column:lng" json:"lng"`
	Email             string                         `gorm:"column:email" json:"email,omitempty"`                           // 邮箱地址
	EnterpriseAddress string                         `gorm:"column:enterprise_address" json:"enterprise_address,omitempty"` // 企业地址
	EnterpriseName    string                         `gorm:"column:enterprise_name" json:"enterprise_name,omitempty"`       // 企业全称
	IsUu              bool                           `gorm:"column:is_uu" json:"is_uu"`                                     // 0、否1、是
	MchID             uint64                         `gorm:"column:mch_id" json:"mch_id"`                                   // 商户id
	Mobile            string                         `gorm:"column:mobile" json:"mobile,omitempty"`                         // 手机号
	SourceID          string                         `gorm:"column:source_id" json:"source_id,omitempty"`                   // 商户编号
	StoreID           int64                          `gorm:"column:store_id" json:"store_id"`                               // 商城id
	WxAppID           string                         `gorm:"column:wx_app_id" json:"wx_app_id,omitempty"`                   // uuappid
	WxAppSecret       string                         `gorm:"column:wx_app_secret" json:"wx_app_secret,omitempty"`           // uuSecret Key
	WxPayPrivateKey   string                         `gorm:"column:wx_pay_private_key" json:"wx_pay_private_key,omitempty"` //
	Logo              string                         `gorm:"column:logo" json:"logo"`
	Bg                string                         `gorm:"column:bg" json:"bg"`
	Video             string                         `gorm:"column:video" json:"video"`
	ServiceTime       any                            `gorm:"-" json:"service_time"`
	Projects          []*customized_project.CProject `gorm:"-" json:"projects"`
}

func (_this Store) TableName() string {
	return "store"
}
