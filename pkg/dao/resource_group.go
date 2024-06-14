package dao

import (
	"rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/resource_group"
)

var _ ResourceGroupDao = (*resourceGroupDao)(nil)

type ResourceGroupDao interface {
	dao.Dao[resource_group.ResourceGroup]
}

type resourceGroupDao struct {
	dao.DaoIMPL[resource_group.ResourceGroup]
}

func NewResourceGroupDao() ResourceGroupDao {
	return &resourceGroupDao{
		DaoIMPL: dao.DaoIMPL[resource_group.ResourceGroup]{
			Db:    db.DB,
			Model: &resource_group.ResourceGroup{},
		},
	}
}
