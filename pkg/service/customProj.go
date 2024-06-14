package service

import (
	. "rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/model/customized_project"
)

type CustomProjService interface {
	Service[customized_project.CProject]
}

type customProjService struct {
	ServiceIMPL[customized_project.CProject]
	BaseDao dao.CustomProjDao
}

func NewCustomProjService() CustomProjService {
	customProjDao := dao.NewCustomProjDao()
	return &customProjService{
		ServiceIMPL: ServiceIMPL[customized_project.CProject]{
			BaseDao: customProjDao,
		},
		BaseDao: customProjDao,
	}
}
