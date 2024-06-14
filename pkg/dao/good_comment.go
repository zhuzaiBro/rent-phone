package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/good"
)

type GoodCommentDao interface {
	Dao[good.GoodComment]
}

type goodCommentDao struct {
	DaoIMPL[good.GoodComment]
}

func NewGoodCommentDao() GoodCommentDao {
	return &goodCommentDao{DaoIMPL[good.GoodComment]{
		Db:    db.DB,
		Model: &good.GoodComment{},
	}}
}
