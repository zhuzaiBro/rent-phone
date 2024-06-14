package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/cart"
)

type CartDao interface {
	Dao[cart.Cart]
	GetListByFilter(filter map[string]any) (list []*cart.Cart, err error)
}

type cartDao struct {
	DaoIMPL[cart.Cart]
}

func NewCartDao() CartDao {
	return &cartDao{DaoIMPL: DaoIMPL[cart.Cart]{
		Db:    db.DB,
		Model: &cart.Cart{},
	}}
}

func (_this cartDao) GetListByFilter(filter map[string]any) (list []*cart.Cart, err error) {
	filter["is_del"] = 0
	err = _this.Db.Model(_this.Model).Where(filter).Find(&list).Error
	return
}
