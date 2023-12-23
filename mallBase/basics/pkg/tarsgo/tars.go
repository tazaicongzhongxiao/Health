package tarsgo

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/pkg/config"
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/model"
	"reflect"
	"strconv"
)

var (
	comm     *tars.Communicator
	servants = make(map[string]string)
)

type Proto struct {
	s model.Servant
}

func init() {
	comm = tars.NewCommunicator()
}

// Start TarsGo 获取
func Start(run func()) {
	if err := config.Config().Bind(app.CfgName, "servants", &servants, func() {
		if config.Config().Bind(app.CfgName, "servants", &servants, nil) == nil {
			go run()
		}
	}); err != nil {
		panic(error.Error(fmt.Errorf("集群 Servants 配置读取错误: %s", err.Error())))
		return
	}
	go run()
}

// GrpcSetServants
// @Description: 设置Servants
// @param m
func GrpcSetServants(m map[string]string) {
	servants = m
	return
}

// GrpcServants
// @Description: 获取Servants
func GrpcServants(name string) string {
	if servants[name] == "" {
		panic(error.Error(fmt.Errorf("集群GRPC servants [%s]未配置", name)))
	}
	return servants[name]
}

// GrpcPluginServants
// @Description: 解析GRPC及插件地址
// @param grpc
// @return proxy
// @return pluginID
func GrpcPluginServants(grpc string) (proxy *tars.ServantProxy, pluginID int32) {
	pluginId, _ := strconv.ParseFloat(grpc, 64)
	if pluginId > 0 {
		proxy = tars.NewServantProxy(comm, GrpcServants("server_plugin_micro"))
		pluginID = int32(pluginId)
	} else {
		proxy = tars.NewServantProxy(comm, grpc)
	}
	return
}

func GrpcServantProxy(name string, t ...int) *tars.ServantProxy {
	p := tars.NewServantProxy(comm, GrpcServants(name))
	if len(t) == 1 {
		p.TarsSetTimeout(t[0])
	}
	return p
}

func GrpcServantProxyGrpc(grpc string, t ...int) *tars.ServantProxy {
	p := tars.NewServantProxy(comm, grpc)
	if len(t) == 1 {
		p.TarsSetTimeout(t[0])
	}
	return p
}

// GRPCList
// @Description: 解析结构体JSON名称
func GRPCList(l any) (res []string) {
	obj := reflect.ValueOf(l)
	t := reflect.TypeOf(obj.Elem().Interface())
	res = make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		if name := t.Field(i).Tag.Get("json"); name != "" {
			res[i] = name
		}
	}
	return res
}
