package database

import (
	"github.com/redis/go-redis/v9"
	"time"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "", // 没有密码，默认值
		DB:           0,  // 默认DB 0
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
}

func GetRedis() *redis.Client {
	return RedisClient
}
