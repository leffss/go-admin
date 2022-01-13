package tasks

import (
	"github.com/go-redis/redis/v8"
	"github.com/gojuukaze/YTask/v2/backends"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/util/yjson"
	"github.com/gojuukaze/YTask/v2/yerrors"
	"github.com/leffss/go-admin/pkg/gredis"
	"time"
)

type MyRedisBackend struct {
}

func NewMyRedisBackend() *MyRedisBackend {
	return &MyRedisBackend{
	}
}

func (br *MyRedisBackend) Activate() {}

func (br *MyRedisBackend) SetPoolSize(n int) {}

func (br *MyRedisBackend) GetPoolSize() int {
	return 1
}

func (br *MyRedisBackend) SetResult(result message.Result, exTime int) error {
	b, err := yjson.YJson.Marshal(result)

	if err != nil {
		return err
	}
	err = gredis.Set(result.GetBackendKey(), b, time.Duration(exTime) * time.Second)
	return err
}

func (br *MyRedisBackend) GetResult(key string) (message.Result, error) {
	var result message.Result

	b, err := gredis.Get(key)
	if err != nil {
		if err == redis.Nil {
			return result, yerrors.ErrNilResult{}
		}
		return result, err
	}

	err = yjson.YJson.Unmarshal([]byte(b), &result)
	return result, err
}

func (br *MyRedisBackend) Clone() backends.BackendInterface{
	return br
}
