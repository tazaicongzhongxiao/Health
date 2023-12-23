package tarsgo

import (
	"MyTestMall/mallBase/basics/pkg/log"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/protocol/push"
	json "github.com/json-iterator/go"
)

var (
	comm1 = tars.NewCommunicator()
	proxy map[string]*push.Client
)

// NoticeBody
// @Description: 在线消息
type NoticeBody struct {
	StoreId       int64    `json:"store_id"`                  // 店铺ID
	ReceiveUserId []string `json:"receive_user_id,omitempty"` // 接收对象
	Mark          string   `json:"mark"`                      // 页面标识
	AppData       string   `json:"app_data"`                  // 应用数据
	SendDate      int64    `json:"send_date"`                 // 发送时间
}

// NoticeResult
// @Description: 在线消息结果
type NoticeResult struct {
	Code int32  `json:"code"` // 状态，200：成功
	Msg  string `json:"msg"`  // 错误消息
	//Data map[string]string `json:"data"` // 返回内容
}

type Client struct {
	client *push.Client
}

func init() {
	proxy = make(map[string]*push.Client)
}

// PushConnect
// @Description: 建立PUSH连接
// @param address
// @param callback
// @return client
func PushConnect(address string, callback ...func(d []byte)) *Client {
	if proxy[address] == nil {
		fun := func(d []byte) {
			log.Info("plugin.tarsgo.notice.PushNotice.data", string(d))
		}
		if len(callback) > 0 {
			fun = callback[0]
		}
		proxy[address] = push.NewClient(fun)
		comm1.StringToProxy(address, proxy[address])
	}
	return &Client{client: proxy[address]}
}

// PushClose
// @Description: 销毁map中的push地址
// @param address
func PushClose(address string) {
	delete(proxy, address)
	return
}

// Send
// @Description: 推送消息
func (c *Client) Send(b []byte) ([]byte, error) {
	return c.client.Connect(b)
}

// Notice
// @Description: 推送在线通知消息
func (c *Client) Notice(req NoticeBody) (res NoticeResult) {
	if req.Mark == "" {
		return NoticeResult{Code: 2, Msg: "页面标识为空"}
	}
	b, err := json.Marshal(req)
	if err != nil {
		return NoticeResult{Code: 2, Msg: err.Error()}
	}
	if data, errs := c.client.Connect(b); errs != nil {
		res.Msg = errs.Error()
	} else if errs = json.Unmarshal(data, &res); errs != nil {
		res.Msg = errs.Error()
	}
	return
}

//func PushNotice(req NoticeBody, _address string) (res NoticeResult) {
//	return PushConnect(_address).Notice(req)
//}
