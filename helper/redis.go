package helper

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var RedisInstance *redis.Client

func GetRedisInstance() *redis.Client {
	if RedisInstance == nil {
		redisURL := viper.GetString("REDIS_URL")
		if redisURL == "" {
			panic("Redis config error")
		}
		opt, error := redis.ParseURL(redisURL)
		if error != nil {
			panic(error)
		}
		opt.PoolSize = 20
		opt.PoolTimeout = 15
		ctx := context.Background()
		if _, err := RedisInstance.Ping(ctx).Result(); err != nil {
			panic(err)
		}
	}
	return RedisInstance
}

func RunRedis(f func(*redis.Conn) error) error {
	conn := GetRedisInstance().Conn()
	defer conn.Close()
	err := f(conn)
	return err
}
