package main

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/pkg/dredis"
	"MyTestMall/mallBase/server/pkg/database/mongo"
	protoHealthy "MyTestMall/protoImp/healthy"
	"github.com/TarsCloud/TarsGo/tars"
	"serverHealthy/handler"
	"serverHealthy/modelHealthy"
)

func init() {
	dredis.Start()
	mongo.Start(modelHealthy.HealthyIndex())
}

func main() {
	new(protoHealthy.Healthy).AddServantWithContext(new(handler.ImpHealthy), app.Cfg.App+"."+app.Cfg.Server+".healthyObj")
	tars.Run()
}
