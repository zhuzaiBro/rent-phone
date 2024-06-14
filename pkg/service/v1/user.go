package v1

import (
	. "rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/model/user"
	"strconv"
)

type UserService interface {
	List(param *BaseSearchParam, filter map[string]any) (list []*user.User, err error)
	Update(userModel *user.User) error
}

type userService struct {
	UserDao dao.UserDao
}

func NewUserService() UserService {
	return &userService{
		UserDao: dao.NewUserDao(),
	}
}

func (_this userService) List(param *BaseSearchParam, filter map[string]any) (list []*user.User, err error) {

	page, err := strconv.Atoi(param.Page)
	size, err := strconv.Atoi(param.Page)
	if err != nil {
		return
	}

	start := (page - 1) * size
	list, err = _this.UserDao.ToList(start, size, filter, "", "")
	if err != nil {
		return nil, err
	}

	return
}

func (_this userService) Update(userModel *user.User) error {

	return _this.UserDao.Update(userModel, userModel.ID)

}
