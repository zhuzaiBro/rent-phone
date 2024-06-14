package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/store"
)

var _ StoreDao = (*storeDao)(nil)

type StoreDao interface {
	Dao[store.Store]
}

type storeDao struct {
	DaoIMPL[store.Store]
}

func NewStoreDao() StoreDao {
	return &storeDao{
		DaoIMPL: DaoIMPL[store.Store]{
			Db:    db.DB,
			Model: &store.Store{},
		},
	}
}
