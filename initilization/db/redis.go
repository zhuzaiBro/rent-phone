package db

import (
	"fmt"
	redis "github.com/redis/go-redis/v9"
	"rentServer/pkg/config"
)

var redisClient *redis.Client

func ConnectRedis() {
	cfg := config.GetConfig()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	})
}

func GetRedis() *redis.Client {
	return redisClient
}
