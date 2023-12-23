package queue

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/pkg/config"
	"MyTestMall/mallBase/basics/pkg/log"
	"MyTestMall/mallBase/basics/tools/dstring"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/go-playground/validator/v10"
	"github.com/nsqio/go-nsq"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type configs struct {
	Address     string `mapstructure:"address"`
	Lookup      string `mapstructure:"lookup"`
	MaxInFlight int    `mapstructure:"max_in_flight"`
}

type NsqHandler struct {
	topic  string
	handle func(*Message) (map[string]string, error)
	failed func(string, *Message, error, int32) error // state 1待处理 2重试 3成功 4失败
}

type Message struct {
	Name      string      `json:"name,omitempty"`       // 消息标识名称
	Body      interface{} `json:"body,omitempty"`       // 队列消息体
	CallTopic string      `json:"call_topic,omitempty"` // 消息结果回调 name 固定为 nsq_callback
	CallId    string      `json:"call_id,omitempty"`    // 消息结果回调识别ID
	LogId     int64       `json:"log_id,omitempty"`     // NSQ日志回调ID
	MaxRetry  uint16      `json:"max_retry,omitempty"`  // 最大重试次数
	AutoRetry int32       `json:"auto_retry,omitempty"` // 错误后间隔多久重试1次 (默认5秒重试一次，最大自动重试5次)
}

// MessageResult
// @Description: 回调结果
type MessageResult struct {
	Code    int8              `json:"code,omitempty"`     // 状态 0成功 1发送失败 2处理失败
	Msg     string            `json:"msg,omitempty"`      // 消息内容(错误提示)
	Name    string            `json:"name,omitempty"`     // 消息标识名称
	ErrCode int32             `json:"err_code,omitempty"` // 错误返回 code
	CallId  string            `json:"call_id,omitempty"`  // 消息结果回调识别ID
	Params  map[string]string `json:"params,omitempty"`   // 回调结果参数在成功处理消息可通知给发送方数据
}

var (
	NsqLog    = "mall.server.base"
	validate  *validator.Validate
	conf      configs
	nsqClient *NsqDriver
)

type NsqDriver struct {
	addrNsqLookups []string
	producer       *nsq.Producer
	logLevel       nsq.LogLevel
	consumers      []*nsq.Consumer
	nsqConfig      *nsq.Config
}

func Start() bool {
	validate = validator.New()
	var err error
	err = config.Config().Bind(app.CfgName, "nsq", &conf, func() {
		if config.Config().Bind(app.CfgName, "nsq", &conf, nil) == nil {
			if nsqClient, err = initStart(); err != nil {
				panic(err)
			}
		}
	})
	if err == nil {
		nsqClient, err = initStart()
	}
	if err != nil {
		log.Error("nsq 错误", err.Error())
		panic(err)
	}
	return err == nil
}

func initStart() (c *NsqDriver, err error) {
	if conf.Address == "" {
		conf.Address = "127.0.0.1:4150"
	}
	c = &NsqDriver{
		addrNsqLookups: strings.Split(conf.Lookup, ","),
		nsqConfig:      nsq.NewConfig(),
	}
	c.nsqConfig.MaxRequeueDelay = time.Minute * 60
	c.nsqConfig.MaxInFlight = conf.MaxInFlight
	c.producer, err = nsq.NewProducer(conf.Address, c.nsqConfig)
	if err != nil {
		log.Warn("nsq connect error", err.Error())
		return c, err
	}
	//logLevel = nsq.LogLevelDebug
	c.logLevel = nsq.LogLevelError
	c.producer.SetLogger(log.NsqLogger(), c.logLevel)
	if err = c.producer.Ping(); err != nil {
		log.Warn("nsq ping error", err.Error())
		return c, err
	}
	log.Info("nsq connect success", conf.Address)
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		c.gracefulStop()
	}()
	return c, nil
}

// Client
// @Description: 获取 NSQ Producer
// @return error
func Client() *NsqDriver {
	return nsqClient
}

// PublishBody
// @Description: 检测是否插件并进行转换
// @param topic 回调地址
// @param body
// @return err
func PublishBody(topic string, m *Message) (newTopic string, err error) {
	pluginId, _ := strconv.ParseFloat(topic, 64)
	if pluginId > 0 {
		Body := struct {
			PluginID    int32  `json:"plugin_id"`    // 使用的插件ID
			PluginLabel string `json:"plugin_label"` // 使用的插件接口标签名
			Data        any    `json:"data"`         // 传递的数据
		}{
			int32(pluginId),
			m.Name,
			m.Body,
		}
		m.Name = "plugin"
		m.Body = Body
		topic = "mall.server.plugin"
	} else if !dstring.IsNSQ(topic) {
		err = errors.New(fmt.Sprintf("NSQ 地址[%s]非法", topic))
	}
	return topic, err
}

// Client
// @Description: 获取客户端
// @return error
func (rd *NsqDriver) Client() *nsq.Producer {
	return rd.producer
}

// ClientID
// @Description: 获取 服务器标识
// @return error
func (rd *NsqDriver) ClientID() string {
	return rd.nsqConfig.ClientID
}

// Ping
// @Description: 检查 NSQ 状态
// @return error
func (rd *NsqDriver) Ping() error {
	return rd.producer.Ping()
}

// Stop
// @Description: 停止NSQ
func (rd *NsqDriver) Stop() {
	rd.gracefulStop()
}

// PublishCallback
// @Description: NSQ结果通知
// @param topic 回调地址
// @param body
// @return err
func PublishCallback(topic string, LogId int64, body MessageResult) {
	if topic != "" {
		data := &Message{
			Name: "nsq_callback",
			Body: body,
		}
		newTopic, err := PublishBody(topic, data)
		if err == nil {
			jsons, errs := json.Marshal(data)
			if errs == nil {
				err = nsqClient.producer.Publish(newTopic, jsons)
			}
		}
		if err != nil {
			log.Error(fmt.Sprintf("NSQ结果通知 %s", err.Error()), body, 3)
		}
	}
	if LogId > 0 {
		body.Name = "nsq_result"
		body.CallId = strconv.FormatInt(LogId, 10)
		data := &Message{
			Name: "nsq_callback",
			Body: body,
		}
		jsons, errs := json.Marshal(data)
		if errs == nil {
			_ = nsqClient.producer.Publish(NsqLog, jsons)
		}
	}
	return
}

func Publish(topic string, body Message, delay ...int32) (err error) {
	if nsqClient == nil {
		return errors.New("NSQ 客户端未初始化连接")
	}
	defer func() {
		if err != nil {
			if body.CallTopic != "" || body.LogId > 0 {
				go PublishCallback(body.CallTopic, body.LogId, MessageResult{
					Code:   1,
					Msg:    err.Error(),
					Name:   body.Name,
					CallId: body.CallId,
				})
			}
			log.Info(fmt.Sprintf("NSQ发送(%s)ERR:%s", topic, err), body, 3)
		} else {
			log.Info(fmt.Sprintf("NSQ发送(%s)OK", topic), body, 3)
		}
	}()
	if nsqClient.producer == nil {
		return fmt.Errorf("NSQ 未连接 topic: %s", topic)
	}
	if topic == "" {
		return fmt.Errorf("topic 未传入")
	}
	if body.Body != nil && reflect.TypeOf(body.Body).Kind() == reflect.Struct {
		err = validate.Struct(body.Body)
		if err != nil {
			return err
		}
	}
	jsons, err := json.Marshal(body)
	if err != nil {
		return err
	}
	for {
		retryCount := 0
	retry:
		if len(delay) == 1 && delay[0] > 0 {
			err = nsqClient.producer.DeferredPublish(topic, time.Duration(delay[0])*time.Second, jsons)
		} else {
			err = nsqClient.producer.Publish(topic, jsons)
		}
		if err != nil && tars.GetErrorCode(err) == 1 {
			if retryCount < 2 {
				retryCount++
				goto retry
			}
		} else {
			break
		}
	}
	return err
}

// Subscribe
// @Description: 创建消费者
// @param topic
// @param channel
// @param handler
// @return err
func Subscribe(topic string, channel string, handle func(*Message) (map[string]string, error), failed func(string, *Message, error, int32) error) (err error) {
	if nsqClient == nil {
		return errors.New("NSQ 客户端未初始化连接")
	}
	h := new(NsqHandler)
	nsqClient.nsqConfig.MaxAttempts = 0
	nsqClient.nsqConfig.DefaultRequeueDelay = 0
	nsqClient.nsqConfig.MaxBackoffDuration = time.Millisecond * 50
	c, err := nsq.NewConsumer(topic, channel, nsqClient.nsqConfig)
	if err != nil {
		log.Error("nsq err", err.Error())
		return err
	}
	c.SetLogger(log.NsqLogger(), nsqClient.logLevel)
	c.AddConcurrentHandlers(h, 1)
	err = c.ConnectToNSQLookupds(nsqClient.addrNsqLookups)
	if err != nil {
		log.Error("nsq err", err.Error())
		return err
	}
	h.handle = handle
	h.failed = failed
	h.topic = topic
	// 启动后发送测试信息给自己
	jsons, err := json.Marshal(Message{
		Name: "test",
	})
	_ = nsqClient.producer.Publish(topic, jsons)
	nsqClient.consumers = append(nsqClient.consumers, c)
	return
}

func (rd *NsqDriver) gracefulStop() {
	rd.producer.Stop()
	var wg sync.WaitGroup
	for _, c := range rd.consumers {
		wg.Add(1)
		go func(c *nsq.Consumer) {
			c.Stop()
			// disconnect from all lookupd
			for _, addr := range rd.addrNsqLookups {
				err := c.DisconnectFromNSQLookupd(addr)
				if err != nil {
					log.Warn("nsq lookupd", err.Error())
				}
			}
			wg.Done()
		}(c)
	}
	wg.Wait()
}

func (s *NsqHandler) HandleMessage(message *nsq.Message) (err error) {
	var msg *Message
	err = json.Unmarshal(message.Body, &msg)
	if err != nil {
		log.Error("nsq 解析错误 "+err.Error(), string(message.Body), 3)
		return nil
	}
	if msg.LogId > 0 && (s.topic == NsqLog || s.topic == "") {
		msg.LogId = 0
	}
	res, err := s.handle(msg)
	if err != nil {
		log.Info(fmt.Sprintf("NSQ接收(%d)ERR:%s", message.Attempts, err), msg, 3)
		if tars.GetErrorCode(err) == 1 {
			message.DisableAutoResponse()
			if msg.MaxRetry > 0 && msg.MaxRetry > message.Attempts {
				if msg.AutoRetry > 0 {
					message.Requeue(time.Duration(msg.AutoRetry) * time.Second)
				} else {
					message.Requeue(-1)
				}
				return nil
			} else {
				if msg.CallTopic != "" || msg.LogId > 0 {
					go PublishCallback(msg.CallTopic, msg.LogId, MessageResult{
						Code:    2,
						Msg:     err.Error(),
						ErrCode: tars.GetErrorCode(err),
						Name:    msg.Name,
						CallId:  msg.CallId,
					})
				}
				if msg.LogId == 0 {
					err = s.failed(s.topic, msg, err, 4)
				}
			}
		} else {
			if msg.CallTopic != "" || msg.LogId > 0 {
				go PublishCallback(msg.CallTopic, msg.LogId, MessageResult{
					Code:    2,
					Msg:     err.Error(),
					ErrCode: tars.GetErrorCode(err),
					Name:    msg.Name,
					CallId:  msg.CallId,
				})
			}
			if msg.LogId == 0 {
				_ = s.failed(s.topic, msg, err, 4)
			}
			err = nil
		}
	} else {
		log.Info(fmt.Sprintf("NSQ接收(%d)OK", message.Attempts), msg, 3)
		if msg.CallTopic != "" || msg.LogId > 0 {
			go PublishCallback(msg.CallTopic, msg.LogId, MessageResult{
				Code:   0,
				Name:   msg.Name,
				CallId: msg.CallId,
				Params: res,
			})
		}
	}
	message.Finish()
	return err
}
