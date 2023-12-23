package redis_queue_delay

import (
	"MyTestMall/mallBase/basics/pkg/dredis"
	"MyTestMall/mallBase/basics/pkg/log"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type handlerFunc func(msg Message)

type consumer struct {
	ctx      context.Context
	duration time.Duration
	ch       chan []string
	handler  handlerFunc
}

func NewConsumer(ctx context.Context, duration int, handler handlerFunc) *consumer {
	if duration == 0 {
		duration = 1
	}
	return &consumer{
		ctx:      ctx,
		duration: time.Duration(duration) * time.Second, // 每隔多少秒获取一次数据
		ch:       make(chan []string, 1000),
		handler:  handler,
	}
}

func (c *consumer) listen(driver dredis.Driver, topic string) {
	// 从 Hashes 中获取数据并处理
	go func() {
		for {
			select {
			case ret := <-c.ch:
				// 批量从hashes中获取数据信息
				key := topic + HashSuffix
				result, err := driver.HMGet(key, ret...)
				if err != nil {
					log.Error("redis_queue_delay", err.Error())
				}
				if len(result) > 0 {
					_, _ = driver.HDel(key, ret...)
				}
				msg := Message{}
				for _, v := range result {
					// 由于hashes 和 scoreSet 非事务操作，会出现删除了set但hashes未删除的情况
					if v == nil {
						continue
					}
					str := v.(string)
					_ = json.Unmarshal([]byte(str), &msg)
					// 处理逻辑
					go c.handler(msg)
				}
			}
		}
	}()
	ticker := time.NewTicker(c.duration)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			log.Error("redis_queue_delay consumer quit", c.ctx.Err())
			return
		case <-ticker.C:
			// read data from redis
			min := strconv.Itoa(0)
			max := strconv.Itoa(int(time.Now().Unix()))
			opt := &redis.ZRangeBy{
				Min: min,
				Max: max,
			}
			key := topic + SetSuffix
			result, err := driver.ZRangeByScore(key, opt)
			if err != nil {
				log.Error("redis_queue_delay read", err.Error())
				return
			}
			// 获取到数据
			if len(result) > 0 {
				// 从 sorted sets 中移除数据
				_, _ = driver.ZRemRangeByScore(key, min, max)
				// 写入 chan, 进行hashes处理
				c.ch <- result
			}
		}
	}
}
