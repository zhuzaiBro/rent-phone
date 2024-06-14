package dao

import (
	"errors"
	"gorm.io/gorm"
	"log"
)

type Dao[T interface{}] interface {
	Insert(data *T) error
	FindByID(data *T, id interface{}) error
	Delete(id interface{}) error
	Update(data *T, id interface{}) error
	ToList(start int, size int, where any, order string, keyword string) (list []*T, err error)
	Exit(id any) bool
	FilterFindOne(data *T, filter map[string]any) error
	UpdateFromMap(fields map[string]any, where map[string]any) error
}

type DaoIMPL[T interface{}] struct {
	Db    *gorm.DB
	Model *T
}

func (d DaoIMPL[T]) Insert(data *T) error {
	return d.Db.Model(data).Create(data).Error
}

func (d DaoIMPL[T]) FindByID(data *T, id interface{}) error {
	return d.Db.Model(data).Where("is_del", 0).Where("id", id).First(data).Error
}

func (d DaoIMPL[T]) Delete(id interface{}) error {
	return d.Db.Model(d.Model).Where("id", id).Update("is_del", 1).Error
}

func (d DaoIMPL[T]) Update(data *T, id interface{}) error {
	return d.Db.Model(d.Model).Where("id", id).Updates(data).Error
}

func (d DaoIMPL[T]) ToList(start int, size int, where any, order string, keyword string) (list []*T, err error) {
	if len(order) == 0 {
		order = "id desc"
	}
	if start <= 0 {
		start = 0
	}
	if size < 0 {
		size = 10
	}

	log.Println(start, size, where, order, keyword)

	err = d.Db.Model(d.Model).Where("is_del", 0).Where(where).Where(keyword).Order(order).Limit(size).Offset(start).Find(&list).Error
	return list, err
}

func (d DaoIMPL[T]) Exit(order any) bool {
	e := d.Db.Model(d.Model).Where(order).First(d.Model).Error
	if errors.Is(e, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (d DaoIMPL[T]) FilterFindOne(data *T, filter map[string]any) error {
	filter["is_del"] = "0"

	return d.Db.Model(data).Where(filter).First(data).Error
}

func (d DaoIMPL[T]) UpdateFromMap(fields map[string]any, where map[string]any) error {

	where["is_del"] = "0"
	return d.Db.Model(d.Model).Where(where).Updates(fields).Error
}
