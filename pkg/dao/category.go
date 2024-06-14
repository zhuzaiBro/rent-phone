package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/category"
)

var _ CategoryDao = (*categoryDao)(nil)

type CategoryDao interface {
	Dao[category.Category]
}

type categoryDao struct {
	DaoIMPL[category.Category]
}

func NewCategoryDao() CategoryDao {
	return &categoryDao{
		DaoIMPL: DaoIMPL[category.Category]{
			Db:    db.DB,
			Model: &category.Category{},
		},
	}
}
