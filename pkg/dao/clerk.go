package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/clerk"
)

var _ ClerkDao = (*clerkDao)(nil)

type ClerkDao interface {
	Dao[clerk.Clerk]
}

type clerkDao struct {
	DaoIMPL[clerk.Clerk]
}

func NewClerkDao() ClerkDao {
	return &clerkDao{
		DaoIMPL: DaoIMPL[clerk.Clerk]{
			Db:    db.DB,
			Model: &clerk.Clerk{},
		},
	}
}
