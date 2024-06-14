package wxpay

import (
	"context"
	"rentServer/pkg/config"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var WxPayClient *core.Client
var WxPayCtx context.Context

func Init() error {

	conf := config.GetConfig()
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(conf.WxPayConfig.CertPath)
	if err != nil {
		return err
	}

	WxPayCtx = context.Background()
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(conf.WxPayConfig.MchId, conf.WxPayConfig.MchCertificateSerialNumber, mchPrivateKey, conf.WxPayConfig.ApiV3Key),
	}
	WxPayClient, err = core.NewClient(WxPayCtx, opts...)
	if err != nil {
		return err
	}
	return nil

}
