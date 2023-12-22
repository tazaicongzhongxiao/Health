package main

import (
	"github.com/TarsCloud/TarsGo/tars"
	"gitlab.mall.com/mallBase/basics/pkg/app"
	"gitlab.mall.com/mallBase/basics/pkg/dredis"
	"gitlab.mall.com/mallBase/server/pkg/database/mongo"
	protoHealthy "mall/mallUser/qinwong/protoImp/healthy"
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
