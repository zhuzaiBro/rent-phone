package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/clerk"
)

type ClerkCommentDao interface {
	Dao[clerk.Comment]
}

type clerkCommentDao struct {
	DaoIMPL[clerk.Comment]
}

func NewClerkCommentDao() ClerkCommentDao {
	return &clerkCommentDao{DaoIMPL[clerk.Comment]{
		Db:    db.DB,
		Model: &clerk.Comment{},
	}}
}
