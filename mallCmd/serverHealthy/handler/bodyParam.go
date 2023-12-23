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

type ImpHealthy struct {
}

func (imp *ImpHealthy) BodyParamPage(ctx context.Context, input protoHealthy.ReqBodyParamPage) (output protoHealthy.ResBodyParamPage, err error) {
	var info orm.IndexPage
	err = app.Unmarshal(input.Page, &info)
	if err != nil {
		return output, err
	}
	list, total, err := service.BodyParamPage(input.Bmi, info)
	if err != nil {
		return output, err
	}
	output.Total = total
	err = app.Unmarshal(list, &output.List)
	return output, err
}

func (imp *ImpHealthy) BodyParamSave(ctx context.Context, input protoHealthy.BodyParam) (output protoHealthy.BodyParam, err error) {
	var info modelHealthy.BodyParam
	if err = app.Unmarshal(input, &info); err != nil {
		return output, err
	}
	if result, err := service.BodyParamSave(info); err != nil {
		return output, err
	} else {
		err = app.Unmarshal(result, &output)
		return output, err
	}
}

func (imp *ImpHealthy) BodyParamList(ctx context.Context, input protoHealthy.ReqId) (output protoHealthy.BodyParam, err error) {
	result, err := service.BodyParamList(input.Id)
	if err != nil {
		return output, err
	}
	err = app.Unmarshal(result, &output)
	return output, err
}

func (imp *ImpHealthy) BodyParamDelete(ctx context.Context, input protoHealthy.BodyParam) (output protoComm.Result, err error) {
	var info modelHealthy.BodyParam
	err = app.Unmarshal(input, &info)
	if err != nil {
		return output, app.Err(app.Fail, "数据反序列化错误")
	}
	err = service.BodyParamDelete(info)
	return output, err
}
