package controllers

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/client/pkg/comm"
	"MyTestMall/mallBase/client/pkg/validator"
	protoHealthy "MyTestMall/protoImp/healthy"
	"clientHealthy/client"
	"clientHealthy/request"
	"github.com/gin-gonic/gin"
)

func BodyParam(ctx *gin.Context) {
	var form request.PageInfo
	if err := validator.Bind(ctx, &form); !err.IsValid() {
		comm.ValidatorResponse(ctx, err.ErrorsInfo)
		return
	}
	var Form protoHealthy.ReqBodyParamPage
	if err := app.UnmarshalJson(form, &Form.Page); err != nil {
		comm.ApiResponse(ctx, nil, err)
		return
	}
	res, err := client.GRPC.ServerHealthy.BodyParamPage(Form)
	comm.ApiResponse(ctx, &res, err)
}

func BodyParamList(ctx *gin.Context) {
	var form request.ReqId
	if err := validator.Bind(ctx, &form); !err.IsValid() {
		comm.ValidatorResponse(ctx, err.ErrorsInfo)
		return
	}
	res, err := client.GRPC.ServerHealthy.BodyParamList(protoHealthy.ReqId{
		Id: form.Id,
	})
	comm.ApiResponse(ctx, &res, err)
}

func BodyParamSave(ctx *gin.Context) {
	var form request.ReqBodyParam
	if err := validator.Bind(ctx, &form); !err.IsValid() {
		comm.ValidatorResponse(ctx, err.ErrorsInfo)
		return
	}
	var Form protoHealthy.BodyParam
	if err := app.UnmarshalJson(form, &Form); err != nil {
		comm.ApiResponse(ctx, nil, err)
		return
	}
	res, err := client.GRPC.ServerHealthy.BodyParamSave(Form)
	comm.ApiResponse(ctx, &res, err)
}

func BodyParamDelete(ctx *gin.Context) {
	var form request.ReqBodyParam
	if err := validator.Bind(ctx, &form); !err.IsValid() {
		comm.ValidatorResponse(ctx, err.ErrorsInfo)
		return
	}
	var Form protoHealthy.BodyParam
	if err := app.UnmarshalJson(form, &Form); err != nil {
		comm.ApiResponse(ctx, nil, err)
		return
	}
	res, err := client.GRPC.ServerHealthy.BodyParamDelete(Form)
	comm.ApiResponse(ctx, &res, err)
}
