package handler

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/server/pkg/database/orm"
	protoComm "MyTestMall/protoImp/comm"
	protoHealthy "MyTestMall/protoImp/healthy"
	"context"
	"serverHealthy/modelHealthy"
	"serverHealthy/service"
)

func (imp *ImpHealthy) SportsPage(ctx context.Context, input protoHealthy.ReqSportsPage) (output protoHealthy.ResSportsPage, err error) {
	var info orm.IndexPage
	err = app.Unmarshal(input.Page, &info)
	if err != nil {
		return output, err
	}
	list, total, err := service.SportsPage(input.Name, info)
	if err != nil {
		return output, err
	}
	output.Total = total
	err = app.Unmarshal(list, &output.List)
	return output, err
}

func (imp *ImpHealthy) SportsSave(ctx context.Context, input protoHealthy.Sports) (output protoHealthy.Sports, err error) {
	var info modelHealthy.Sports
	if err = app.Unmarshal(input, &info); err != nil {
		return output, err
	}
	if result, err := service.SportsSave(info); err != nil {
		return output, err
	} else {
		err = app.Unmarshal(result, &output)
		return output, err
	}
}

func (imp *ImpHealthy) SportsList(ctx context.Context, input protoHealthy.ReqId) (output protoHealthy.Sports, err error) {
	result, err := service.SportsList(input.Id)
	if err != nil {
		return output, err
	}
	err = app.Unmarshal(result, &output)
	return output, err
}

func (imp *ImpHealthy) SportsDelete(ctx context.Context, input protoHealthy.Sports) (output protoComm.Result, err error) {
	var info modelHealthy.Sports
	err = app.Unmarshal(input, &info)
	if err != nil {
		return output, app.Err(app.Fail, "数据反序列化错误")
	}
	err = service.SportsDelete(info)
	return output, err
}
