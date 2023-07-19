package helper

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisInstance *redis.Client

func init() {
	RedisInstance = redis.NewClient(&redis.Options{
		Addr:        os.Getenv("REDIS_URL"),
		PoolSize:    20,
		PoolTimeout: 15,
	})
}

func GetRedisInstance() *redis.Client {
	return RedisInstance
}

func RunRedis(f func(*redis.Conn) error) error {
	conn := RedisInstance.Conn()
	defer conn.Close()
	err := f(conn)
	return err
}
