package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/address"
)

var _ AddressDao = (*addressDao)(nil)

type AddressDao interface {
	Dao[address.Address]
}

type addressDao struct {
	DaoIMPL[address.Address]
}

func NewAddressDao() AddressDao {
	return &addressDao{
		DaoIMPL[address.Address]{
			Db:    db.DB,
			Model: &address.Address{},
		},
	}
}
