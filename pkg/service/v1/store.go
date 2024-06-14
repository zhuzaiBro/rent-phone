package v1

import (
	. "rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/model/store"
)

var _ StoreService = (*storeServiceIMPL)(nil)

type StoreService interface {
	Service[store.Store]
}

type storeServiceIMPL struct {
	ServiceIMPL[store.Store]
}

func NewStoreService() StoreService {
	storeDao := dao.NewStoreDao()
	return &storeServiceIMPL{
		ServiceIMPL: ServiceIMPL[store.Store]{
			BaseDao: storeDao,
		},
	}
}
