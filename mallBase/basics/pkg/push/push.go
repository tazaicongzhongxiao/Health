package push

import (
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/protocol/push"
)

var (
	comm  = tars.GetCommunicator()
	proxy map[string]*push.Client
)

func init() {
	proxy = make(map[string]*push.Client)
}

// Connect
// @Description: 建立PUSH连接
// @param address
// @param callback
// @return *Client
func Connect(address string, callback func(data []byte)) *push.Client {
	if proxy[address] == nil {
		proxy[address] = push.NewClient(callback)
		comm.StringToProxy(address, proxy[address])
	}
	return proxy[address]
}

// ConnectSend
// @Description: 建立PUSH连接并发送数据
// @param address
// @param data
// @return err
func ConnectSend(address string, data []byte) (err error) {
	if proxy[address] == nil {
		proxy[address] = push.NewClient(nil)
		comm.StringToProxy(address, proxy[address])
	}
	if proxy[address] != nil {
		_, err = proxy[address].Connect(data)
	}
	return fmt.Errorf("PUSH连接未正确获取")
}

func Close(address string) {
	delete(proxy, address)
	return
}
