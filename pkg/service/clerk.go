package service

import (
	. "rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/model/clerk"
)

var _ ClerkService = (*clerkService)(nil)

type ClerkService interface {
	Service[clerk.Clerk]
}

type clerkService struct {
	ServiceIMPL[clerk.Clerk]
	ClerkDao dao.ClerkDao
}

func NewClerkService() ClerkService {
	clerkDao := dao.NewClerkDao()
	return &clerkService{
		ServiceIMPL: ServiceIMPL[clerk.Clerk]{
			BaseDao: clerkDao,
		},
		ClerkDao: clerkDao,
	}
}
