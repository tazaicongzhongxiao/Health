// Package redis_queue
// @Description: redis 消息队列
// https://github.com/toegg/redis_queue
package redis_queue

import (
	"MyTestMall/mallBase/basics/pkg/dredis"
	"errors"
	"fmt"
	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"log"
	"runtime"
	"sync"
)

type Queueable interface {
	Execute(*QueuePayload) *QueueResult
}

// QueuePayload 消息载体
type QueuePayload struct {
	ID     string      `json:"id"`
	IsFast bool        `json:"is_fast"`
	Topic  string      `json:"topic"`
	Group  string      `json:"group"`
	Body   interface{} `json:"body"`
}

// QueueResult 执行结果
type QueueResult struct {
	State   bool        `json:"state"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewQueueResult(state bool, msg string, data interface{}) *QueueResult {
	return &QueueResult{State: state, Message: msg, Data: data}
}

var instanceQueueManager *QueueManager
var onceQueueManager sync.Once

// QueueManager 队列管理器
type QueueManager struct {
	driver    dredis.Driver
	MaxRetry  int
	RecoverCh chan RecoverData
	Handlers  map[string]interface{}
}

// RecoverData 队列恢复的信息
type RecoverData struct {
	Topic   string
	Group   string
	Handler interface{}
}

// NewQueueManager 初始化队列管理器
func NewQueueManager() *QueueManager {
	onceQueueManager.Do(func() {
		instanceQueueManager = &QueueManager{}
		instanceQueueManager.MaxRetry = 3
		instanceQueueManager.Handlers = make(map[string]interface{})
	})
	return instanceQueueManager
}

func (r *QueueManager) SetRedis(client dredis.Driver) {
	r.driver = client
}

func (r *QueueManager) SetRecoverLis(ch chan RecoverData) {
	r.RecoverCh = ch
}

func (r *QueueManager) GetRecoverLis() chan RecoverData {
	return r.RecoverCh
}

func (r *QueueManager) GetQueueName(topic string, group string) string {
	var name string
	if len(group) > 0 {
		name = fmt.Sprintf("Queue_%s::%s", topic, group)
	} else {
		name = fmt.Sprintf("Queue_%s", topic)
	}
	return name
}

// RegisterQueue 注册队列
func (r *QueueManager) RegisterQueue(topic string, group string, handler interface{}) error {
	name := r.GetQueueName(topic, group)
	if _, ok := r.Handlers[name]; ok {
		return errors.New("is exits")
	} else {
		r.Handlers[name] = handler
		go r.QueueConsume(topic, group)
	}
	return nil
}

// RecoverQueue 重启队列
func (r *QueueManager) RecoverQueue(recoverData RecoverData) {
	name := r.GetQueueName(recoverData.Topic, recoverData.Group)
	if _, ok := r.Handlers[name]; ok {
		go r.QueueConsume(recoverData.Topic, recoverData.Group)
	}
}

// QueuePublish 生产者执行入队列
func (r *QueueManager) QueuePublish(payload *QueuePayload) error {
	if len(payload.Topic) <= 0 {
		return errors.New("TopicId can not be empty")
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	payload.ID = id.String()
	payloadStr, _ := json.Marshal(payload)
	_, _ = r.driver.LPush(r.GetQueueName(payload.Topic, payload.Group), payloadStr)
	return nil
}

// QueueConsume 消费者执行出队列
func (r *QueueManager) QueueConsume(topic string, group string) {
	defer func() {
		if err := recover(); err != nil {
			var stacktrace string
			for i := 1; ; i++ {
				_, f, l, got := runtime.Caller(i)
				if !got {
					break
				}
				stacktrace += fmt.Sprintf("%s:%d\n", f, l)
			}
			// when stack finishes
			logMessage := fmt.Sprintf("Trace: %s\n", err)
			logMessage += fmt.Sprintf("\n%s", stacktrace)
			log.Println(logMessage)
			//执行恢复函数
			r.handleRecover(topic, group)
		}
	}()
	for {
		//消费者执行出列
		var payload QueuePayload
		result, _ := r.driver.BRPop(0, r.GetQueueName(topic, group))
		if len(result) > 0 {
			err := json.Unmarshal([]byte(result[1]), &payload)
			if err != nil {
				log.Println("BRPOP json.Unmarshal Error:", err)
				continue
			}
			//执行回调函数
			r.handleCallBack(&payload)
		}
	}
}

// handleRecover 执行恢复函数
func (r *QueueManager) handleRecover(topic string, group string) {
	handleName := r.GetQueueName(topic, group)
	handler, ok := r.Handlers[handleName]
	if r.RecoverCh != nil && ok {
		r.RecoverCh <- RecoverData{topic, group, handler}
	}
}

// handleCallBack 执行回调函数
func (r *QueueManager) handleCallBack(payload *QueuePayload) {
	handleName := r.GetQueueName(payload.Topic, payload.Group)
	it := r.Handlers[handleName]
	if it != nil {
		if ob, ok := it.(Queueable); ok {
			//同步执行Max次，保证队列顺序，失败则丢弃消息,
			for i := 0; i < r.MaxRetry; i++ {
				rs := ob.Execute(payload)
				if rs.State {
					break
				}
			}
		} else {
			log.Println("no ExecuteFunc，pop：", payload)
		}
	}
}
