package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/customized_project"
)

var _ CustomProjDao = (*customProjDao)(nil)

type CustomProjDao interface {
	Dao[customized_project.CProject]
}

type customProjDao struct {
	DaoIMPL[customized_project.CProject]
}

func NewCustomProjDao() CustomProjDao {
	return &customProjDao{
		DaoIMPL: DaoIMPL[customized_project.CProject]{
			Db:    db.DB,
			Model: &customized_project.CProject{},
		},
	}
}
