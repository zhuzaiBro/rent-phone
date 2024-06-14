package dao

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache[T interface{}] interface {
	CacheSave(data *T, key string, duration time.Duration) error
	CacheGetValue(data *T, key string) error
	CacheChangeValue(key string, data *T) error
}

type RedisCacheIMPL[T interface{}] struct {
	Client *redis.Client
}

func (_this RedisCacheIMPL[T]) CacheSave(data *T, key string, duration time.Duration) error {

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return _this.Client.Set(context.Background(), key, value, duration).Err()
}

func (_this RedisCacheIMPL[T]) CacheGetValue(data *T, key string) error {

	record, err := _this.Client.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(record), data)
}

func (_this RedisCacheIMPL[T]) GetStringValue(key string) (str string, err error) {
	return _this.Client.Get(context.Background(), key).Result()
}

func (_this RedisCacheIMPL[T]) CacheChangeValue(key string, data *T) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// 获取键的剩余过期时间
	ttlResult, err := _this.Client.TTL(context.Background(), key).Result()
	if err != nil {
		return err
	}

	// 获取剩余过期时间的秒数
	ttl := ttlResult.Seconds()
	if ttl <= 0 {
		return errors.New("key已过期")
	}
	return _this.Client.Set(context.Background(), key, string(value), ttlResult).Err()
}
