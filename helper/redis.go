package helper

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisInstance *redis.Client

func GetRedisInstance() *redis.Client {
	if RedisInstance == nil {
		if os.Getenv("REDIS_URL") == "" {
			panic("Redis config error")
		}
		RedisInstance = redis.NewClient(&redis.Options{
			Addr:        os.Getenv("REDIS_URL"),
			PoolSize:    20,
			PoolTimeout: 15,
		})
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
