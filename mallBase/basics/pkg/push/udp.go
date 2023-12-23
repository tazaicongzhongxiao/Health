package push

import (
	"MyTestMall/mallBase/basics/pkg/log"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	udpProxy = make(map[string]*net.UDPConn)
)

// UdpServer
// @Description: 建立UDP服务器连接
// @param address UDP地址 如 127.0.0.1:8080
// @param data
// @return err
func UdpServer(address string, port int, callback func(addr string, data []byte)) (err error) {
	if port == 0 {
		return errors.New("UDP端口错误")
	}
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(fmt.Sprintf("0.0.0.0:%d", port)),
		Port: port,
	})
	if err != nil {
		return fmt.Errorf("UDP listen failed, err:%s", err.Error())
	}
	defer func(listen *net.UDPConn) {
		_ = listen.Close()
		log.Error(fmt.Sprintf("UDP服务器[%s:%d]断开连接", address, port), nil)
	}(listen)
	log.Info(fmt.Sprintf("UDP服务器[%s:%d]连接成功", address, port), nil)
	for {
		var data [4096]byte
		n, addr, err := listen.ReadFromUDP(data[:]) // 接收数据
		if err != nil {
			fmt.Println("read udp failed, err:", err)
			continue
		}
		callback(addr.IP.String(), data[:n])
	}
}

// UdpSend
// @Description: 建立UDP连接并发生数据
// @param address UDP地址 如 127.0.0.1:8080
// @param data
// @return err
func UdpSend(address string, data []byte) (err error) {
	if udpProxy[address] == nil {
		args := strings.Split(address, ":")
		if len(args) == 2 {
			port, _ := strconv.Atoi(args[1])
			socket, _ := net.DialUDP("udp", nil, &net.UDPAddr{
				IP:   net.ParseIP(args[0]),
				Port: port,
			})
			udpProxy[address] = socket
		} else {
			return fmt.Errorf("传递的地址不为合法的IP地址")
		}
	}
	if _, err = udpProxy[address].Write(data); err != nil {
		err = fmt.Errorf("发送数据失败，err:%s", err.Error())
	}
	return
}

// UdpClose
// @Description: 断开指定UDP地址
// @param address
func UdpClose(address string) {
	if udpProxy[address] == nil {
		_ = udpProxy[address].Close()
		delete(udpProxy, address)
	}
	return
}
