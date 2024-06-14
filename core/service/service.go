package service

import (
	"log"
	baseDao "rentServer/core/dao"
	"strconv"
)

type BaseSearchParam struct {
	Page     string `json:"page" form:"page"`
	Size     string `json:"size" form:"size"`
	Order    string `json:"order" form:"order"`
	Keywords string `json:"keywords" `
	Where    any    `json:"where" `
}

type Service[T interface{}] interface {
	Save(data *T, id uint64) error
	List(param *BaseSearchParam) ([]*T, error)
	Delete(id []uint64) error
	FindByID(data *T, id any) error
	Update(data *T, id any) error
	FilterFindOne(data *T, filter map[string]any) error
	FilterList(param *BaseSearchParam, filter map[string]any) (list []*T, err error)
}

type ServiceIMPL[T interface{}] struct {
	BaseDao baseDao.Dao[T]
}

func (s ServiceIMPL[T]) Save(data *T, id uint64) error {
	log.Println(data, id, "asdasdadada===========")
	if id >= 1 {
		return s.BaseDao.Update(data, id)
	}
	return s.BaseDao.Insert(data)
}

func (s ServiceIMPL[T]) Delete(id []uint64) error {
	return s.BaseDao.Delete(id)
}

func (s ServiceIMPL[T]) List(param *BaseSearchParam) (list []*T, err error) {
	page, err := strconv.Atoi(param.Page)
	size, err := strconv.Atoi(param.Size)

	start := (page - 1) * size
	return s.BaseDao.ToList(start, size, param.Where, param.Order, param.Keywords)

}

func (s ServiceIMPL[T]) FindByID(data *T, id any) error {
	return s.BaseDao.FindByID(data, id)
}

func (s ServiceIMPL[T]) Update(data *T, id any) error {
	return s.BaseDao.Update(data, id)
}

func (s ServiceIMPL[T]) FilterFindOne(data *T, filter map[string]any) error {

	return s.BaseDao.FilterFindOne(data, filter)
}

func (s ServiceIMPL[T]) FilterList(param *BaseSearchParam, filter map[string]any) (list []*T, err error) {
	page, err := strconv.Atoi(param.Page)
	size, err := strconv.Atoi(param.Size)

	if err != nil {
		return nil, err
	}
	start := (page - 1) * size

	list, err = s.BaseDao.ToList(start, size, filter, "", "")
	if err != nil {
		return nil, err
	}
	return
}
