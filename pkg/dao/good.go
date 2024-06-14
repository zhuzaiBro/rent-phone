package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/good"
)

type GoodDao interface {
	Dao[good.Product]
}

type goodDao struct {
	DaoIMPL[good.Product]
}

func NewGoodDao() GoodDao {
	return &goodDao{DaoIMPL[good.Product]{
		Db:    db.DB,
		Model: &good.Product{},
	}}
}
