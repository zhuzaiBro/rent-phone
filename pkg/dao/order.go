package dao

import (
	"gorm.io/gorm"
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/order"
)

type OrderDao interface {
	Dao[order.Order]
	RedisCache[order.Order]
	GetOrderCount(filter map[string]any) (count int64, err error)
	TimeConditionList(tx *gorm.DB, filter map[string]any) (list []*order.Order, err error)
}

type orderDao struct {
	DaoIMPL[order.Order]
	RedisCacheIMPL[order.Order]
}

func NewOrderDao() OrderDao {
	return &orderDao{
		DaoIMPL: DaoIMPL[order.Order]{
			Db:    db.DB,
			Model: &order.Order{},
		},
		RedisCacheIMPL: RedisCacheIMPL[order.Order]{
			db.GetRedis(),
		},
	}
}

func (_this orderDao) GetOrderCount(filter map[string]any) (count int64, err error) {
	err = _this.Db.Model(_this.Model).Where(filter).Count(&count).Error
	return
}

func (_this orderDao) TimeConditionList(tx *gorm.DB, filter map[string]any) (list []*order.Order, err error) {

	err = tx.Where(filter).Find(&list).Error

	return
}
