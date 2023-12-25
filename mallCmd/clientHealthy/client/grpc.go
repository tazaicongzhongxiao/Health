package client

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/pkg/tarsgo"
	protoHealthy "MyTestMall/protoImp/healthy"
)

var GRPC = &struct {
	ServerHealthy protoHealthy.Healthy `json:"server_healthy"`
}{}

func GrpcInit(name string) {
	p := tarsgo.GrpcServantProxy(name)
	switch name {
	case "server_healthy":
		GRPC.ServerHealthy.SetServant(p)
	}
}

func init() {
	tarsgo.Start(func() {
		for _, name := range tarsgo.GRPCList(GRPC) {
			app.Println(4, name)
			GrpcInit(name)
		}
	})
}
