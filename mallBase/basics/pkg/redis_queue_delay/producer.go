package redis_queue_delay

import (
	"MyTestMall/mallBase/basics/pkg/dredis"
	"context"
	"strconv"
)

type producer struct {
	ctx context.Context
}

func NewProducer(ctx context.Context) *producer {
	return &producer{
		ctx: ctx,
	}
}

func (p *producer) Publish(driver dredis.Driver, topic string, msg *Message) (int64, error) {
	// stored sets 写入
	key := topic + SetSuffix
	n, err := driver.ZAdd(key, msg.GetScore(), strconv.FormatInt(msg.GetId(), 10))
	if err != nil {
		return n, err
	}
	// hashes 写入
	key = topic + HashSuffix
	return driver.HSet(key, msg.GetId(), msg)
}
