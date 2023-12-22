package handler

import (
	"context"
	"gitlab.mall.com/mallBase/basics/pkg/app"
	"gitlab.mall.com/mallBase/server/pkg/database/orm"
	protoComm "gitlab.mall.com/mallModel/protoImp/base/comm"
	protoHealthy "mall/mallUser/qinwong/protoImp/healthy"
	"serverHealthy/modelHealthy"
	"serverHealthy/service"
)

func (imp *ImpHealthy) FoodPage(ctx context.Context, input protoHealthy.ReqFoodPage) (output protoHealthy.ResFoodPage, err error) {
	var info orm.IndexPage
	err = app.Unmarshal(input.Page, &info)
	if err != nil {
		return output, err
	}
	list, total, err := service.FoodPage(input.Name, info)
	if err != nil {
		return output, err
	}
	output.Total = total
	err = app.Unmarshal(list, &output.List)
	return output, err
}

func (imp *ImpHealthy) FoodSave(ctx context.Context, input protoHealthy.Food) (output protoHealthy.Food, err error) {
	var info modelHealthy.Food
	if err = app.Unmarshal(input, &info); err != nil {
		return output, err
	}
	if result, err := service.FoodSave(info); err != nil {
		return output, err
	} else {
		err = app.Unmarshal(result, &output)
		return output, err
	}
}

func (imp *ImpHealthy) FoodList(ctx context.Context, input protoHealthy.ReqId) (output protoHealthy.Food, err error) {
	result, err := service.FoodList(input.Id)
	if err != nil {
		return output, err
	}
	err = app.Unmarshal(result, &output)
	return output, err
}

func (imp *ImpHealthy) FoodDelete(ctx context.Context, input protoHealthy.Food) (output protoComm.Result, err error) {
	var info modelHealthy.Food
	err = app.Unmarshal(input, &info)
	if err != nil {
		return output, app.Err(app.Fail, "数据反序列化错误")
	}
	err = service.FoodDelete(info)
	return output, err
}
