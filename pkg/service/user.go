package service

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	. "rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/middle_ware"
	"rentServer/pkg/model/user"
	"rentServer/pkg/model/wx_api"
	"rentServer/pkg/oauth"
)

type UserService interface {
	Service[user.User]
	WxMpLogin(dto wx_api.RequestDto) (token string, err error)
}

type userServiceIMPL struct {
	ServiceIMPL[user.User]
	UserDao dao.UserDao
}

func NewUserService() UserService {
	userDao := dao.NewUserDao()
	return &userServiceIMPL{
		ServiceIMPL: ServiceIMPL[user.User]{
			BaseDao: userDao,
		},
		UserDao: userDao,
	}
}

// WxMpLogin 微信小程序授权登陆
func (_this userServiceIMPL) WxMpLogin(dto wx_api.RequestDto) (token string, err error) {

	var userModel user.User
	profile, err := oauth.GetUserProfile(dto)
	if err != nil {
		return "", err
	}

	filter := map[string]any{
		"wx_open_id": profile.OpenID,
	}

	err = _this.UserDao.FilterFindOne(&userModel, filter)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 注册用户
		wxMpData, err := json.Marshal(&profile)
		if err != nil {
			return "", err
		}
		userModel.Nickname = profile.NickName
		userModel.Avatar = profile.AvatarURL
		userModel.WxMpData = string(wxMpData)
		userModel.WxOpenID = profile.OpenID

		err = _this.UserDao.Insert(&userModel)
		if err != nil {
			return "", err
		}
	}

	if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return "", err
	}

	token, err = middle_ware.GenerateToken(userModel.ID)
	if err != nil {
		return "", err
	}

	return
}
