package tasks

import (
	"github.com/go-redis/redis/v8"
	"github.com/gojuukaze/YTask/v2/brokers"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/util/yjson"
	"github.com/gojuukaze/YTask/v2/yerrors"
	"github.com/leffss/go-admin/pkg/gredis"
	"time"
)

type MyRedisBroker struct {
	BLPopTimeout time.Duration
}

func NewMyRedisBroker(blPopTimeout time.Duration) *MyRedisBroker {
	return &MyRedisBroker{
		BLPopTimeout: blPopTimeout,
	}
}

func (mr *MyRedisBroker) Activate() {}

func (mr *MyRedisBroker) SetPoolSize(n int) {}

func (mr *MyRedisBroker) GetPoolSize() int {
	return 1
}

func (mr *MyRedisBroker) Next(queueName string) (message.Message, error) {
	var msg message.Message
	values, err := gredis.BLPop(queueName, mr.BLPopTimeout)
	if err != nil {
		if err == redis.Nil {
			return msg, yerrors.ErrEmptyQuery{}
		}
		return msg, err
	}

	err = yjson.YJson.UnmarshalFromString(values[1], &msg)
	return msg, err
}

func (mr *MyRedisBroker) Send(queueName string, msg message.Message) error {
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = gredis.RPush(queueName, b)
	return err
}

func (mr *MyRedisBroker) LSend(queueName string, msg message.Message) error {
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = gredis.LPush(queueName, b)
	return err
}

func (mr *MyRedisBroker) Clone() brokers.BrokerInterface {
	return mr
}
