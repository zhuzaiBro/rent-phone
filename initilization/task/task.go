package task

import "rentServer/pkg/service"

// 用于检查并执行一些由于服务端程序关闭而取消的定时任务

func Init() error {
	orderServ := service.NewOrderService()

	return orderServ.OrderTask()
}
