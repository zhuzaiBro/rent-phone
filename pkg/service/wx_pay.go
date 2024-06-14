package service

import (
	"encoding/json"
	wx_pay_utils "github.com/wechatpay-apiv3/wechatpay-go/utils"
	"rentServer/pkg/config"
	"rentServer/pkg/model/wx_api"
)

type WxPayService interface {
	DecodeNotify(notify *wx_api.WechatPayNotify, wxCallback *wx_api.WxPayCallback) error
}

type wxPayServiceIMPL struct {
}

func NewWxPayService() WxPayService {
	return wxPayServiceIMPL{}
}

// DecodeNotify 解密微信回调密文
func (_this wxPayServiceIMPL) DecodeNotify(notify *wx_api.WechatPayNotify, wxCallback *wx_api.WxPayCallback) error {
	conf := config.GetConfig()
	pt, err := wx_pay_utils.DecryptAES256GCM(conf.WxPayConfig.ApiV3Key, notify.Resource.AssociatedData, notify.Resource.Nonce, notify.Resource.Cipher)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(pt), &wxCallback)

}
