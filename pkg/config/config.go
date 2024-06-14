package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Credential struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

// Configuration 项目配置
type Configuration struct {
	// 日志级别，info或者debug
	LogLevel string `yaml:"log_level"`

	JwtSecret string `yaml:"jwt_secret"`

	HttpPort string `yaml:"http_port"`

	MysqlConf struct {
		Dsn          string `yaml:"dsn"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxOpenConns int    `yaml:"maxOpenConns"`
	} `yaml:"mysql"`
	ServerUrl string `yaml:"server_url"`

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
	} `yaml:"redis"`

	WechatAppConfiguration struct {
		AppId     string `yaml:"app_id"`
		AppSecret string `yaml:"app_secret"`
	} `yaml:"wechat_app_config"`

	CosConfig struct {
		CosUrl    string `yaml:"cos_url"`
		SecretID  string `yaml:"secret_id"`
		SecretKey string `yaml:"secret_key"`
	} `yaml:"cos_config"`

	WxPayConfig struct {
		MchId                      string `yaml:"mch_id"`
		MchCertificateSerialNumber string `yaml:"mch_certificate_serial_number"`
		ApiV3Key                   string `yaml:"api_v_3_key"`
		CertPath                   string `yaml:"cert_path"`
		NotifyUrl                  string `yaml:"notify_url"`
		Description                string `yaml:"description"`
		Currency                   string `yaml:"currency"`
	} `yaml:"wx_pay_config"`

	WxappImg struct {
		Path   string `yaml:"path"`
		Domain string `yaml:"domain"`
	} `yaml:"wxapp_img"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{}
		data, err := ioutil.ReadFile("config.yml")
		if err != nil {
			log.Fatal(err)
		}
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			log.Fatal(err)
		}

		// 如果环境变量有配置，读取环境变量
		logLevel := os.Getenv("LOG_LEVEL")
		if logLevel != "" {
			config.LogLevel = logLevel
		}

	})

	// 一些默认值
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	return config
}

func GetConfig() *Configuration {
	return config
}
