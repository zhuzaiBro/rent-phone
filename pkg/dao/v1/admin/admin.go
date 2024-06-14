package admin

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/admin"
)

type AdminDao interface {
	Dao[admin.Admin]
}

type adminDao struct {
	DaoIMPL[admin.Admin]
}

func NewAdminDao() AdminDao {
	return adminDao{DaoIMPL[admin.Admin]{
		Db:    db.DB,
		Model: &admin.Admin{},
	}}
}
