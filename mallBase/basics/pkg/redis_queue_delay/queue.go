package redis_queue_delay

import (
	"MyTestMall/mallBase/basics/pkg/dredis"
	"context"
	"sync"
)

const (
	HashSuffix = ":hash"
	SetSuffix  = ":set"
)

var once sync.Once

type Options struct {
	Topic    string
	Handler  handlerFunc
	Duration int
}

type Queue struct {
	ctx      context.Context
	driver   *dredis.Driver
	topic    string
	producer *producer
	consumer *consumer
}

func NewQueue(ctx context.Context, client dredis.Driver, opts Options) *Queue {
	var queue *Queue
	once.Do(func() {
		queue = &Queue{
			ctx:      ctx,
			driver:   &client,
			topic:    opts.Topic,
			producer: NewProducer(ctx),
			consumer: NewConsumer(ctx, opts.Duration, opts.Handler),
		}
	})
	return queue
}

func (q *Queue) Start() {
	go q.consumer.listen(*q.driver, q.topic)
}

func (q *Queue) Publish(msg *Message) (int64, error) {
	return q.producer.Publish(*q.driver, q.topic, msg)
}
