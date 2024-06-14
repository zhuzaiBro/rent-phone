package service

import (
	"fmt"
	"math"
	"rentServer/initilization/service_time"
	"rentServer/pkg/dao"
	"rentServer/pkg/model/order"
	"rentServer/pkg/model/reservation"
	"time"
)

var _ ReservationService = (*reservationService)(nil)

//
//TimeMap := map[time.D]string {
//		time.
//}

type ReservationService interface {
	Handlerent(rent *reservation.Reservation, in *order.Order, userID uint64) error
	GetServiceTime(timestamp int64) (list map[string]time.Time)
}

type reservationService struct {
	ClerkDao dao.ClerkDao
	OrderService
}

func NewrentService() ReservationService {
	return &reservationService{
		ClerkDao:     dao.NewClerkDao(),
		OrderService: NewOrderService(),
	}
}

func (_this reservationService) Handlerent(rent *reservation.Reservation, in *order.Order, userID uint64) error {

	// 通过师傅来查询价格
	err := _this.ClerkDao.FindByID(&rent.Clerk, rent.Clerk.ID)
	if err != nil {
		return err
	}
	// 生成服务订单
	err = _this.OrderService.GenerateServiceOrder(in, rent)
	if err != nil {
		return err
	}

	return nil
}

func (_this reservationService) GetWorkTime(clerkID uint64) {
	d := service_time.TimeMap["17:15"]

	d.Second()
	//time.
	//_this.
}

func IsTomorrowOrLater(timestamp time.Time) bool {
	now := time.Now()
	startOfTomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	return timestamp.After(startOfTomorrow)
}

func (_this reservationService) GetServiceTime(timestamp int64) (list map[string]time.Time) {

	timeMap := generateTime(timestamp)
	// 通过 timeStep 判断用户传进来的日期 先比较日期
	//if IsTomorrowOrLater(time.UnixMilli(timestamp)) {
	//	for key, _ := range timeMap {
	//		list = append(list, key)
	//	}
	//	return
	//}
	//
	//now := time.Now()
	//for key, item := range timeMap {
	//	if item.Second() > now.Second() {
	//		// 在这个时间之前要留下
	//		list = append(list, key)
	//	}
	//}
	return timeMap
}

func generateTime(_timestamp int64) map[string]time.Time {
	// 创建一个映射
	timeMap := make(map[string]time.Time)

	// 获取当前日期
	now := time.Now()
	year, month, day := now.Date()

	// 构建时间字符串和时间戳的映射关系
	for i := 0; i < 96; i++ {
		// 15分钟一个模块
		hour := 0.25 * float64(i)

		hours := int(math.Floor(hour))
		min := (i % 4) * 15

		// 构建时间字符串
		timeStr := fmt.Sprintf("%02d:%02d", hours, min)

		// 构建时间戳
		timestamp := time.Date(year, month, day, hours, min, 0, 0, time.Local)

		if timestamp.UnixMilli() > _timestamp {
			// 将时间字符串和时间戳添加到映射中
			timeMap[timeStr] = timestamp
		}

	}

	return timeMap
}
