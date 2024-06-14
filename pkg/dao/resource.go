package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/resource"
)

var _ ResourceDao = (*resourceDao)(nil)

type ResourceDao interface {
	Dao[resource.Resource]
}

type resourceDao struct {
	DaoIMPL[resource.Resource]
}

func NewResourceDao() ResourceDao {
	return &resourceDao{
		DaoIMPL[resource.Resource]{
			Db:    db.DB,
			Model: &resource.Resource{},
		},
	}
}
