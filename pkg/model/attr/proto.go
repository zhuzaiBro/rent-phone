package attr

import (
	"encoding/json"
	. "rentServer/core/model"
)

type Attr struct {
	BaseModel
	D         string   `gorm:"column:detail;default:'[]';type:varchar(1000)" json:"-"`
	Detail    []string `json:"detail" gorm:"-"`
	Value     string   `gorm:"column:value" json:"value"`
	ProductID string   `gorm:"column:product_id" json:"-"`
	Index     int      `gorm:"-" json:"index"`
}

func (_this Attr) TableName() string {
	return "product_attr"
}

func (_this Attr) Detail2D() (str string, err error) {
	d, err := json.Marshal(&_this.Detail)
	if err != nil {
		return "", err
	}
	str = string(d)
	return
}

func (_this Attr) D2Detail() (list []string, err error) {
	err = json.Unmarshal([]byte(_this.D), &list)
	return
}

//
//func (_this Attr) BeforeCreate(tx *gorm.DB) (err error) {
//
//	_this.D, err = _this.Detail2D()
//	if err != nil {
//		if _this.D == "" {
//			_this.D = "[]"
//		} else {
//			return err
//		}
//
//	}
//	logger.Info(_this.D, _this.Detail)
//	return nil
//}
//
//func (_this Attr) BeforeUpdate(tx *gorm.DB) (err error) {
//	_this.D, err = _this.Detail2D()
//	if err != nil {
//		if _this.D == "" {
//			_this.D = "[]"
//		} else {
//			return err
//		}
//
//	}
//	return nil
//}
//
//func (_this Attr) AfterFind(tx *gorm.DB) (err error) {
//	_this.Detail, err = _this.D2Detail()
//	return
//}
