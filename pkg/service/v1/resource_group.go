package v1

import (
	. "rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/model/resource_group"
)

var _ ResourceGroupService = (*resourceGroupServiceIMPL)(nil)

type ResourceGroupService interface {
	Service[resource_group.ResourceGroup]
}

type resourceGroupServiceIMPL struct {
	ServiceIMPL[resource_group.ResourceGroup]
	ResourceGroupDao dao.ResourceGroupDao
}

func NewResourceGroupService() ResourceGroupService {
	resourceGroupDao := dao.NewResourceGroupDao()

	return &resourceGroupServiceIMPL{
		ServiceIMPL: ServiceIMPL[resource_group.ResourceGroup]{
			BaseDao: resourceGroupDao,
		},
		ResourceGroupDao: resourceGroupDao,
	}
}
