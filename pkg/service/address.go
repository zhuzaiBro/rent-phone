package service

import (
	. "rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/model/address"
)

type AddressService interface {
	Service[address.Address]
}

type addressServiceIMPL struct {
	ServiceIMPL[address.Address]
	AddressDao dao.AddressDao
}

func NewAddressService() AddressService {
	addressDao := dao.NewAddressDao()
	return &addressServiceIMPL{
		ServiceIMPL: ServiceIMPL[address.Address]{
			BaseDao: addressDao,
		},
		AddressDao: addressDao,
	}
}
