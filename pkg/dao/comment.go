package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/comment"
)

type CommentDao interface {
	Dao[comment.Comment]
}

type commentDao struct {
	DaoIMPL[comment.Comment]
}

func NewCommentDao() CommentDao {
	return &commentDao{
		DaoIMPL: DaoIMPL[comment.Comment]{
			Db:    db.DB,
			Model: &comment.Comment{},
		},
	}
}
