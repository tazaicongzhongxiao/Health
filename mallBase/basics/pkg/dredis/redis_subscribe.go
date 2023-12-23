package dredis

import (
	json "github.com/json-iterator/go"
	"reflect"
)

// Subscribe REDIS 订阅消息
func Subscribe(rd Driver, channels []string, run func(channel, payload string)) {
	pubSub := rd.Subscribe(channels...)
	defer pubSub.Close()
	for msg := range pubSub.Channel() {
		go run(msg.Channel, msg.Payload)
	}
	return
}

// Publish REDIS 发布订阅消息
func Publish(rd Driver, channel string, message interface{}) (int64, error) {
	switch reflect.TypeOf(message).Kind() {
	case reflect.String:
		return rd.Publish(channel, message.(string))
	default:
		bs, _ := json.Marshal(message)
		return rd.Publish(channel, string(bs))
	}
}
