package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/attr"
)

var _ AttrDao = (*attrDao)(nil)

type AttrDao interface {
	Dao[attr.Attr]
	FindAttrIDSByProductID(goodID uint64) (list []uint64, err error)
}

type attrDao struct {
	DaoIMPL[attr.Attr]
}

func NewAttrDao() AttrDao {
	return &attrDao{
		DaoIMPL[attr.Attr]{
			Db:    db.DB,
			Model: &attr.Attr{},
		},
	}
}

func (_this attrDao) FindAttrIDSByProductID(goodID uint64) (list []uint64, err error) {
	filter := map[string]any{
		"is_del":     "0",
		"product_id": goodID,
	}
	err = _this.Db.Model(_this.Model).Select("id").Where(filter).Find(&list).Error
	return
}
