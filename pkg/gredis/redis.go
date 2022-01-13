package gredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/leffss/go-admin/pkg/setting"
	"time"
)

var RedisConn redis.UniversalClient
var ctx = context.Background()

// Setup Initialize the Redis instance
func Setup() error {
	redisSetting := setting.GetRedisSetting()
	RedisConn = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: redisSetting.Host,
		MasterName: redisSetting.MasterName,
		Password: redisSetting.Password,
		IdleTimeout: redisSetting.IdleTimeout,
		PoolSize: redisSetting.PoolSize,
	})

	return nil
}

func CloseRedis() error {
	return RedisConn.Close()
}

// Set a key/value
func Set(key string, data interface{}, expire time.Duration) error {
	return RedisConn.Set(ctx, key, data, expire).Err()
}

// Get get a key
func Get(key string) (string, error) {
	return RedisConn.Get(ctx, key).Result()
}

// Delete delete a kye
func Delete(key string) error {
	return RedisConn.Del(ctx, key).Err()
}

func RPush(key string, data interface{}) error {
	return RedisConn.RPush(ctx, key, data).Err()
}

func LPush(key string, data interface{}) error {
	return RedisConn.LPush(ctx, key, data).Err()
}

func BLPop(key string, timeout time.Duration) ([]string, error) {
	return RedisConn.BLPop(ctx, timeout, key).Result()
}
