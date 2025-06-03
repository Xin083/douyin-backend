package model

import (
	"context"
	"douyin-backend/app/global/variable"
	"sync"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	dbContainer = make(map[string]*gorm.DB)
	redisClient *redis.Client
	redisOnce   sync.Once
)

// GetRedisClient 获取 Redis 客户端实例
func GetRedisClient() *redis.Client {
	redisOnce.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     variable.ConfigYml.GetString("Redis.Host") + ":" + variable.ConfigYml.GetString("Redis.Port"),
			Password: variable.ConfigYml.GetString("Redis.Auth"), // 使用 Auth 而不是 Password
			DB:       0,                                          // 使用默认 DB
		})

		// 测试连接
		ctx := context.Background()
		if err := redisClient.Ping(ctx).Err(); err != nil {
			panic("Redis connection failed: " + err.Error())
		}
	})
	return redisClient
}

// ... 其他现有代码 ...
