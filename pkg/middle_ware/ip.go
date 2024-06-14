package middle_ware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

var ipCountsMutex sync.Mutex // 全局互斥锁，保护对ipCounts的访问

// IPRequestCounter 存储IP地址及其最后一次访问时间和访问次数
type IPRequestCounter struct {
	sync.Mutex
	Count     int
	LastVisit time.Time
}

// Blacklist 存储被禁止的IP地址
var Blacklist = make(map[string]bool)

// IPMiddleware 是一个中间件，用于限制IP的访问频率
func IPMiddleware() gin.HandlerFunc {
	// 创建一个map来存储IP地址及其访问计数和最后访问时间
	// 这里我们使用全局变量作为示例，但在实际应用中，你可能希望使用更持久的存储（如Redis）
	ipCounts := make(map[string]*IPRequestCounter)
	// 设置时间间隔，例如每分钟检查一次
	interval := 1 * time.Minute

	// 启动一个goroutine来定期清理过期的IP计数
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			ipCountsMutex.Lock()
			for ip, counter := range ipCounts {
				if time.Now().Sub(counter.LastVisit) > interval {
					counter.Count = 0 // 重置计数器
					if _, ok := Blacklist[ip]; ok {
						// 如果IP在黑名单中但已经过期，可以将其从黑名单中移除（可选）
						delete(Blacklist, ip)
					}
				}
			}
			ipCountsMutex.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		// 检查IP是否在黑名单中
		if Blacklist[ip] {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests from this IP, access denied."})
			c.Abort()
			return
		}

		ipCountsMutex.Lock()
		counter, ok := ipCounts[ip]
		if !ok {
			counter = &IPRequestCounter{}
			ipCounts[ip] = counter
		}
		counter.Lock() // 对单个IP的counter进行加锁
		counter.Count++
		counter.LastVisit = time.Now()
		counter.Unlock() // 解锁单个IP的counter
		ipCountsMutex.Unlock()

		// 检查是否超过限制（例如每分钟最多100次请求）
		if counter.Count > 100 { // 这里你可以设置你自己的限制
			Blacklist[ip] = true // 将IP添加到黑名单中
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests from this IP, temporary access denied."})
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}
