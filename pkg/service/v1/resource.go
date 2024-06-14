package v1

import (
	. "rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/model/resource"
)

var _ ResourceService = (*resourceService)(nil)

type ResourceService interface {
	Service[resource.Resource]
}

type resourceService struct {
	ServiceIMPL[resource.Resource]
	ResourceDao dao.ResourceDao
}

func NewResourceService() ResourceService {
	resourceDao := dao.NewResourceDao()
	return resourceService{
		ServiceIMPL[resource.Resource]{
			BaseDao: resourceDao,
		},
		resourceDao,
	}
}
