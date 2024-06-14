package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/attr_value"
)

type AttrValueDao interface {
	Dao[attr_value.AttrValue]
	FindAttrIDSByProductID(goodID uint64) (list []uint64, err error)
}

type attrValueDao struct {
	DaoIMPL[attr_value.AttrValue]
}

func NewAttrValueDao() AttrValueDao {
	return &attrValueDao{DaoIMPL[attr_value.AttrValue]{
		Db:    db.DB,
		Model: &attr_value.AttrValue{},
	}}
}

// FindAttrIDSByProductID 通过产品ID来查询所关联的标签ID
func (_this attrValueDao) FindAttrIDSByProductID(goodID uint64) (list []uint64, err error) {
	filter := map[string]any{
		"is_del":     "0",
		"product_id": goodID,
	}
	err = _this.Db.Model(_this.Model).Select("id").Where(filter).Find(&list).Error
	return
}
