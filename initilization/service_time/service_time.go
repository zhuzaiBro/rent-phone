package service_time

import (
	"fmt"
	"math"
	"time"
)

var TimeMap map[string]time.Time

func Init() {
	// 创建一个映射
	timeMap := make(map[string]time.Time)

	// 获取当前日期
	now := time.Now()
	year, month, day := now.Date()

	// 构建时间字符串和时间戳的映射关系
	for i := 0; i < 96; i++ {
		// 15分钟一个模块
		hour := 0.25 * float64(i)

		// 构建时间字符串
		timeStr := fmt.Sprintf("%02d:%02d", int(math.Floor(hour)), (i%4)*15)

		// 构建时间戳
		timestamp := time.Date(year, month, day, i, 0, 0, 0, time.Local)

		// 将时间字符串和时间戳添加到映射中
		timeMap[timeStr] = timestamp
	}

	TimeMap = timeMap
}
