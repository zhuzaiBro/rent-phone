package initilization

import (
	"rentServer/initilization/db"
	"rentServer/initilization/http"
	"rentServer/initilization/resource"
	"rentServer/initilization/service_time"
	"rentServer/initilization/task"
	"rentServer/initilization/wxpay"
	"rentServer/pkg/config"
	"rentServer/pkg/logger"
)

func LoadConf() {

	config.LoadConfig()
	logger.InitLogger(config.GetConfig().LogLevel)
	db.InitDB()

	service_time.Init()
	//err := rabit.Init()
	//if err != nil {
	//	panic(err)
	//}

	err := task.Init()
	if err != nil {
		panic(err)
	}
	err = resource.Init()
	if err != nil {
		panic(err)
	}
	err = wxpay.Init()
	if err != nil {
		panic(err)
	}
	err = http.NewRouter().Run(":" + config.GetConfig().HttpPort)
	if err != nil {
		panic(err)
	}

}
