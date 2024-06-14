package v2

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	. "rentServer/core/controller"
	cosClient "rentServer/initilization/resource"
	"rentServer/pkg/config"
	"rentServer/pkg/logger"
	"rentServer/pkg/model/user"
	"rentServer/pkg/model/wx_api"
	"rentServer/pkg/oauth"
	"rentServer/pkg/service"
	"strconv"
	"strings"
	"time"
)

type UserController interface {
	GetUserInfo(c *Context)
	MpLogin(c *Context)
	UpdateUserInfo(c *Context)
	BindUserWxPhone(c *Context)
	UploadAvatar(c *Context)
}

type userController struct {
	UserSvc  service.UserService
	OrderSvc service.OrderService
}

func NewUserController() UserController {
	return &userController{
		UserSvc:  service.NewUserService(),
		OrderSvc: service.NewOrderService(),
	}
}

func (_this userController) MpLogin(c *Context) {

	var dto wx_api.RequestDto

	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSONOK(ErrBadRequest)
		return
	}

	token, err := _this.UserSvc.WxMpLogin(dto)
	if err != nil {
		c.JSONOK(ErrInternalServer)
		return
	}

	c.JSONOK(token)
}

func (_this userController) GetUserInfo(c *Context) {

	filter := map[string]any{
		"id": c.GetUint64(Uid),
	}
	var userModel user.User
	err := _this.UserSvc.FilterFindOne(&userModel, filter)
	if err != nil {

		c.JSONOK(err.Error())
		return
	}
	filter1 := map[string]any{
		"user_id": userModel.ID,
		"type":    service.SERVICE,
		"is_del":  false,
	}
	userModel.ServiceOrderCount, err = _this.OrderSvc.GetMyOrderCount(filter1)

	filter2 := map[string]any{
		"user_id": userModel.ID,
		"type":    service.GOOD,
		"is_del":  false,
	}
	userModel.ProductOrderCount, err = _this.OrderSvc.GetMyOrderCount(filter2)

	if err != nil {

		c.JSONOK(err.Error())
		return
	}

	today := time.Now().Format("2006-01-02")
	conditions := map[string]any{
		"create_at": today,
		"user_id":   userModel.ID,
		"is_del":    false,
		// 其他条件...
	}
	err = _this.OrderSvc.FilterFindOne(&userModel.CurrentOrder, conditions)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSONOK(err.Error())
		return
	}

	c.JSONOK(userModel)
}

func (_this userController) UpdateUserInfo(c *Context) {
	var userModel user.User

	err := c.ShouldBindJSON(&userModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = _this.UserSvc.Update(&userModel, c.GetUint64(Uid))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK(userModel)
}

func (_this userController) BindUserWxPhone(c *Context) {
	var (
		userModel user.User
		req       wx_api.RequestDto
		phoneDto  wx_api.WxPhoneDto
	)

	err := c.ShouldBindJSON(&req)

	phoneStr, err := oauth.GetPhone(req)
	if err != nil {
		c.JSONOK(err.Error())
		return
	}

	err = json.Unmarshal(phoneStr, &phoneDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = _this.UserSvc.FindByID(&userModel, c.GetUint64(Uid))
	if err != nil {
		c.JSONOK(err.Error())
		return
	}

	userModel.Mobile = phoneDto.PurePhoneNumber
	err = _this.UserSvc.Update(&userModel, c.GetUint64(Uid))
	if err != nil {
		c.JSONOK(err.Error())
		return
	}

	c.JSONOK(userModel)
}

func (_this userController) UploadAvatar(c *Context) {

	_file, err := c.FormFile(Upload)

	if err != nil {
		logger.Error(err.Error())
		c.JSONOK(ErrBadRequest)
		return
	}
	n := time.Now()

	ns := strconv.Itoa(int(n.UnixMicro()))

	r := rand.Int()
	rs := strconv.Itoa(r)

	fileList := strings.Split(_file.Filename, ".")

	key := "/upload/" + ns + rs[0:6] + "." + fileList[1]
	if err != nil {
		print(err.Error(), 123)
		return
	}

	file, err3 := _file.Open()
	if err3 != nil {
		print(err3.Error())
		return
	}
	logger.Info("开始上传图片资源")
	_, err2 := cosClient.CosClient.Object.Put(context.Background(), key, file, &cos.ObjectPutOptions{})

	if err2 != nil {
		logger.Error(err2.Error())
		c.JSONOK(ErrInternalServer)
		return
	}
	c.JSONOK(config.GetConfig().CosConfig.CosUrl + key)
}
