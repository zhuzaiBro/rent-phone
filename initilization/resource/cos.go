package resource

import (
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"rentServer/pkg/config"
)

var CosClient *cos.Client

func Init() error {

	conf := config.GetConfig()
	u, _ := url.Parse(conf.CosConfig.CosUrl)

	b := &cos.BaseURL{BucketURL: u}
	CosClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: os.Getenv(conf.CosConfig.SecretID),
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: os.Getenv(conf.CosConfig.SecretKey),
		},
	})

	return nil
}
