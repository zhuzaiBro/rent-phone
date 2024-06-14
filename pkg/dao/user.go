package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/user"
)

type UserDao interface {
	Dao[user.User]
}

type userDao struct {
	DaoIMPL[user.User]
}

func NewUserDao() UserDao {
	return &userDao{
		DaoIMPL: DaoIMPL[user.User]{
			Db:    db.DB,
			Model: &user.User{},
		},
	}
}
